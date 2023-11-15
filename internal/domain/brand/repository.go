package brand

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO add erros
type BrandRepository struct {
	db *pgxpool.Pool
}

func NewBrandRepository(db *pgxpool.Pool) *BrandRepository {
	return &BrandRepository{db: db}
}

func (br *BrandRepository) CreateBrand(title string, slug string, description *string, imgPath *string) error {
	query := "INSERT INTO brand (title, slug, description, img_path) VALUES ($1, $2, $3, $4);"

	_, err := br.db.Exec(context.Background(), query, title, slug, description, imgPath)
	if err != nil {
		return err
	}
	return nil
}

func (br *BrandRepository) FindByFeild(field string, value any) (*Brand, error) {
	query := fmt.Sprintf("SELECT brand_id, title, slug, description, img_path FROM brand WHERE %s = $1;", field)
	row := br.db.QueryRow(context.Background(), query, value)

	brand := Brand{}

	err := row.Scan(&brand.Id, &brand.Title, &brand.Slug, &brand.Description, &brand.ImgPath)
	if err != nil {
		return nil, err
	}
	return &brand, nil
}
