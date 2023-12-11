package category

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cr *CategoryRepository) FindByFieldWithSubcategories(field string, value any) (*Category, exception.Error) {
	query := fmt.Sprintf(`
	WITH RECURSIVE category_tree AS (
		SELECT category_id, title, slug, short_title, img_path, parent_category_id, 1 AS level
		FROM category
		WHERE %s = $1
	
		UNION ALL
	
		SELECT c.category_id, c.title, c.slug, c.short_title, c.img_path, c.parent_category_id, ct.level +1 as level
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT * FROM category_tree;
`, field)

	ctx := context.Background()

	rows, err := cr.db.Query(ctx, query, value)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	defer rows.Close()

	result := Category{}
	secondMap := make(map[int]Category)
	thirdMap := make(map[int]Category)
	fourthMap := make(map[int]Category)
	founded := false

	for rows.Next() {

		category := Category{}

		err := rows.Scan(&category.Id, &category.Title, &category.Slug,
			&category.ShortTitle, &category.ImgPath, &category.ParentId, &category.Level)

		if err != nil {
			return nil, exception.ServerError(err.Error())
		}

		if category.Level == 1 {
			result = category
		} else {
			switch category.Level {
			case 2:
				secondMap[category.Id] = category
				break
			case 3:
				thirdMap[category.Id] = category
				break
			case 4:
				fourthMap[category.Id] = category
			}
		}
		if !founded {
			founded = true
		}

	}
	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, exception.NewErr(categoryNotFound, 404)
	}

	for _, value := range fourthMap {
		third := thirdMap[*value.ParentId]
		third.Subcategories = append(third.Subcategories, value)
		thirdMap[*value.ParentId] = third
	}

	for _, value := range thirdMap {
		second := secondMap[*value.ParentId]
		second.Subcategories = append(second.Subcategories, value)
		secondMap[*value.ParentId] = second
	}

	for _, value := range secondMap {
		result.Subcategories = append(result.Subcategories, value)
	}

	return &result, nil
}

func (cr *CategoryRepository) FindByField(field string, value any) (*CategoryDb, exception.Error) {
	query := fmt.Sprintf(`
	SELECT category_id, parent_category_id, slug, title, short_title, img_path
	FROM category WHERE %s = $1;
	`, field)

	row := cr.db.QueryRow(context.Background(), query, value)

	cat := CategoryDb{}

	err := row.Scan(&cat.Id, &cat.ParentId, &cat.Slug, &cat.Title, &cat.ShortTitle, &cat.ImgPath)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, exception.NewErr(categoryNotFound, 404)
		}
		return nil, exception.ServerError(err.Error())
	}

	return &cat, nil
}

func (cr *CategoryRepository) GetParentTopLevel(id int) (*CategoryDb, exception.Error) {
	query := `
	WITH RECURSIVE recursive_cte AS (
		SELECT category_id, parent_category_id, slug, title, short_title, img_path
		FROM category
		WHERE category_id = $1
		UNION ALL
		SELECT t.category_id, t.parent_category_id, t.slug, t.title, t.short_title, t.img_path
		FROM category t
		INNER JOIN recursive_cte r ON r.parent_category_id = t.category_id
	)
	SELECT *
	FROM recursive_cte
	WHERE parent_category_id IS NULL;
	`
	row := cr.db.QueryRow(context.Background(), query, id)

	categoryDb := CategoryDb{}

	err := row.Scan(&categoryDb.Id, &categoryDb.ParentId, &categoryDb.Slug, &categoryDb.Title, &categoryDb.ShortTitle, &categoryDb.ImgPath)
	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	return &categoryDb, nil
}

func (cr *CategoryRepository) GetParentSubLevel(id int) (*CategoryDb, exception.Error) {
	query := `
	WITH RECURSIVE recursive_cte AS (
		SELECT category_id, parent_category_id, slug, title, short_title, img_path, 1 AS level
		FROM category
		WHERE category_id = $1
		UNION ALL
		SELECT t.category_id, t.parent_category_id, t.slug, t.title, t.short_title, t.img_path, r.level +1
		FROM category t
		INNER JOIN recursive_cte r ON r.parent_category_id = t.category_id
		WHERE r.level < 3
	)
	SELECT *
	FROM recursive_cte
	WHERE level = 3;
	`
	row := cr.db.QueryRow(context.Background(), query, id)

	categoryDb := CategoryDb{}

	err := row.Scan(&categoryDb.Id, &categoryDb.ParentId, &categoryDb.Slug,
		&categoryDb.Title, &categoryDb.ShortTitle, &categoryDb.ImgPath, nil)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	return &categoryDb, nil
}

func (cr *CategoryRepository) CreateCategory(dto CreateCategoryDto, slug string) exception.Error {

	ctx := context.Background()

	query := "INSERT INTO category (title, img_path, parent_category_id, slug, short_title) VALUES ($1, $2, $3, $4, $5);"

	_, err := cr.db.Exec(ctx, query, dto.Title, dto.ImgPath, dto.ParentId, slug, dto.ShortTitle)

	if err != nil {
		return exception.NewErr(categoryCreateError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (cr *CategoryRepository) GetAll() ([]Category, exception.Error) {
	query := `
	WITH RECURSIVE category_tree AS (
		SELECT category_id, title, slug, short_title, img_path, parent_category_id, 1 AS level
		FROM category	
		WHERE parent_category_id is NULL
		UNION ALL
		SELECT c.category_id, c.title, c.slug, c.short_title, c.img_path, c.parent_category_id, ct.level +1 as level
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT * FROM category_tree;
`

	ctx := context.Background()

	rows, err := cr.db.Query(ctx, query)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	defer rows.Close()

	firstMap := make(map[int]Category)
	secondMap := make(map[int]Category)
	thirdMap := make(map[int]Category)
	fourthMap := make(map[int]Category)

	var founded bool = false

	for rows.Next() {

		category := Category{}

		err := rows.Scan(&category.Id, &category.Title, &category.Slug,
			&category.ShortTitle, &category.ImgPath, &category.ParentId, &category.Level)

		if err != nil {
			return nil, exception.ServerError(err.Error())
		}

		switch category.Level {
		case 1:
			firstMap[category.Id] = category
		case 2:
			secondMap[category.Id] = category
			break
		case 3:
			thirdMap[category.Id] = category
			break
		case 4:
			fourthMap[category.Id] = category
		}

		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, exception.NewErr(categoryNotFound, exception.STATUS_NOT_FOUND)
	}

	for _, value := range fourthMap {
		third := thirdMap[*value.ParentId]
		third.Subcategories = append(third.Subcategories, value)
		thirdMap[*value.ParentId] = third
	}

	for _, value := range thirdMap {
		second := secondMap[*value.ParentId]
		second.Subcategories = append(second.Subcategories, value)
		secondMap[*value.ParentId] = second
	}

	for _, value := range secondMap {
		first := firstMap[*value.ParentId]
		first.Subcategories = append(first.Subcategories, value)
		firstMap[first.Id] = first
	}

	result := make([]Category, 0, len(firstMap))

	for _, v := range firstMap {
		result = append(result, v)
	}

	return result, nil
}

func (cr *CategoryRepository) UpdateCategory(dto UpdateCategoryDto, newSlug *string, id int) exception.Error {

	if dto.ImgPath != nil {
		_, err := cr.db.Exec(context.Background(), "UPDATE category SET img_path = $1 WHERE category_id = $2;", dto.ImgPath, id)
		if err != nil {
			return exception.NewErr(categoryUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if dto.ShortTitle != nil {
		_, err := cr.db.Exec(context.Background(), "UPDATE category SET short_title = $1 WHERE category_id = $2;", dto.ShortTitle, id)
		if err != nil {
			return exception.NewErr(categoryUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if dto.Title != nil {
		_, err := cr.db.Exec(context.Background(), "UPDATE category SET title = $1 WHERE category_id = $2;", dto.Title, id)
		if err != nil {
			return exception.NewErr(categoryUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}

	if newSlug != nil {
		_, err := cr.db.Exec(context.Background(), "UPDATE category SET slug = $1 WHERE category_id = $2;", newSlug, id)
		if err != nil {
			return exception.NewErr(categoryUpdateError, exception.STATUS_INTERNAL_ERROR)
		}
	}
	return nil
}
