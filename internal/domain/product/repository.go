package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (pr *ProductRepository) FindBySlug(slug string) (*Product, error) {
	query := `
	select p.product_id as p_id, p.title as p_title, p.slug as p_slug,
	p.description as p_description,
	c.category_id as c_id, c.title as c_title, c.slug as c_slug, c.short_title as c_short_title,
	c.img_path as c_img_path,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, b.img_path as b_img_path,
	pimg.product_img_id as pimg_id, pimg.img_path as pimg_img_path, pimg.main as pimg_main
	from product as p
	inner join category as c on p.category_id = c.category_id
	inner join brand as b on p.brand_id = b.brand_id
	left join product_img as pimg on pimg.product_id = p.product_id
	where p.slug = $1;
	`

	rows, err := pr.db.Query(context.Background(), query, slug)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	p := Product{}
	processedRows := 0

	for rows.Next() {

		productImg := ProductImg{}
		err := rows.Scan(&p.Id, &p.Title, &p.Slug, &p.Description,
			&p.Category.Id, &p.Category.Title, &p.Category.Slug, &p.Category.ShortTitle, &p.Category.ImgPath,
			&p.Brand.Id, &p.Brand.Title, &p.Brand.Slug, &p.Brand.ImgPath, &productImg.Id, &productImg.ImgPath, &productImg.Main,
		)

		if err != nil {
			return nil, err
		}

		if productImg.Id != nil {
			p.Images = append(p.Images, productImg)
		}
		processedRows++

	}

	if rows.Err() != nil {
		return nil, err
	}

	if processedRows == 0 {
		return nil, errors.New("Product Not found")
	}

	return &p, nil
}

func (pr *ProductRepository) AddPhoto(dto CreateProductImg) error {
	query := "INSERT INTO product_img (img_path, main, product_id) VALUES ($1,$2,$3);"

	_, err := pr.db.Exec(context.Background(), query, dto.ImgPath, dto.Main, dto.ProductId)
	if err != nil {
		return err
	}

	return nil
}
