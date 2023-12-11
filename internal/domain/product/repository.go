package product

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) CreateProduct(dto CreateProductDto, slug string) exception.Error {
	query := `
  INSERT INTO product (title, slug, description, category_id, brand_id)
  VALUES ($1, $2, $3, $4, $5);
  `

	_, err := pr.db.Exec(context.Background(), query, dto.Title, slug, dto.Description, dto.CategoryID, dto.BrandID)
	if err != nil {
		return exception.NewErr(productCreateError, exception.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (pr *ProductRepository) FindById(id int) (*ProductWithoutRelations, exception.Error) {
	query := `
	select p.product_id as p_id, p.title as p_title, p.slug as p_slug, p.description as p_descr,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug, c.short_title as c_short_title,
	c.img_path as c_img_path,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, b.img_path as b_img_path
	from product as p
	inner join category as c on p.category_id = c.category_id
	inner join brand as b on p.brand_id = b.brand_id
	where product_id = $1; 	
	`
	row := pr.db.QueryRow(context.Background(), query, id)

	p := ProductWithoutRelations{}

	err := row.Scan(&p.Id, &p.Title, &p.Slug, &p.Description, &p.Category.Id, &p.Category.Title, &p.Category.Slug,
		&p.Category.ShortTitle, &p.Category.ImgPath, &p.Brand.Id, &p.Brand.Title, &p.Brand.Slug, &p.Brand.ImgPath,
	)
	if err != nil {
		return nil, exception.NewErr(productNotFound, exception.STATUS_NOT_FOUND)
	}

	return &p, nil
}

func (pr *ProductRepository) GetCatalogModels(categorySlug string, sql utils.GeneratedCatalogQuery) string {

	mainJoins := `FROM product p INNER JOIN category_tree ct ON p.category_id = ct.category_id 
	INNER JOIN brand b on p.brand_id = b.brand_id
	INNER JOIN product_model pm ON pm.product_id = p.product_id
	inner join model_sizes ms on ms.product_model_id = pm.product_model_id
	inner join sizes sz on ms.size_id = sz.size_id
	left join product_model_img as pimg on pimg.product_model_id = pm.product_model_id
	`

	query := fmt.Sprintf(`
	WITH RECURSIVE category_tree AS (
		SELECT category_id, title, slug, parent_category_id
		FROM category
		WHERE slug = 'men'
		UNION ALL
		SELECT c.category_id, c.title, c.slug, c.parent_category_id
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT p.title as product_title, p.slug as product_slug,
	b.title as brand_title, pm.product_model_id as model_id, pm.price as model_price, pm.discount as model_discount,
	pm.main_image_path as pm_img,
	pimg.product_img_id as pimg_id, pimg.img_path as pimg_img_path, pimg.product_model_id as pimg_model_id,
	sz.size_id as size_id, sz.size_value as size_value,
	(select count(distinct pm.product_model_id)%s %s
	) as total_count,
(select round(avg(f.rate), 2) from feedback as f
inner join product_model fpm on fpm.product_model_id = f.product_model_id
where fpm.product_model_id = pm.product_model_id
) as avg_rate,
(select count(f.feedback_id) from feedback as f
inner join product_model fpm on fpm.product_model_id = f.product_model_id
where fpm.product_model_id = pm.product_model_id
)as rate_count
	%s %s %s;`, mainJoins, sql.MainQuery, mainJoins, sql.MainQuery, sql.SortStatement)

	return query
}

func (pr *ProductRepository) FindModelsColored(slug string) ([]ProductModelColors, exception.Error) {
	query := `
	select 
  pm.product_model_id as pm_id, p.slug as p_slug, pm.main_image_path as pm_img 
  v.value as v_value, pmop.product_model_id as pmop_model_id_v
	from product as p
	inner join product_model as pm on pm.product_id = p.product_id
	inner join product_model_option as pmop on pmop.product_model_id = pm.product_model_id
	inner join option as op on op.option_id = pmop.option_id
	inner join option_value as v on v.option_value_id = pmop.option_value_id
	where p.slug = $1 and op.slug = 'color';
	`
	rows, err := pr.db.Query(context.Background(), query, slug)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	defer rows.Close()

	modelMap := make(map[int]ProductModelColors)

	found := false

	for rows.Next() {

		model := ProductModelColors{}

		err := rows.Scan(&model.Id, &model.ProductSlug, &model.Image, &model.Color.Value, &model.Color.ModelId)

		if err != nil {
			return nil, exception.ServerError(err.Error())
		}

		modelMap[*model.Id] = model

		if !found {
			found = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	if !found {
		return []ProductModelColors{}, nil
	}

	models := make([]ProductModelColors, 0, len(modelMap))

	for _, v := range modelMap {
		models = append(models, v)
	}

	return models, nil

}

func (pr *ProductRepository) FindModelById(modelId int) (*ProductModelWithoutRelations, exception.Error) {
	query := `
	select pm.product_model_id, pm.price, pm.discount, pm.main_image_path, pm.product_id
	from product_model as pm
	where pm.product_model_id = $1;
	`

	row := pr.db.QueryRow(context.Background(), query, modelId)

	model := ProductModelWithoutRelations{}

	err := row.Scan(&model.Id, &model.Price, &model.Discount, &model.ImagePath, &model.ProductId)

	if err != nil {
		return nil, exception.NewErr(productNotFound, exception.STATUS_NOT_FOUND)
	}

	return &model, nil

}

func (pr *ProductRepository) FindModelByIdWithRelations(modelId int) (*Product, exception.Error) {
	query := `
	select p.product_id as p_id, p.title as p_title, p.slug as p_slug,
	p.description as p_description,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug, c.short_title as c_short_title,
	c.img_path as c_img_path,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, b.img_path as b_img_path,
	pm.product_model_id as pm_id, pm.price as pm_price, pm.discount as pm_discount, pm.product_id as pm_product_id, pm.main_image_path as pm_image_main,
	pimg.product_img_id as pimg_id, pimg.img_path as pimg_img_path,pimg.product_model_id as pimg_model_id,
	op.option_id as op_id, op.title as op_title, op.slug as op_slug, pmop.product_model_id as pmop_model_id,
	v.option_value_id as v_id, v.value as v_value, v.info as v_info, v.option_id as v_option_id, pmop.product_model_id as pmop_model_id_v,
  sz.size_id as size_id, ms.literal_size as ls, ms.in_stock as in_stock, ms.product_model_id as size_model_id,
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
	where pm.product_model_id = $1;
	`
	rows, err := pr.db.Query(context.Background(), query, modelId)

	if err != nil {

		return nil, exception.ServerError(err.Error())
	}

	defer rows.Close()

	p := Product{}
	productModel := ProductModel{}
	founded := false

	imagesMap := make(map[int]ProductModelImg)
	optionsMap := make(map[int]ProductModelOption)
	valuesMap := make(map[int]ProductModelOptionValue)
	sizesMap := make(map[int]ProductModelSize)

	for rows.Next() {

		productModelImg := ProductModelImg{}
		opt := ProductModelOption{}
		val := ProductModelOptionValue{}
		size := ProductModelSize{}

		err := rows.Scan(&p.Id, &p.Title, &p.Slug, &p.Description,
			&p.Category.Id, &p.Category.Title, &p.Category.Slug, &p.Category.ShortTitle, &p.Category.ImgPath,
			&p.Brand.Id, &p.Brand.Title, &p.Brand.Slug, &p.Brand.ImgPath,
			&productModel.Id, &productModel.Price, &productModel.Discount, &productModel.ProductId, &productModel.ImagePath,
			&productModelImg.Id, &productModelImg.ImgPath, &productModelImg.ProductModelId,
			&opt.Id, &opt.Title, &opt.Slug, &opt.ProductModelId,
			&val.Id, &val.Value, &val.Info, &val.OptionId, &val.ProductModelId,
			&size.SizeId, &size.Literal, &size.InStock, &size.ModelId, &size.SizeModelId,
		)

		if err != nil {

			return nil, exception.ServerError(err.Error())
		}

		imagesMap[*productModelImg.Id] = productModelImg

		optionsMap[*opt.Id] = opt

		valuesMap[*val.Id] = val

		sizesMap[size.SizeId] = size

		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, exception.NewErr(productNotFound, exception.STATUS_NOT_FOUND)
	}

	for _, v := range valuesMap {
		opt := optionsMap[*v.OptionId]
		opt.Values = append(opt.Values, v)
		optionsMap[*v.OptionId] = opt
	}

	options := make([]ProductModelOption, 0, len(optionsMap))

	for _, v := range optionsMap {
		options = append(options, v)
	}

	images := make([]ProductModelImg, 0, len(imagesMap))

	for _, v := range imagesMap {
		images = append(images, v)
	}

	sizes := make([]ProductModelSize, 0, len(sizesMap))

	for _, v := range sizesMap {
		sizes = append(sizes, v)
	}

	productModel.Images = images
	productModel.Options = options
	productModel.Sizes = sizes
	p.CurrentModel = &productModel

	return &p, nil
}

func (pr *ProductRepository) CreateModel(dto CreateProductModelDto) exception.Error {
	query := "INSERT INTO product_model (price, discount, main_image_path, product_id) VALUES ($1, $2, $3, $4);"

	_, err := pr.db.Exec(context.Background(), query, dto.Price, dto.Discount, dto.ImagePath, dto.ProductId)

	if err != nil {
		return exception.NewErr(productCreateError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (pr *ProductRepository) AddPhoto(dto CreateProducModelImg) exception.Error {
	query := "INSERT INTO product_model_img (img_path, product_model_id) VALUES ($1,$2);"

	_, err := pr.db.Exec(context.Background(), query, dto.ImgPath, dto.ProductModelId)
	if err != nil {
		return exception.NewErr(productAddPhotoError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (pr *ProductRepository) RemovePhoto(photoId int) exception.Error {
	query := "DELETE FROM product_model_img WHERE product_img_id = $1;"

	_, err := pr.db.Exec(context.Background(), query, photoId)
	if err != nil {
		return exception.NewErr(productAddPhotoError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (pr *ProductRepository) AdminGetProducts(page int, brandId *int, categoryId *int) (*AdminProductResponse, exception.Error) {
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
		SELECT category_id, title, slug, short_title,img_path
		FROM category
		%s
		UNION ALL
		SELECT c.category_id, c.title, c.slug, c.short_title, c.img_path
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT p.product_id as p_id, p.title as p_title, p.slug as p_slug, p.description as p_descr,
	ct.category_id as c_id, ct.title as c_title, ct.slug as c_slug, ct.short_title as c_short, ct.img_path as c_img,
	b.brand_id as b_id, b.title as b_title,  b.slug as b_slug,  b.img_path as b_img,
	m.product_model_id as m_id, m.price as m_price, m.discount as m_discount, m.main_image_path as m_img, 
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
	rows, err := pr.db.Query(context.Background(), query, limit, offset)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	defer rows.Close()

	var founded bool = false

	modelsMap := make(map[int]AdminProductModelRelation)
	productsMap := make(map[int]AdminProduct)
	var totalCount int
	for rows.Next() {

		product := AdminProduct{}
		model := AdminProductModelRelation{}

		err := rows.Scan(&product.Id, &product.Title, &product.Slug, &product.Description, &product.Category.Id,
			&product.Category.Title, &product.Category.Slug, &product.Category.ShortTitle, &product.Category.ImgPath,
			&product.Brand.Id, &product.Brand.Title, &product.Brand.Slug, &product.Brand.ImgPath,
			&model.Id, &model.Price, &model.Discount, &model.ImagePath, &model.ProductId, &totalCount,
		)

		if err != nil {
			return nil, exception.ServerError(err.Error())
		}

		productsMap[product.Id] = product

		if model.Id != nil {
			modelsMap[*model.Id] = model
		}

		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	if !founded {
		return &AdminProductResponse{Total: 0, Products: []AdminProduct{}}, nil
	}

	for _, v := range modelsMap {
		p := productsMap[*v.ProductId]
		p.Models = append(p.Models, v)
		productsMap[p.Id] = p
	}

	result := make([]AdminProduct, 0, len(productsMap))

	for _, v := range productsMap {
		result = append(result, v)
	}

	return &AdminProductResponse{Products: result, Total: totalCount}, nil
}

func (pr *ProductRepository) UpdateProduct(dto UpdateProductDto, slug *string, id int) exception.Error {
	if dto.Description != nil {
		_, err := pr.db.Exec(context.Background(), "UPDATE product SET description = $1 WHERE product_id = $2;", dto.Description, id)
		if err != nil {
			return exception.NewErr(productUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if dto.Title != nil {
		_, err := pr.db.Exec(context.Background(), "UPDATE product SET title = $1 WHERE product_id = $2;", dto.Title, id)
		if err != nil {
			return exception.NewErr(productUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if slug != nil {
		_, err := pr.db.Exec(context.Background(), "UPDATE product SET slug = $1 WHERE product_id = $2;", slug, id)
		if err != nil {
			return exception.NewErr(productUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	return nil
}

func (pr *ProductRepository) UpdateProductModel(dto UpdateProductModelDto, modelId int) exception.Error {

	if dto.Discount != nil {
		_, err := pr.db.Exec(context.Background(), "UPDATE product_model SET discount = $1 WHERE product_model_id = $2;", dto.Discount, modelId)
		if err != nil {
			return exception.NewErr(productUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if dto.Price != nil {
		_, err := pr.db.Exec(context.Background(), "UPDATE product_model SET price = $1 WHERE product_model_id = $2;", dto.Price, modelId)
		if err != nil {
			return exception.NewErr(productUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if dto.ImagePath != nil {
		_, err := pr.db.Exec(context.Background(), "UPDATE product_model SET main_image_path = $1 WHERE product_model_id = $2;",
			dto.ImagePath, modelId)
		if err != nil {
			return exception.NewErr(productUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}
	return nil
}

func (pr *ProductRepository) DeleteProduct(id int) exception.Error {
	query := `
	DELETE FROM product WHERE product_id = $1;
	`

	_, err := pr.db.Exec(context.Background(), query, id)

	if err != nil {
		return exception.ServerError(err.Error())
	}
	return nil
}

func (pr *ProductRepository) DeleteProductModel(id int) exception.Error {
	query := `
	DELETE FROM product_model WHERE product_model_id = $1;
	`

	_, err := pr.db.Exec(context.Background(), query, id)

	if err != nil {
		return exception.ServerError(err.Error())
	}
	return nil
}
