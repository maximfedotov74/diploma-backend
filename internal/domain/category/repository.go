package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

// todo add errors
type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cr *CategoryRepository) GetCatalogCategories() ([]CatalogCategory, error) {
	ctx := context.Background()

	query := `select t.title as third_lvl_title, t.category_id as third_lvl_id, t.parent_category_id as third_lvl_id_parent,
	s.title as second_lvl_title, s.category_id as second_lvl_id, s.parent_category_id as second_lvl_id_parent,
	f.title as first_lvl_title, f.category_id as first_lvl_id,
	count(*) as products_found from product as p 
	inner join category as t on t.category_id = p.category_id
	inner join category as s on s.category_id = t.parent_category_id
	inner join category as f on f.category_id = s.parent_category_id
	group by (third_lvl_title, second_lvl_title, first_lvl_title, third_lvl_id, second_lvl_id, first_lvl_id, third_lvl_id_parent, second_lvl_id_parent);`

	rows, err := cr.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	firstMap := make(map[int]CatalogCategory)
	secondMap := make(map[int]CatalogCategory)
	thirdMap := make(map[int]CatalogCategory)

	for rows.Next() {
		first := CatalogCategory{}
		second := CatalogCategory{}
		third := CatalogCategory{}
		err := rows.Scan(&third.Title, &third.Id, &third.ParentId, &second.Title, &second.Id, &second.ParentId,
			&first.Title, &first.Id, &third.Count)
		if err != nil {
			return nil, err
		}

		_, ok := firstMap[first.Id]
		if !ok {
			firstMap[first.Id] = first
		}

		_, ok = secondMap[second.Id]
		if !ok {
			secondMap[second.Id] = second
		}

		_, ok = thirdMap[third.Id]
		if !ok {
			thirdMap[third.Id] = third
		}

	}

	if rows.Err() != nil {
		return nil, err
	}

	for _, value := range thirdMap {
		second := secondMap[value.ParentId]
		second.Subcategories = append(second.Subcategories, value)
		second.Count += value.Count
		secondMap[value.ParentId] = second
	}

	for _, value := range secondMap {
		first := firstMap[value.ParentId]
		first.Subcategories = append(first.Subcategories, value)
		first.Count += value.Count
		firstMap[value.ParentId] = first
	}

	res := make([]CatalogCategory, 0, len(firstMap))

	for _, v := range firstMap {
		res = append(res, v)
	}

	return res, nil
}

func (cr *CategoryRepository) FindByFieldWithSubcategories(field string, value any) (*Category, error) {
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
		return nil, err
	}
	defer rows.Close()

	result := Category{}
	secondMap := make(map[int]Category)
	thirdMap := make(map[int]Category)
	fourthMap := make(map[int]Category)

	for rows.Next() {

		category := Category{}

		err := rows.Scan(&category.Id, &category.Title, &category.Slug,
			&category.ShortTitle, &category.ImgPath, &category.ParentId, &category.Level)

		if err != nil {
			return nil, err
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

	}

	if rows.Err() != nil {
		return nil, err
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

func (cr *CategoryRepository) FindByField(field string, value any) (*CategoryDb, error) {
	query := fmt.Sprintf(`
	SELECT category_id, parent_category_id, slug, title, short_title, img_path
	FROM category WHERE %s = $1;
	`, field)

	row := cr.db.QueryRow(context.Background(), query, value)

	cat := CategoryDb{}

	err := row.Scan(&cat.Id, &cat.ParentId, &cat.Slug, &cat.Title, &cat.ShortTitle, &cat.ImgPath)
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (cr *CategoryRepository) GetParentTopLevel(id int) (*CategoryDb, error) {
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
		return nil, err
	}

	return &categoryDb, nil
}

func (cr *CategoryRepository) GetParentSubLevel(id int) (*CategoryDb, error) {
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
		return nil, err
	}

	return &categoryDb, nil
}

func (cr *CategoryRepository) CreateCategory(dto CreateCategoryDto, slug string) error {

	ctx := context.Background()

	query := "INSERT INTO category (title, img_path, parent_category_id, slug, short_title) VALUES ($1, $2, $3, $4, $5);"

	_, err := cr.db.Exec(ctx, query, dto.Title, dto.ImgPath, dto.ParentId, slug, dto.ShortTitle)

	if err != nil {
		return errors.New(messages.CATEGORY_CREATE_ERROR)
	}

	return nil
}
