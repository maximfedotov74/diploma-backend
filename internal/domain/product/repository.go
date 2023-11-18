package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// todo add errors
// todo add in_stock
type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) CreateProduct(dto CreateProductDto, slug string) error {
	query := `
  INSERT INTO product (title, slug, description, category_id, brand_id)
  VALUES ($1, $2, $3, $4, $5);
  `

	_, err := pr.db.Exec(context.Background(), query, dto.Title, slug, dto.Description, dto.CategoryID, dto.BrandID)
	if err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) FindById(id int) (*ProductWithoutRelations, error) {
	query := `
	select product_id, title, slug, description from product where product_id = $1; 	
	`
	row := pr.db.QueryRow(context.Background(), query, id)

	p := ProductWithoutRelations{}

	err := row.Scan(&p.Id, &p.Title, &p.Slug, &p.Description)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (pr *ProductRepository) GetCatalogModels(categorySlug string, queryParams ...string) error {
	return nil
}

func (pr *ProductRepository) FindByProductSlugAndModelId(slug string, modelId int) (*Product, error) {
	query := `
	select p.product_id as p_id, p.title as p_title, p.slug as p_slug,
	p.description as p_description,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug, c.short_title as c_short_title,
	c.img_path as c_img_path,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, b.img_path as b_img_path,
	pm.product_model_id as pm_id, pm.price as pm_price, pm.discount as pm_discount, pm.product_id as pm_product_id,
	pimg.product_img_id as pimg_id, pimg.img_path as pimg_img_path, pimg.main as pimg_main, pimg.product_model_id as pimg_model_id,
	op.option_id as op_id, op.title as op_title, op.slug as op_slug, pmop.product_model_id as pmop_model_id,
	v.option_value_id as v_id, v.value as v_value, v.info as v_info, v.option_id as v_option_id, pmop.product_model_id as pmop_model_id_v
	from product as p
	inner join category as c on p.category_id = c.category_id
	inner join brand as b on p.brand_id = b.brand_id
	left join product_model as pm on pm.product_id = p.product_id
	left join product_model_img as pimg on pimg.product_model_id = pm.product_model_id
	left join product_model_option as pmop on pmop.product_model_id = pm.product_model_id
	left join option as op on op.option_id = pmop.option_id
	left join option_value as v on v.option_value_id = pmop.option_value_id
	where p.slug = $1;
	`

	rows, err := pr.db.Query(context.Background(), query, slug)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	p := Product{}
	processedRows := 0

	imagesMap := make(map[int]ProductModelImg)
	modelsMap := make(map[int]ProductModel)
	optionsMap := make(map[int]ProductModelOption)
	valuesMap := make(map[int]ProductModelOptionValue)

	for rows.Next() {

		productModel := ProductModel{}
		productModelImg := ProductModelImg{}
		opt := ProductModelOption{}
		val := ProductModelOptionValue{}

		err := rows.Scan(&p.Id, &p.Title, &p.Slug, &p.Description,
			&p.Category.Id, &p.Category.Title, &p.Category.Slug, &p.Category.ShortTitle, &p.Category.ImgPath,
			&p.Brand.Id, &p.Brand.Title, &p.Brand.Slug, &p.Brand.ImgPath,
			&productModel.Id, &productModel.Price, &productModel.Discount, &productModel.ProductId,
			&productModelImg.Id, &productModelImg.ImgPath, &productModelImg.Main, &productModelImg.ProductModelId,
			&opt.Id, &opt.Title, &opt.Slug, &opt.ProductModelId,
			&val.Id, &val.Value, &val.Info, &val.OptionId, &val.ProductModelId,
		)

		if err != nil {
			return nil, err
		}

		if productModel.Id != nil {
			modelsMap[*productModel.Id] = productModel
		}

		if productModelImg.Id != nil {
			imagesMap[*productModelImg.Id] = productModelImg
		}

		if opt.Id != nil {
			if *opt.ProductModelId == modelId {
				optionsMap[*opt.ProductModelId] = opt
			}
		}

		if val.Id != nil {
			if *val.ProductModelId == modelId {
				valuesMap[*val.ProductModelId] = val
			}
		}

		processedRows++

	}

	if rows.Err() != nil {
		return nil, err
	}

	if processedRows == 0 {
		return nil, errors.New("Product Not found")
	}

	for key, v := range valuesMap {
		opt := optionsMap[key]
		opt.Value = v
		optionsMap[key] = opt
	}

	for key, v := range optionsMap {
		m := modelsMap[key]
		m.Options = append(m.Options, v)
		modelsMap[*m.Id] = m
	}

	for _, v := range imagesMap {
		m := modelsMap[*v.ProductModelId]
		m.Images = append(m.Images, v)
		modelsMap[*m.Id] = m
	}

	models := make([]ProductModel, 0, len(modelsMap))

	current := ProductModel{}

	for _, v := range modelsMap {

		if *v.Id == modelId {
			current = v
		} else {
			models = append(models, v)
		}
	}

	p.CurrentModel = &current
	p.Models = models

	return &p, nil
}

func (pr *ProductRepository) CreateModel(dto CreateProductModelDto) error {
	query := "INSERT INTO product_model (price, discount, product_id) VALUES ($1, $2, $3);"

	_, err := pr.db.Exec(context.Background(), query, dto.Price, dto.Discount, dto.ProductId)

	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) AddPhoto(dto CreateProducModelImg) error {
	query := "INSERT INTO product_model_img (img_path, main, product_model_id) VALUES ($1,$2,$3);"

	_, err := pr.db.Exec(context.Background(), query, dto.ImgPath, dto.Main, dto.ProductModelId)
	if err != nil {
		return err
	}

	return nil
}
