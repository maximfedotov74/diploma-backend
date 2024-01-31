package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type OptionRepository struct {
	db db.PostgresClient
}

func NewOptionRepository(db db.PostgresClient) *OptionRepository {
	return &OptionRepository{db: db}
}

func (r *OptionRepository) CheckValueInOption(ctx context.Context, valueId int, optionId int) fall.Error {
	q := "SELECT option_value_id FROM option_value WHERE option_value_id = $1 AND option_id = $2;"

	row := r.db.QueryRow(ctx, q, valueId, optionId)

	var id string

	err := row.Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fall.NewErr(msg.ValueNotInOption, fall.STATUS_BAD_REQUEST)
		}
		return fall.ServerError(err.Error())
	}

	return nil
}

func (r *OptionRepository) CreateOption(ctx context.Context, dto model.CreateOptionDto) fall.Error {
	query := `
  INSERT INTO option (title, slug) VALUES ($1,$2);
  `

	_, err := r.db.Exec(ctx, query, dto.Title, dto.Slug)

	if err != nil {
		return fall.ServerError(msg.OptionCreateError)
	}

	return nil
}

func (r *OptionRepository) CreateValue(ctx context.Context, dto model.CreateOptionValueDto) fall.Error {
	query := `
  INSERT INTO option_value (value, info, option_id) VALUES ($1,$2,$3);
  `
	_, err := r.db.Exec(ctx, query, dto.Value, dto.Info, dto.OptionId)

	if err != nil {
		return fall.ServerError(msg.OptionValueCreateError)
	}

	return nil
}

func (r *OptionRepository) DeleteOption(ctx context.Context, id int) fall.Error {
	query := `
  DELETE FROM option WHERE option_id = $1;
  `
	_, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fall.ServerError(msg.OptionDeleteError)
	}
	return nil
}

func (r *OptionRepository) DeleteValue(ctx context.Context, id int) fall.Error {
	query := `
  DELETE FROM option_value WHERE option_value_id = $1;
  `
	_, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fall.ServerError(msg.ValueDeleteError)
	}

	return nil
}

func (r *OptionRepository) DeleteSize(ctx context.Context, id int) fall.Error {
	q := "DELETE FROM sizes WHERE size_id = $1;"

	_, err := r.db.Exec(ctx, q, id)

	if err != nil {
		return fall.ServerError(msg.SizeDeleteError)
	}
	return nil
}

func (r *OptionRepository) DeleteSizeFromProductModel(ctx context.Context, modelSizeId int) fall.Error {
	q := "DELETE FROM model_sizes WHERE model_size_id = $1;"

	_, err := r.db.Exec(ctx, q, modelSizeId)
	if err != nil {
		return fall.ServerError(msg.ModelSizeDeleteError)
	}
	return nil
}

func (r *OptionRepository) DeleteOptionFromProductModel(ctx context.Context, productModelOptionId int) fall.Error {
	q := "DELETE FROM product_model_option WHERE product_model_option_id = $1;"

	_, err := r.db.Exec(ctx, q, productModelOptionId)
	if err != nil {
		return fall.ServerError(msg.OptionModelDeleteError)
	}
	return nil
}

func (r *OptionRepository) GetAll(ctx context.Context) ([]model.Option, fall.Error) {
	query := `
	SELECT op.option_id as op_id, op.title as op_title, op.slug as op_slug, op.for_catalog,
  v.option_value_id as v_id, v.value as v_value, v.info as v_info, v.option_id as v_option_id
  FROM option as op
  LEFT JOIN option_value as v ON v.option_id = op.option_id
	`
	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	optionsMap := make(map[int]model.Option)
	valuesMap := make(map[int]model.OptionValue)
	var optionsOrder []int
	var valuesOrder []int

	var founded bool = false

	for rows.Next() {
		opt := model.Option{}
		v := model.OptionValue{}

		err := rows.Scan(&opt.Id, &opt.Title, &opt.Slug, &opt.ForCatalog, &v.Id, &v.Value, &v.Info, &v.OptionId)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		if v.Id != nil {
			_, ok := valuesMap[*v.Id]
			if !ok {
				valuesMap[*v.Id] = v
				valuesOrder = append(valuesOrder, *v.Id)
			}
		}
		_, ok := optionsMap[opt.Id]
		if !ok {
			optionsMap[opt.Id] = opt
			optionsOrder = append(optionsOrder, opt.Id)
		}
		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	result := make([]model.Option, 0, len(optionsMap))

	if !founded {
		return result, nil
	}

	for _, v := range valuesOrder {
		val := valuesMap[v]
		opt := optionsMap[*val.OptionId]
		opt.Values = append(opt.Values, val)
		optionsMap[opt.Id] = opt
	}

	for _, v := range optionsOrder {
		opt := optionsMap[v]
		result = append(result, opt)
	}

	return result, nil

}

func (r *OptionRepository) FindByField(ctx context.Context, field string, value any) (*model.Option, fall.Error) {
	query := fmt.Sprintf(`
  SELECT op.option_id as op_id, op.title as op_title, op.slug as op_slug, op.for_catalog,
  v.option_value_id as v_id, v.value as v_value, v.info as v_info, v.option_id as v_option_id
  FROM option as op
  LEFT JOIN option_value as v ON v.option_id = op.option_id
  WHERE op.%s = $1;
  `, field)

	rows, err := r.db.Query(ctx, query, value)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	opt := model.Option{}

	valMap := make(map[int]model.OptionValue)
	var valOrder []int

	founded := false

	for rows.Next() {
		v := model.OptionValue{}

		err := rows.Scan(&opt.Id, &opt.Title, &opt.Slug, &opt.ForCatalog, &v.Id, &v.Value, &v.Info, &v.OptionId)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		if v.Id != nil {
			valMap[*v.Id] = v
			valOrder = append(valOrder, *v.Id)
		}
		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, fall.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.OptionNotFound, fall.STATUS_NOT_FOUND)
	}

	values := make([]model.OptionValue, 0, len(valMap))

	for _, key := range valOrder {
		v := valMap[key]
		values = append(values, v)
	}

	opt.Values = values

	return &opt, nil
}

func (r *OptionRepository) AddOptionToProductModel(ctx context.Context, dto model.AddOptionToProductModelDto) fall.Error {
	query := `
	INSERT INTO product_model_option (product_model_id, option_id, option_value_id) VALUES ($1,$2,$3);
	`
	_, err := r.db.Exec(ctx, query, dto.ProductModelId, dto.OptionId, dto.ValueId)
	if err != nil {
		return fall.ServerError(msg.AddOptionToProductError)
	}

	return nil
}

func (r *OptionRepository) AddSizeToProductModel(ctx context.Context, dto model.AddSizeToProductModelDto) fall.Error {
	query := `INSERT INTO model_sizes (product_model_id, size_id, literal_size, in_stock) VALUES ($1,$2,$3,$4);`

	_, err := r.db.Exec(ctx, query, dto.ProductModelId, dto.SizeId, dto.Literal, dto.InStock)
	if err != nil {
		log.Println(err.Error())
		return fall.ServerError(msg.AddSizeToProductError)
	}

	return nil
}

func (r *OptionRepository) CreateSize(ctx context.Context, value string) fall.Error {
	query := `INSERT INTO sizes (size_value) VALUES ($1);`

	_, err := r.db.Exec(ctx, query, value)
	if err != nil {
		return fall.ServerError(msg.SizeCreateError)
	}

	return nil
}

func (r *OptionRepository) GetCatalogFilters(ctx context.Context, categorySlug string) (*model.CatalogFilters, fall.Error) {

	mainQuery := `
	FROM product as p
	inner join category_tree as ct on p.category_id = ct.category_id
	inner join brand as b on p.brand_id = b.brand_id
	inner join product_model as pm on pm.product_id = p.product_id
	inner join product_model_option as pmop on pmop.product_model_id = pm.product_model_id
	inner join option as op on op.option_id = pmop.option_id
	inner join option_value as v on v.option_value_id = pmop.option_value_id
	inner join model_sizes as ms on ms.product_model_id = pm.product_model_id
	inner join sizes as sz on ms.size_id = sz.size_id
	`
	query := fmt.Sprintf(`
	WITH RECURSIVE category_tree AS (
		SELECT category_id, title, slug, parent_category_id
		FROM category
		WHERE slug = $1
		UNION ALL
		SELECT c.category_id, c.title, c.slug, c.parent_category_id
		FROM category c
		INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
	)
	SELECT op.option_id as option_id, op.title, op.slug, v.option_value_id as value_id, v.value as option_value,
	v.option_id as value_option_id,
	sz.size_value as size_value, sz.size_id as size_id,
	b.brand_id as b_id, b.title as b_title,
	(select max(pm.price)
	%[1]s
 	) as max_price,
	(select min(pm.price)
	%[1]s
	) as min_price
	%[1]s
	where op.for_catalog = true;
	`, mainQuery)

	rows, err := r.db.Query(ctx, query, categorySlug)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	optionsMap := make(map[int]model.CatalogOption)
	var optionsOrder []int
	valuesMap := make(map[int]model.CatalogValue)
	var valuesOrder []int
	sizeMap := make(map[int]model.CatalogSize)
	brandsMap := make(map[int]model.CatalogBrand)
	var brandOrder []int

	var max int
	var min int

	for rows.Next() {
		opt := model.CatalogOption{}
		val := model.CatalogValue{}
		sz := model.CatalogSize{}
		b := model.CatalogBrand{}

		err := rows.Scan(&opt.Id, &opt.Title, &opt.Slug, &val.Id, &val.Value, &val.OptionId, &sz.Value, &sz.Id, &b.Id, &b.Title, &max, &min)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}

		_, ok := optionsMap[opt.Id]
		if !ok {
			optionsMap[opt.Id] = opt
			optionsOrder = append(optionsOrder, opt.Id)
		}

		_, ok = sizeMap[sz.Id]
		if !ok {
			sizeMap[sz.Id] = sz
		}

		_, ok = valuesMap[val.Id]
		if !ok {
			valuesMap[val.Id] = val
			valuesOrder = append(valuesOrder, val.Id)
		}
		_, ok = brandsMap[b.Id]
		if !ok {
			brandsMap[b.Id] = b
			brandOrder = append(brandOrder, b.Id)
		}
	}

	if rows.Err() != nil {
		return nil, fall.ServerError(rows.Err().Error())
	}

	options := make([]model.CatalogOption, 0, len(optionsMap))
	sizes := make([]model.CatalogSize, 0, len(sizeMap))
	brands := make([]model.CatalogBrand, 0, len(brandsMap))

	for _, key := range valuesOrder {
		v := valuesMap[key]
		opt := optionsMap[v.OptionId]
		opt.Values = append(opt.Values, v)
		optionsMap[opt.Id] = opt
	}

	for _, v := range sizeMap {
		sizes = append(sizes, v)
	}

	for _, key := range optionsOrder {
		opt := optionsMap[key]
		options = append(options, opt)
	}

	for _, key := range brandOrder {
		b := brandsMap[key]
		brands = append(brands, b)
	}

	sort.Slice(sizes, func(i, j int) bool {
		a := sizes[i].Value
		b := sizes[j].Value
		return a < b
	})

	return &model.CatalogFilters{
		Options: options,
		Sizes:   sizes,
		Brands:  brands,
		Price:   model.CatalogPrice{Max: max, Min: min},
	}, nil

}

func (r *OptionRepository) UpdateOption(ctx context.Context, dto model.UpdateOptionDto, id int) fall.Error {

	var queries []string

	if dto.ForCatalog != nil {
		queries = append(queries, fmt.Sprintf("for_catalog = %t", *dto.ForCatalog))
	}

	if dto.Title != nil {
		queries = append(queries, fmt.Sprintf("title = '%s'", *dto.Title))
	}

	if dto.Slug != nil {
		queries = append(queries, fmt.Sprintf("slug = '%s'", *dto.Slug))
	}

	if len(queries) > 0 {
		q := "UPDATE option SET " + strings.Join(queries, ",") + " WHERE option_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.OptionUpdateError, err.Error()))
		}
	}

	return nil
}

func (r *OptionRepository) UpdateOptionValue(ctx context.Context, dto model.UpdateOptionValueDto, id int) fall.Error {
	var queries []string

	if dto.Info != nil {
		queries = append(queries, fmt.Sprintf("info = '%s'", *dto.Info))
	}

	if dto.Value != nil {
		queries = append(queries, fmt.Sprintf("value = '%s'", *dto.Value))
	}

	if len(queries) > 0 {
		q := "UPDATE option_value SET " + strings.Join(queries, ",") + " WHERE option_value_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.OptionValueUpdateError, err.Error()))
		}
	}
	return nil
}
