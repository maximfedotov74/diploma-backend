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
	"github.com/maximfedotov74/diploma-backend/internal/shared/generator"
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

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) fall.Error {
	query := `
	DELETE FROM product WHERE product_id = $1;
	`
	_, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fall.ServerError(fmt.Sprintf("%s, details: %s", msg.ProductDeleteError, err.Error()))
	}
	return nil
}

func (r *ProductRepository) DeleteProductModel(ctx context.Context, id int) fall.Error {
	query := `
	DELETE FROM product_model WHERE product_model_id = $1;
	`

	_, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fall.ServerError(fmt.Sprintf("%s, details: %s", msg.ProductModelDeleteError, err.Error()))
	}
	return nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, dto model.UpdateProductDto, id int) fall.Error {

	var queries []string

	if dto.Description != nil {
		queries = append(queries, fmt.Sprintf("description = '%s'", *dto.Description))
	}

	if dto.Title != nil {
		queries = append(queries, fmt.Sprintf("title = '%s'", *dto.Title))
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

func (r *ProductRepository) AddPhoto(ctx context.Context, dto model.CreateProducModelImg) fall.Error {
	query := "INSERT INTO product_model_img (img_path, product_model_id) VALUES ($1,$2);"

	_, err := r.db.Exec(ctx, query, dto.ImgPath, dto.ProductModelId)
	if err != nil {
		return fall.NewErr(msg.ProductAddPhotoError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (r *ProductRepository) RemovePhoto(ctx context.Context, photoId int) fall.Error {
	query := "DELETE FROM product_model_img WHERE product_img_id = $1;"

	_, err := r.db.Exec(ctx, query, photoId)
	if err != nil {
		return fall.NewErr(msg.ProductAddPhotoError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (r *ProductRepository) GetProductPage(ctx context.Context, slug string) (*model.ProductRelation, fall.Error) {

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
	rows, err := r.db.Query(ctx, query, slug)

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

func (r *ProductRepository) AdminGetProductModels(ctx context.Context, id int) ([]model.AdminProductModelRelation, fall.Error) {

	q := "SELECT product_model_id,price,slug,article,discount,main_image_path,product_id FROM product_model WHERE product_id = $1;"

	rows, err := r.db.Query(ctx, q, id)
	defer rows.Close()

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	var models = []model.AdminProductModelRelation{}

	for rows.Next() {
		m := model.AdminProductModelRelation{}

		err := rows.Scan(&m.Id, &m.Price, &m.Slug, &m.Article, &m.Discount, &m.ImagePath, &m.ProductId)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		models = append(models, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return models, nil
}

func (r *ProductRepository) AdminGetProducts(ctx context.Context, page int, brandId *int, categoryId *int) (*model.AdminProductResponse, fall.Error) {
	//todo add sort
	limit := 8

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
	SELECT distinct p.product_id as p_id, p.title as p_title, p.description as p_descr,
	ct.category_id as c_id, ct.title as c_title, ct.slug as c_slug, ct.short_title as c_short,
	ct.img_path as c_img, ct.parent_category_id as ct_parent_id,
	b.brand_id as b_id, b.title as b_title,  b.slug as b_slug,  b.img_path as b_img, b.description as b_description,
	(select count(distinct p.product_id)
		from product as p
		inner join category_tree as ct on ct.category_id = p.category_id
		inner join brand as b on b.brand_id = p.brand_id
		%s
	) as total
	from product as p
	inner join category_tree as ct on ct.category_id = p.category_id
	inner join brand as b on b.brand_id = p.brand_id
	%s
	ORDER BY p.product_id
	LIMIT $1 OFFSET $2;
	`, whereCategory, whereBrand, whereBrand)
	rows, err := r.db.Query(ctx, query, limit, offset)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	var founded bool = false

	var products []*model.AdminProduct
	var totalCount int
	for rows.Next() {

		product := model.AdminProduct{}

		err := rows.Scan(&product.Id, &product.Title, &product.Description, &product.Category.Id,
			&product.Category.Title, &product.Category.Slug, &product.Category.ShortTitle, &product.Category.ImgPath, &product.Category.ParentId,
			&product.Brand.Id, &product.Brand.Title, &product.Brand.Slug, &product.Brand.ImgPath, &product.Brand.Description, &totalCount,
		)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		products = append(products, &product)

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

	return &model.AdminProductResponse{Products: products, Total: totalCount}, nil
}

func (r *ProductRepository) GetCatalogModels(ctx context.Context, categorySlug string, sql generator.GeneratedCatalogQuery) (*model.CatalogResponse, fall.Error) {

	// TODO: сначала получить id подходящих моделей, а потом их получить и вернуть.

	mainJoins := `FROM product p INNER JOIN category_tree ct ON p.category_id = ct.category_id 
	INNER JOIN brand b on p.brand_id = b.brand_id
	INNER JOIN product_model pm ON pm.product_id = p.product_id
	inner join model_sizes ms on ms.product_model_id = pm.product_model_id
	inner join sizes sz on ms.size_id = sz.size_id
	inner join product_model_img as pimg on pimg.product_model_id = pm.product_model_id
	`

	query := fmt.Sprintf(`
	WITH RECURSIVE category_tree AS (
		SELECT category_id, title, slug, parent_category_id
		FROM category
		WHERE slug = $1
		UNION ALL
		SELECT c.category_id, c.title, c.slug, c.parent_category_id
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT p.product_id as p_id, p.title as p_title,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, ct.category_id as ct_id, ct.title as ct_title, ct.slug as ct_slug,
	pm.product_model_id as model_id, pm.slug as m_slug, pm_artice as m_artice, pm.price as model_price, pm.discount as model_discount,
	pm.main_image_path as pm_main_img,
	pimg.product_img_id as pimg_id, pimg.product_model_id as pimg_model_id, pimg.img_path as pimg_img_path, 
	sz.size_id as size_id, sz.size_value as size_value, ms.literal_size as literal_size,
	ms.product_model_id as ms_pm_id, ms.in_stock as ms_in_stock,
	ms.model_size_id as ms_m_sz_id,
	(select count(distinct pm.product_model_id)%s %s
	) as total_count
	%s %s %s %s;`, mainJoins, sql.MainQuery, mainJoins, sql.MainQuery, sql.SortStatement, sql.Pagination)

	rows, err := r.db.Query(ctx, query, categorySlug)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	imagesMap := make(map[int]*model.ProductModelImg)
	sizesMap := make(map[int]*model.ProductModelSize)
	modelsMap := make(map[int]*model.CatalogProductModel)
	var total int
	var modelOrder []int
	var imgOrder []int
	var sizeOrder []int

	for rows.Next() {
		sz := model.ProductModelSize{}
		img := model.ProductModelImg{}
		m := model.CatalogProductModel{}

		err := rows.Scan(&m.ProductId, &m.Title, &m.Brand.Id, &m.Brand.Title, &m.Brand.Slug,
			&m.Category.Id, &m.Category.Title, &m.Category.Slug, &m.ModelId, &m.Slug, &m.Article, &m.Price, &m.Discount,
			&m.MainImagePath, &img.Id, &img.ProductModelId, &img.ImgPath, &sz.SizeId, &sz.Value, &sz.Literal, &sz.ModelId, &sz.InStock, &sz.SizeModelId, &total,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		_, ok := modelsMap[m.ModelId]
		if !ok {
			modelsMap[m.ModelId] = &m
			modelOrder = append(modelOrder, m.ModelId)
		}
		_, ok = imagesMap[img.Id]
		if !ok {
			imagesMap[img.Id] = &img
			imgOrder = append(imgOrder, img.Id)
		}
		_, ok = sizesMap[sz.SizeModelId]
		if !ok {
			sizesMap[sz.SizeModelId] = &sz
			sizeOrder = append(sizeOrder, sz.SizeModelId)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	for _, v := range imgOrder {
		img := imagesMap[v]
		m := modelsMap[img.ProductModelId]
		m.Images = append(m.Images, img)
	}

	for _, v := range sizeOrder {
		sz := sizesMap[v]
		m := modelsMap[sz.ModelId]
		m.Sizes = append(m.Sizes, sz)
	}

	var result []*model.CatalogProductModel

	for _, id := range modelOrder {
		m := modelsMap[id]
		result = append(result, m)
	}

	return &model.CatalogResponse{
		Models:     result,
		TotalCount: total,
	}, nil

}

func (r *ProductRepository) FindModelSizeById(ctx context.Context, id int) (*model.OrderProductModelSize, fall.Error) {
	query := `
	SELECT ms.model_size_id,ms.product_model_id,ms.in_stock, m.price, m.discount
	FROM model_sizes as ms
	INNER JOIN product_model as m ON ms.product_model_id = m.product_model_id
	WHERE ms.model_size_id = $1;
	`
	row := r.db.QueryRow(ctx, query, id)
	m := model.OrderProductModelSize{}

	err := row.Scan(&m.SizeModelId, &m.ModelId, &m.InStock, &m.Price, &m.Discount)
	if err != nil {
		return nil, fall.NewErr(msg.ProductNotFound, fall.STATUS_NOT_FOUND)
	}
	return &m, nil
}

func (r *ProductRepository) ReduceQuantityInStock(ctx context.Context, modelSizeId int, quantity int, tx db.Transaction) fall.Error {
	m, ex := r.FindModelSizeById(ctx, modelSizeId)
	if ex != nil {
		return ex
	}

	newQuantity := m.InStock - quantity

	if newQuantity < 0 {
		return fall.NewErr(msg.ProductInStockCannotBeLessThanZero, fall.STATUS_BAD_REQUEST)
	}

	query := "UPDATE model_sizes SET in_stock = $1 WHERE model_size_id = $2"

	if tx != nil {
		_, err := tx.Exec(ctx, query, newQuantity, modelSizeId)
		if err != nil {
			return fall.ServerError(err.Error())
		}
		return nil
	}
	_, err := r.db.Exec(ctx, query, newQuantity, modelSizeId)
	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}

func (r *ProductRepository) ReturnQuantityInStock(ctx context.Context, modelSizeId int, quantity int, tx db.Transaction) fall.Error {
	m, ex := r.FindModelSizeById(ctx, modelSizeId)

	if ex != nil {
		return ex
	}

	newQuantity := m.InStock + quantity

	query := "UPDATE model_sizes SET in_stock = $1 WHERE model_size_id = $2"

	if tx != nil {
		_, err := tx.Exec(ctx, query, newQuantity, modelSizeId)
		if err != nil {
			return fall.ServerError(err.Error())
		}
		return nil
	}
	_, err := r.db.Exec(ctx, query, newQuantity, modelSizeId)
	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}
