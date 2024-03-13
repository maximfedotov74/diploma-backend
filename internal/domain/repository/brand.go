package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type BrandRepository struct {
	db db.PostgresClient
}

func NewBrandRepository(db db.PostgresClient) *BrandRepository {
	return &BrandRepository{db: db}
}

func (r *BrandRepository) GetBrandsByGender(ctx context.Context, slug string) ([]model.Brand, fall.Error) {
	q := `
	WITH RECURSIVE category_tree AS (
		SELECT category_id, slug, parent_category_id
		FROM category
		WHERE slug = $1
		UNION ALL
		SELECT c.category_id, c.slug, c.parent_category_id
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT DISTINCT b.brand_id, b.title, b.slug, b.description, b.img_path
	FROM product p
	INNER JOIN category_tree ct ON p.category_id = ct.category_id 
	INNER JOIN brand b on p.brand_id = b.brand_id;
	`

	rows, err := r.db.Query(ctx, q, slug)
	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	var result []model.Brand
	for rows.Next() {
		brand := model.Brand{}
		err := rows.Scan(&brand.Id, &brand.Title, &brand.Slug, &brand.Description, &brand.ImgPath)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		result = append(result, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return result, nil
}

func (r *BrandRepository) CreateBrand(ctx context.Context, dto model.CreateBrandDto, slug string) fall.Error {
	query := "INSERT INTO brand (title, slug, description, img_path) VALUES ($1, $2, $3, $4);"

	_, err := r.db.Exec(ctx, query, dto.Title, slug, dto.Description, dto.ImgPath)
	if err != nil {
		return fall.NewErr(msg.BrandCreateError, fall.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (r *BrandRepository) FindByFeild(ctx context.Context, field string, value any) (*model.Brand, fall.Error) {
	query := fmt.Sprintf("SELECT brand_id, title, slug, description, img_path FROM brand WHERE %s = $1;", field)
	row := r.db.QueryRow(ctx, query, value)

	brand := model.Brand{}

	err := row.Scan(&brand.Id, &brand.Title, &brand.Slug, &brand.Description, &brand.ImgPath)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.BrandNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}
	return &brand, nil
}

func (r *BrandRepository) GetAll(ctx context.Context) ([]model.Brand, fall.Error) {
	query := `
	SELECT brand_id, title, slug, description, img_path FROM brand
	ORDER BY brand_id;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	var result []model.Brand
	var founded bool = false
	for rows.Next() {
		brand := model.Brand{}
		err := rows.Scan(&brand.Id, &brand.Title, &brand.Slug, &brand.Description, &brand.ImgPath)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		result = append(result, brand)
		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.BrandNotFound, fall.STATUS_NOT_FOUND)
	}

	return result, nil
}

func (r *BrandRepository) UpdateBrand(ctx context.Context, dto model.UpdateBrandDto, newSlug *string, id int) fall.Error {

	var queries []string

	if dto.ImgPath != nil {
		queries = append(queries, fmt.Sprintf("img_path = '%s'", *dto.ImgPath))
	}

	if dto.Title != nil {
		queries = append(queries, fmt.Sprintf("title = '%s'", *dto.Title))
	}

	if dto.Description != nil {
		queries = append(queries, fmt.Sprintf("description = '%s'", *dto.Description))
	}

	if newSlug != nil {
		queries = append(queries, fmt.Sprintf("slug = '%s'", *newSlug))
	}

	if len(queries) > 0 {
		q := "UPDATE brand SET " + strings.Join(queries, ",") + " WHERE brand_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.BrandUpdateError, err.Error()))
		}
		return nil
	}

	return nil
}

func (r *BrandRepository) Delete(ctx context.Context, slug string) fall.Error {
	q := "DELETE FROM brand WHERE slug = $1;"

	_, err := r.db.Exec(ctx, q, slug)

	if err != nil {
		return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.BrandDeleteError, err.Error()))
	}

	return nil
}
