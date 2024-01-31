package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type ProductRepository struct {
	db db.PostgresClient
}

func NewProductRepository(db db.PostgresClient) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, dto model.CreateProductDto) fall.Error {
	q := `
  INSERT INTO product (title,description,category_id,brand_id) VALUES ($1,$2,$3,$4);
  `

	_, err := r.db.Exec(ctx, q, dto.Title, dto.Description, dto.CategoryId, dto.BrandId)

	if err != nil {
		return fall.ServerError(msg.ProductCreateError)
	}

	return nil
}

func (r *ProductRepository) CreateModel(ctx context.Context, dto model.CreateProductModelDto, slug string) fall.Error {

	q := fmt.Sprintf(`
	WITH generated AS (
		SELECT uuid_generate_v4() AS generated_article,
		$1::integer as price_param,
		$2::smallint as discount_param,
		$3::text as img_param,
		$4::integer as product_id_param
	)
	INSERT INTO product_model (article, slug, price, discount, main_image_path, product_id)
	SELECT
	LEFT(REPLACE(generated_article::text, '-', ''), 12),
	$5 || '-' || LEFT(REPLACE(generated_article::text, '-', ''), 12),
	price_param,discount_param,img_param,product_id_param
	FROM generated;
	`)

	_, err := r.db.Exec(ctx, q, dto.Price, dto.Discount, dto.ImagePath, dto.ProductId, slug)

	if err != nil {
		log.Println(err.Error())
		return fall.NewErr(msg.ProductCreateModelError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, id int) fall.Error {
	query := `
	DELETE FROM product WHERE product_id = $1;
	`
	_, err := pr.db.Exec(ctx, query, id)

	if err != nil {
		return fall.ServerError(fmt.Sprintf("%s, details: %s", msg.ProductDeleteError, err.Error()))
	}
	return nil
}

func (pr *ProductRepository) DeleteProductModel(ctx context.Context, id int) fall.Error {
	query := `
	DELETE FROM product_model WHERE product_model_id = $1;
	`

	_, err := pr.db.Exec(ctx, query, id)

	if err != nil {
		return fall.ServerError(fmt.Sprintf("%s, details: %s", msg.ProductModelDeleteError, err.Error()))
	}
	return nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, dto model.UpdateProductDto, id int) fall.Error {

	var queries []string

	if dto.Description != nil {
		queries = append(queries, fmt.Sprintf("description = %s", *dto.Description))
	}

	if dto.Title != nil {
		queries = append(queries, fmt.Sprintf("title = %s", *dto.Title))
	}

	if len(queries) > 0 {
		q := "UPDATE product SET " + strings.Join(queries, ",") + " WHERE product_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.ProductUpdateError, err.Error()))
		}
	}

	return nil
}

func (r *ProductRepository) UpdateProductModel(ctx context.Context, dto model.UpdateProductModelDto, modelId int) fall.Error {

	var queries []string

	if dto.Discount != nil {
		queries = append(queries, fmt.Sprintf("discount = %d", *dto.Discount))
	}

	if dto.Price != nil {
		queries = append(queries, fmt.Sprintf("price = %d", *dto.Price))
	}

	if dto.ImagePath != nil {
		queries = append(queries, fmt.Sprintf("main_image_path = '%s'", *dto.ImagePath))
	}

	if len(queries) > 0 {
		q := "UPDATE product_model SET " + strings.Join(queries, ",") + " WHERE product_model_id = $1;"
		_, err := r.db.Exec(ctx, q, modelId)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.ProductModelUpdateError, err.Error()))
		}
	}

	return nil
}

func (r *ProductRepository) FindProductById(ctx context.Context, id int) (*model.Product, fall.Error) {
	query := `
	select p.product_id as p_id, p.title as p_title, p.description as p_descr,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug, c.short_title as c_short_title,
	c.img_path as c_img_path, c.parent_category_id as c_parent_id,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, b.img_path as b_img_path, b.description as d_description
	from product as p
	inner join category as c on p.category_id = c.category_id
	inner join brand as b on p.brand_id = b.brand_id
	where product_id = $1; 	
	`
	row := r.db.QueryRow(ctx, query, id)

	p := model.Product{}

	err := row.Scan(&p.Id, &p.Title, &p.Description, &p.Category.Id, &p.Category.Title, &p.Category.Slug,
		&p.Category.ShortTitle, &p.Category.ImgPath, &p.Category.ParentId, &p.Brand.Id, &p.Brand.Title, &p.Brand.Slug, &p.Brand.ImgPath, &p.Brand.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.ProductNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &p, nil
}

func (r *ProductRepository) FindProductModelById(ctx context.Context, id int) (*model.ProductModel, fall.Error) {
	return r.findProductModelByField(ctx, "product_model_id", id)
}

func (r *ProductRepository) FindProductModelBySlug(ctx context.Context, slug string) (*model.ProductModel, fall.Error) {
	return r.findProductModelByField(ctx, "slug", slug)
}

func (r *ProductRepository) FindProductModelByArticle(ctx context.Context, article string) (*model.ProductModel, fall.Error) {
	return r.findProductModelByField(ctx, "article", article)
}

func (r *ProductRepository) findProductModelByField(ctx context.Context, field string, value any) (*model.ProductModel, fall.Error) {
	query := fmt.Sprintf(`
	SELECT product_model_id,price,discount,main_image_path, slug, article, product_id FROM product_model WHERE %s = $1; 	
	`, field)
	row := r.db.QueryRow(ctx, query, value)

	m := model.ProductModel{}

	err := row.Scan(&m.Id, &m.Price, &m.Discount, &m.ImagePath, &m.Slug, &m.Article, &m.ProductId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.ProductModelNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &m, nil
}

func (pr *ProductRepository) AddPhoto(ctx context.Context, dto model.CreateProducModelImg) fall.Error {
	query := "INSERT INTO product_model_img (img_path, product_model_id) VALUES ($1,$2);"

	_, err := pr.db.Exec(ctx, query, dto.ImgPath, dto.ProductModelId)
	if err != nil {
		return fall.NewErr(msg.ProductAddPhotoError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (pr *ProductRepository) RemovePhoto(ctx context.Context, photoId int) fall.Error {
	query := "DELETE FROM product_model_img WHERE product_img_id = $1;"

	_, err := pr.db.Exec(ctx, query, photoId)
	if err != nil {
		return fall.NewErr(msg.ProductAddPhotoError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (pr *ProductRepository) GetProductPage(ctx context.Context, slug string) (*model.ProductRelation, fall.Error) {

	query := `
	select p.product_id as p_id, p.title as p_title,
	p.description as p_description,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug, c.short_title as c_short_title,
	c.img_path as c_img_path, c.parent_category_id as c_parent_id,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, b.img_path as b_img_path, b.description as d_description,
	pm.product_model_id as pm_id, pm.slug as pm_slug, pm.article as pm_article, pm.price as pm_price, pm.discount as pm_discount, pm.product_id as pm_product_id, pm.main_image_path as pm_image_main,
	pimg.product_img_id as pimg_id, pimg.img_path as pimg_img_path,pimg.product_model_id as pimg_model_id,
	op.option_id as op_id, op.title as op_title, op.slug as op_slug, pmop.product_model_id as pmop_model_id,
	v.option_value_id as v_id, v.value as v_value, v.info as v_info, v.option_id as v_option_id, pmop.product_model_id as pmop_model_id_v,
  sz.size_id as size_id, sz.size_value as sz_value, ms.literal_size as ls, ms.in_stock as in_stock, ms.product_model_id as size_model_id,
	ms.model_size_id as order_model_size
	from product as p
	inner join category as c on p.category_id = c.category_id
	inner join brand as b on p.brand_id = b.brand_id
	inner join product_model as pm on pm.product_id = p.product_id
	inner join product_model_img as pimg on pimg.product_model_id = pm.product_model_id
	inner join product_model_option as pmop on pmop.product_model_id = pm.product_model_id
	inner join option as op on op.option_id = pmop.option_id
	inner join option_value as v on v.option_value_id = pmop.option_value_id
  inner join model_sizes as ms on ms.product_model_id = pm.product_model_id
  inner join sizes as sz on ms.size_id = sz.size_id
	where pm.slug = $1;
	`
	rows, err := pr.db.Query(ctx, query, slug)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	p := model.ProductRelation{}
	productModel := model.ProductModelRelation{}
	founded := false

	imagesMap := make(map[int]model.ProductModelImg)
	var imagesOrder []int
	optionsMap := make(map[int]*model.ProductModelOption)
	var optionsOrder []int
	valuesMap := make(map[int]model.ProductModelOptionValue)
	var valuesOrder []int
	sizesMap := make(map[int]model.ProductModelSize)

	for rows.Next() {

		productModelImg := model.ProductModelImg{}
		opt := model.ProductModelOption{}
		val := model.ProductModelOptionValue{}
		size := model.ProductModelSize{}

		err := rows.Scan(&p.Id, &p.Title, &p.Description,
			&p.Category.Id, &p.Category.Title, &p.Category.Slug, &p.Category.ShortTitle, &p.Category.ImgPath, &p.Category.ParentId,
			&p.Brand.Id, &p.Brand.Title, &p.Brand.Slug, &p.Brand.ImgPath, &p.Brand.Description,
			&productModel.Id, &productModel.Slug, &productModel.Article, &productModel.Price, &productModel.Discount, &productModel.ProductId, &productModel.ImagePath,
			&productModelImg.Id, &productModelImg.ImgPath, &productModelImg.ProductModelId,
			&opt.Id, &opt.Title, &opt.Slug, &opt.ProductModelId,
			&val.Id, &val.Value, &val.Info, &val.OptionId, &val.ProductModelId,
			&size.SizeId, &size.Value, &size.Literal, &size.InStock, &size.ModelId, &size.SizeModelId,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		_, ok := imagesMap[productModelImg.Id]
		if !ok {
			imagesMap[productModelImg.Id] = productModelImg
			imagesOrder = append(imagesOrder, productModelImg.Id)
		}

		_, ok = optionsMap[opt.Id]
		if !ok {
			optionsMap[opt.Id] = &opt
			optionsOrder = append(optionsOrder, opt.Id)
		}

		_, ok = valuesMap[val.Id]
		if !ok {
			valuesMap[val.Id] = val
			valuesOrder = append(valuesOrder, val.Id)
		}

		_, ok = sizesMap[size.SizeId]
		if !ok {
			sizesMap[size.SizeId] = size
		}

		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.ProductNotFound, fall.STATUS_NOT_FOUND)
	}

	for _, key := range optionsOrder {
		value := valuesMap[key]
		opt := optionsMap[value.OptionId]
		opt.Values = append(opt.Values, value)
	}

	options := make([]*model.ProductModelOption, 0, len(optionsMap))

	for _, key := range optionsOrder {
		opt := optionsMap[key]
		options = append(options, opt)
	}

	images := make([]model.ProductModelImg, 0, len(imagesMap))

	for _, key := range imagesOrder {
		img := imagesMap[key]
		images = append(images, img)
	}

	sizes := make([]model.ProductModelSize, 0, len(sizesMap))

	for _, v := range sizesMap {
		sizes = append(sizes, v)
	}

	sort.Slice(sizes, func(i, j int) bool {
		a := sizes[i].Value
		b := sizes[j].Value
		return a < b
	})

	productModel.Images = images
	productModel.Options = options
	productModel.Sizes = sizes
	p.CurrentModel = productModel

	return &p, nil
}

func (pr *ProductRepository) FindModelsColored(ctx context.Context, id int) ([]model.ProductModelColors, fall.Error) {
	query := `
	select 
  pm.product_model_id as pm_id, pm.slug as pm_slug, pm.main_image_path as pm_img, 
  v.value as v_value, pmop.product_model_id as pmop_model_id_v
	from product as p
	inner join product_model as pm on pm.product_id = p.product_id
	inner join product_model_option as pmop on pmop.product_model_id = pm.product_model_id
	inner join option as op on op.option_id = pmop.option_id
	inner join option_value as v on v.option_value_id = pmop.option_value_id
	where p.product_id = $1 and op.slug = 'color';
	`
	rows, err := pr.db.Query(ctx, query, id)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var models []model.ProductModelColors

	found := false

	for rows.Next() {

		model := model.ProductModelColors{}

		err := rows.Scan(&model.Id, &model.Slug, &model.Image, &model.Color.Value, &model.Color.ModelId)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		models = append(models, model)

		if !found {
			found = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !found {
		return []model.ProductModelColors{}, nil
	}

	return models, nil

}

func (r *ProductRepository) AdminGetProducts(page int, brandId *int, categoryId *int) (*model.AdminProductResponse, fall.Error) {
	//todo add sort
	limit := 32

	offset := page*limit - limit

	whereCategory := ""
	whereBrand := ""

	if categoryId != nil {

		whereCategory = fmt.Sprintf("WHERE category_id = %d", *categoryId)
	}

	if brandId != nil {

		whereBrand = fmt.Sprintf("WHERE b.brand_id = %d", *brandId)
	}

	query := fmt.Sprintf(`
	WITH RECURSIVE category_tree AS (
		SELECT category_id, title, slug, short_title,img_path,parent_category_id
		FROM category
		%s
		UNION ALL
		SELECT c.category_id, c.title, c.slug, c.short_title, c.img_path, c.parent_category_id
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT p.product_id as p_id, p.title as p_title, p.description as p_descr,
	ct.category_id as c_id, ct.title as c_title, ct.slug as c_slug, ct.short_title as c_short,
	ct.img_path as c_img, ct.parent_category_id as ct_parent_id,
	b.brand_id as b_id, b.title as b_title,  b.slug as b_slug,  b.img_path as b_img, b.description as b_description,
	m.product_model_id as m_id,  m.slug as m_slug, m.article as m_article, m.price as m_price, m.discount as m_discount, m.main_image_path as m_img, 
	m.product_id as m_pid,
	(select count(distinct p.product_id)
		from product as p
		inner join category_tree as ct on ct.category_id = p.category_id
		inner join brand as b on b.brand_id = p.brand_id
		left join product_model as m on p.product_id = m.product_id
		%s
	) as total
	from product as p
	inner join category_tree as ct on ct.category_id = p.category_id
	inner join brand as b on b.brand_id = p.brand_id
	left join product_model as m on p.product_id = m.product_id
	%s
	LIMIT $1 OFFSET $2;
	`, whereCategory, whereBrand, whereBrand)
	rows, err := r.db.Query(context.Background(), query, limit, offset)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	var founded bool = false

	modelsMap := make(map[int]model.AdminProductModelRelation)
	var modelsOrder []int
	productsMap := make(map[int]*model.AdminProduct)
	var productOrder []int
	var totalCount int
	for rows.Next() {

		product := model.AdminProduct{}
		model := model.AdminProductModelRelation{}

		err := rows.Scan(&product.Id, &product.Title, &product.Description, &product.Category.Id,
			&product.Category.Title, &product.Category.Slug, &product.Category.ShortTitle, &product.Category.ImgPath, &product.Category.ParentId,
			&product.Brand.Id, &product.Brand.Title, &product.Brand.Slug, &product.Brand.ImgPath, &product.Brand.Description,
			&model.Id, &model.Slug, &model.Article, &model.Price, &model.Discount, &model.ImagePath, &model.ProductId, &totalCount,
		)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		_, ok := productsMap[product.Id]
		if !ok {
			productsMap[product.Id] = &product
			productOrder = append(productOrder, product.Id)
		}

		if model.Id != nil {
			_, ok := modelsMap[*model.Id]
			if !ok {
				modelsMap[*model.Id] = model
				modelsOrder = append(modelsOrder, *model.Id)
			}
		}

		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !founded {
		return &model.AdminProductResponse{Total: 0, Products: []*model.AdminProduct{}}, nil
	}

	for _, key := range modelsOrder {
		m := modelsMap[key]
		p := productsMap[*m.ProductId]
		p.Models = append(p.Models, m)
	}

	result := make([]*model.AdminProduct, 0, len(productsMap))

	for _, v := range productOrder {
		p := productsMap[v]
		result = append(result, p)
	}
	return &model.AdminProductResponse{Products: result, Total: totalCount}, nil
}
