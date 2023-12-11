package brand

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type BrandRepository struct {
	db *pgxpool.Pool
}

func NewBrandRepository(db *pgxpool.Pool) *BrandRepository {
	return &BrandRepository{db: db}
}

func (br *BrandRepository) CreateBrand(title string, slug string, description *string, imgPath *string) exception.Error {
	query := "INSERT INTO brand (title, slug, description, img_path) VALUES ($1, $2, $3, $4);"

	_, err := br.db.Exec(context.Background(), query, title, slug, description, imgPath)
	if err != nil {
		return exception.NewErr(brandCreateError, exception.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (br *BrandRepository) FindByFeild(field string, value any) (*Brand, exception.Error) {
	query := fmt.Sprintf("SELECT brand_id, title, slug, description, img_path FROM brand WHERE %s = $1;", field)
	row := br.db.QueryRow(context.Background(), query, value)

	brand := Brand{}

	err := row.Scan(&brand.Id, &brand.Title, &brand.Slug, &brand.Description, &brand.ImgPath)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, exception.NewErr(brandNotFound, exception.STATUS_NOT_FOUND)
		}
		return nil, exception.ServerError(err.Error())
	}
	return &brand, nil
}

func (br *BrandRepository) GetAll() ([]Brand, exception.Error) {
	query := `
	SELECT brand_id, title, slug, description, img_path FROM brand
	`

	rows, err := br.db.Query(context.Background(), query)
	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	defer rows.Close()

	var result []Brand
	var founded bool = false
	for rows.Next() {
		brand := Brand{}
		err := rows.Scan(&brand.Id, &brand.Title, &brand.Slug, &brand.Description, &brand.ImgPath)
		if err != nil {
			return nil, exception.ServerError(err.Error())
		}

		result = append(result, brand)
		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(err.Error())
	}

	if !founded {
		return nil, exception.NewErr(brandNotFound, exception.STATUS_NOT_FOUND)
	}

	return result, nil
}

func (cr *BrandRepository) UpdateBrand(dto UpdateBrandDto, newSlug *string, id int) exception.Error {

	if dto.ImgPath != nil {
		_, err := cr.db.Exec(context.Background(), "UPDATE brand SET img_path = $1 WHERE brand_id = $2", dto.ImgPath, id)
		if err != nil {
			return exception.NewErr(brandUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if dto.Title != nil {
		_, err := cr.db.Exec(context.Background(), "UPDATE brand SET description = $1 WHERE brand_id = $2", dto.Description, id)
		if err != nil {
			return exception.NewErr(brandUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if dto.Description != nil {
		_, err := cr.db.Exec(context.Background(), "UPDATE brand SET title = $1 WHERE brand_id = $2", dto.Title, id)
		if err != nil {
			return exception.NewErr(brandUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if newSlug != nil {
		_, err := cr.db.Exec(context.Background(), "UPDATE brand SET slug = $1 WHERE brand_id = $2", newSlug, id)
		if err != nil {
			return exception.NewErr(brandUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}
	return nil
}
