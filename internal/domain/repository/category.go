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

type CategoryRepository struct {
	db db.PostgresClient
}

func NewCategoryRepository(db db.PostgresClient) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetChildrenCount(ctx context.Context, id int) (*int, fall.Error) {
	q := "select count(*) from category where parent_category_id = $1;"

	row := r.db.QueryRow(ctx, q, id)

	var count int

	err := row.Scan(&count)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	return &count, nil
}

func (r *CategoryRepository) GetWithoutChildren(ctx context.Context) ([]model.CategoryModel, fall.Error) {
	q := `SELECT category_id, parent_category_id, slug, title, short_title, img_path
	FROM category c
	WHERE NOT EXISTS (
		SELECT 1
		FROM category
		WHERE parent_category_id = c.category_id
	);`

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	var cats []model.CategoryModel

	for rows.Next() {
		c := model.CategoryModel{}
		err := rows.Scan(&c.Id, &c.ParentId, &c.Slug, &c.Title, &c.ShortTitle, &c.ImgPath)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		cats = append(cats, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return cats, nil
}

func (r *CategoryRepository) GetTopLevels(ctx context.Context) ([]model.CategoryModel, fall.Error) {
	q := "SELECT category_id, parent_category_id, slug, title, short_title, img_path FROM category  WHERE parent_category_id IS NULL;"

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var cats []model.CategoryModel

	for rows.Next() {
		c := model.CategoryModel{}
		err := rows.Scan(&c.Id, &c.ParentId, &c.Slug, &c.Title, &c.ShortTitle, &c.ImgPath)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		cats = append(cats, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return cats, nil
}

func (r *CategoryRepository) Create(ctx context.Context, dto model.CreateCategoryDto, slug string) fall.Error {
	query := "INSERT INTO category (title, img_path, parent_category_id, slug, short_title) VALUES ($1, $2, $3, $4, $5);"

	_, err := r.db.Exec(ctx, query, dto.Title, dto.ImgPath, dto.ParentId, slug, dto.ShortTitle)

	if err != nil {
		return fall.NewErr(msg.CategoryCreateError, fall.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (r *CategoryRepository) FindByField(ctx context.Context, field string, value any) (*model.CategoryModel, fall.Error) {
	query := fmt.Sprintf(`
	SELECT category_id, parent_category_id, slug, title, short_title, img_path
	FROM category WHERE %s = $1;
	`, field)

	row := r.db.QueryRow(ctx, query, value)

	cat := model.CategoryModel{}

	err := row.Scan(&cat.Id, &cat.ParentId, &cat.Slug, &cat.Title, &cat.ShortTitle, &cat.ImgPath)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.CategoryNotFound, 404)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &cat, nil
}

func (r *CategoryRepository) FindByFieldRelation(ctx context.Context, field string, value any) (*model.Category, fall.Error) {
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

	rows, err := r.db.Query(ctx, query, value)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	result := model.Category{}
	secondMap := make(map[int]*model.Category)
	var secondOrder []int
	thirdMap := make(map[int]*model.Category)
	var thirdOrder []int
	fourthMap := make(map[int]*model.Category)
	var fourthOrder []int
	fiveMap := make(map[int]*model.Category)
	var fiveOrder []int
	founded := false

	for rows.Next() {

		category := model.Category{}
		category.Subcategories = []*model.Category{}

		err := rows.Scan(&category.Id, &category.Title, &category.Slug,
			&category.ShortTitle, &category.ImgPath, &category.ParentId, &category.Level)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		switch category.Level {
		case 1:
			result = category
		case 2:
			secondMap[category.Id] = &category
			secondOrder = append(secondOrder, category.Id)
		case 3:
			thirdMap[category.Id] = &category
			thirdOrder = append(thirdOrder, category.Id)
		case 4:
			fourthMap[category.Id] = &category
			fourthOrder = append(fourthOrder, category.Id)
		case 5:
			fiveMap[category.Id] = &category
			fiveOrder = append(fiveOrder, category.Id)
		}

		if !founded {
			founded = true
		}

	}
	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.CategoryNotFound, 404)
	}

	for _, id := range fiveOrder {
		p := fiveMap[id]
		fourth := fourthMap[*p.ParentId]
		fourth.Subcategories = append(fourth.Subcategories, p)
	}

	for _, id := range fourthOrder {
		p := fourthMap[id]
		third := thirdMap[*p.ParentId]
		third.Subcategories = append(third.Subcategories, p)
	}

	for _, id := range thirdOrder {
		p := thirdMap[id]
		second := secondMap[*p.ParentId]
		second.Subcategories = append(second.Subcategories, p)
	}

	for _, id := range secondOrder {
		p := secondMap[id]
		result.Subcategories = append(result.Subcategories, p)
	}

	return &result, nil
}

func (r *CategoryRepository) GetParentSubLevel(ctx context.Context, id int) (*model.CategoryModel, fall.Error) {
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
	row := r.db.QueryRow(ctx, query, id)

	c := model.CategoryModel{}

	err := row.Scan(&c.Id, &c.ParentId, &c.Slug,
		&c.Title, &c.ShortTitle, &c.ImgPath, nil)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.CategoryNotFound, fall.STATUS_NOT_FOUND)
		}

		return nil, fall.ServerError(err.Error())
	}

	return &c, nil
}

func (r *CategoryRepository) GetParentTopLevel(ctx context.Context, id int) (*model.CategoryModel, fall.Error) {
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
	row := r.db.QueryRow(ctx, query, id)

	c := model.CategoryModel{}

	err := row.Scan(&c.Id, &c.ParentId, &c.Slug, &c.Title, &c.ShortTitle, &c.ImgPath)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.CategoryNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &c, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]*model.Category, fall.Error) {
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
	SELECT * FROM category_tree ORDER BY category_id;
`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	firstMap := make(map[int]*model.Category)
	var firstOrder []int
	secondMap := make(map[int]*model.Category)
	var secondOrder []int
	thirdMap := make(map[int]*model.Category)
	var thirdOrder []int
	fourthMap := make(map[int]*model.Category)
	var fourthOrder []int
	fiveMap := make(map[int]*model.Category)
	var fiveOrder []int
	var founded bool = false

	for rows.Next() {

		category := model.Category{}
		category.Subcategories = []*model.Category{}

		err := rows.Scan(&category.Id, &category.Title, &category.Slug,
			&category.ShortTitle, &category.ImgPath, &category.ParentId, &category.Level)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		switch category.Level {
		case 1:
			firstMap[category.Id] = &category
			firstOrder = append(firstOrder, category.Id)
		case 2:
			secondMap[category.Id] = &category
			secondOrder = append(secondOrder, category.Id)
		case 3:
			thirdMap[category.Id] = &category
			thirdOrder = append(thirdOrder, category.Id)
		case 4:
			fourthMap[category.Id] = &category
			fourthOrder = append(fourthOrder, category.Id)
		case 5:
			fiveMap[category.Id] = &category
			fiveOrder = append(fiveOrder, category.Id)
		}

		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.CategoryNotFound, fall.STATUS_NOT_FOUND)
	}

	for _, id := range fiveOrder {
		p := fiveMap[id]
		fourth := fourthMap[*p.ParentId]
		fourth.Subcategories = append(fourth.Subcategories, p)
	}

	for _, id := range fourthOrder {
		p := fourthMap[id]
		third := thirdMap[*p.ParentId]
		third.Subcategories = append(third.Subcategories, p)
	}

	for _, id := range thirdOrder {
		p := thirdMap[id]
		second := secondMap[*p.ParentId]
		second.Subcategories = append(second.Subcategories, p)
	}

	for _, id := range secondOrder {
		p := secondMap[id]
		first := firstMap[*p.ParentId]
		first.Subcategories = append(first.Subcategories, p)
	}

	result := make([]*model.Category, 0, len(firstMap))

	for _, id := range firstOrder {
		first := firstMap[id]
		result = append(result, first)
	}

	return result, nil
}

func (r *CategoryRepository) Update(ctx context.Context, dto model.UpdateCategoryDto, newSlug *string, id int) fall.Error {

	var queries []string

	if dto.ImgPath != nil {
		queries = append(queries, fmt.Sprintf("img_path = '%s'", *dto.ImgPath))
	}

	if dto.ShortTitle != nil {
		queries = append(queries, fmt.Sprintf("short_title = '%s'", *dto.ShortTitle))
	}

	if dto.Title != nil {
		queries = append(queries, fmt.Sprintf("title = '%s'", *dto.Title))
	}

	if newSlug != nil {
		queries = append(queries, fmt.Sprintf("slug = '%s'", *newSlug))
	}

	if len(queries) > 0 {
		q := "UPDATE category SET " + strings.Join(queries, ",") + " WHERE category_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.CategoryUpdateError, err.Error()))
		}
		return nil
	}

	return nil
}

func (cr *CategoryRepository) GetCatalogCategories(ctx context.Context, id int, activeSlug string) (*model.CatalogCategoryResponse, fall.Error) {
	query := `
	WITH RECURSIVE category_tree AS (
		SELECT category_id, title, slug, short_title, img_path, parent_category_id, 1 AS level
		FROM category
		WHERE category_id = $1
		UNION ALL
		SELECT c.category_id, c.title, c.slug, c.short_title, c.img_path, c.parent_category_id, ct.level +1 as level
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT * FROM category_tree;
`
	rows, err := cr.db.Query(ctx, query, id)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	result := model.СatalogCategory{}
	current := model.CategoryModel{}
	secondMap := make(map[int]*model.СatalogCategory)
	var secondOrder []int
	thirdMap := make(map[int]*model.СatalogCategory)
	var thirdOrder []int
	fourthMap := make(map[int]*model.СatalogCategory)
	var fourthOrder []int
	fiveMap := make(map[int]*model.СatalogCategory)
	var fiveOrder []int
	founded := false

	for rows.Next() {

		category := model.СatalogCategory{}
		category.Subcategories = []*model.СatalogCategory{}

		err := rows.Scan(&category.Id, &category.Title, &category.Slug,
			&category.ShortTitle, &category.ImgPath, &category.ParentId, &category.Level)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		if category.Slug == activeSlug {
			category.Active = true
			current.Id = category.Id
			current.ImgPath = category.ImgPath
			current.ParentId = category.ParentId
			current.ShortTitle = category.ShortTitle
			current.Slug = category.Slug
			current.Title = category.Title
		} else {
			category.Active = false
		}

		switch category.Level {
		case 1:
			result = category
		case 2:
			secondMap[category.Id] = &category
			secondOrder = append(secondOrder, category.Id)
		case 3:
			thirdMap[category.Id] = &category
			thirdOrder = append(thirdOrder, category.Id)
		case 4:
			fourthMap[category.Id] = &category
			fourthOrder = append(fourthOrder, category.Id)
		case 5:
			fiveMap[category.Id] = &category
			fiveOrder = append(fiveOrder, category.Id)
		}

		if !founded {
			founded = true
		}

	}
	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.CategoryNotFound, fall.STATUS_NOT_FOUND)
	}

	for _, id := range fiveOrder {
		p := fiveMap[id]
		fourth := fourthMap[*p.ParentId]
		fourth.Subcategories = append(fourth.Subcategories, p)
		if p.Active && !fourth.Active {
			fourth.Active = true
		}
	}

	for _, id := range fourthOrder {
		p := fourthMap[id]
		third := thirdMap[*p.ParentId]
		third.Subcategories = append(third.Subcategories, p)
		if p.Active && !third.Active {
			third.Active = true
		}
	}

	for _, id := range thirdOrder {
		p := thirdMap[id]
		second := secondMap[*p.ParentId]
		second.Subcategories = append(second.Subcategories, p)
		if p.Active && !second.Active {
			second.Active = true
		}
	}

	for _, id := range secondOrder {
		p := secondMap[id]
		if p.Active && !result.Active {
			result.Active = true
		}
		result.Subcategories = append(result.Subcategories, p)
	}

	return &model.CatalogCategoryResponse{
		CatalogCategories: result,
		Current:           current,
	}, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, slug string) fall.Error {
	q := "DELETE FROM category WHERE slug = $1;"

	_, err := r.db.Exec(ctx, q, slug)

	if err != nil {
		return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.CategoryDeleteError, err.Error()))
	}

	return nil
}
