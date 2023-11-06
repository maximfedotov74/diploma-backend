package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cr *CategoryRepository) FindTypeByTitle(title string) (*model.CategoryType, error) {
	ctx := context.Background()

	query := "SELECT category_type_id, title FROM category_type WHERE title = $1;"

	row := cr.db.QueryRow(ctx, query, title)

	categoryType := model.CategoryType{}

	err := row.Scan(&categoryType.Id, &categoryType.Title)

	if err != nil {
		return nil, errors.New(messages.CATEGORY_NOT_FOUND)
	}
	return &categoryType, nil
}

func (cr *CategoryRepository) FindCategoryByTitle(title string) (*model.Category, error) {
	ctx := context.Background()

	query := `
	select parent.category_id as parent_id, parent.title as parent_title, parent.img_path as parent_img_path,
	parent.parent_category_id as parent_parent_id,
	child.category_id as child_id, child.title as child_title, child.img_path as child_img_path,
	child.parent_category_id as child_parent_id
	from category as parent
	left join category as child on parent.category_id = child.parent_category_id where parent.title = $1;
	`

	rows, err := cr.db.Query(ctx, query, title)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parent := model.Category{}
	processedRows := 0

	for rows.Next() {
		child := model.Category{}
		err := rows.Scan(&parent.Id, &parent.Title, &parent.ImgPath, &parent.ParentId,
			&child.Id, &child.Title, &child.ImgPath, &child.ParentId)

		if err != nil {
			return nil, err
		}
		parent.Subcategories = append(parent.Subcategories, child)
		processedRows++
	}

	if rows.Err() != nil {
		return nil, err
	}

	if processedRows == 0 {
		return nil, errors.New(messages.CATEGORY_NOT_FOUND)
	}

	return &parent, nil
}

func (cr *CategoryRepository) CreateCategoryType(title string) error {
	ctx := context.Background()

	query := "INSERT INTO category_type (title) VALUES ($1);"

	_, err := cr.db.Exec(ctx, query, title)

	if err != nil {
		return errors.New(messages.CATEGORY_CREATE_ERROR)
	}

	return nil
}

func (cr *CategoryRepository) CreateCategory(title string, img *string, parentId *int) error {
	ctx := context.Background()

	query := "INSERT INTO category (title, img_path, parent_category_id) VALUES ($1, $2, $3);"

	_, err := cr.db.Exec(ctx, query, title, img, parentId)

	if err != nil {
		return errors.New(messages.CATEGORY_CREATE_ERROR)
	}

	return nil
}
