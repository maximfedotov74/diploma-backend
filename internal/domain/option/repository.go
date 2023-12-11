package option

import (
	"context"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type OptionRepository struct {
	db *pgxpool.Pool
}

func NewOptionRepository(db *pgxpool.Pool) *OptionRepository {
	return &OptionRepository{db: db}
}

func (or *OptionRepository) CreateOption(dto CreateOptionDto) exception.Error {
	query := `
  INSERT INTO option (title, slug) VALUES ($1,$2);
  `

	_, err := or.db.Exec(context.Background(), query, dto.Title, dto.Slug)

	if err != nil {
		return exception.NewErr(optionCreateError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (or *OptionRepository) UpdateOption(dto UpdateOptionDto, id int) exception.Error {
	query := `
  UPDATE option SET title = $1
  WHERE option_id = $2;
  `
	_, err := or.db.Exec(context.Background(), query, dto.Title, id)

	if err != nil {
		return exception.NewErr(optionUpdateError, exception.STATUS_INTERNAL_ERROR)
	}
	return nil
}

func (or *OptionRepository) GetAll() ([]Option, exception.Error) {
	query := `
	SELECT op.option_id as op_id, op.title as op_title, op.slug as op_slug,
  v.option_value_id as v_id, v.value as v_value, v.info as v_info, v.option_id as v_option_id
  FROM option as op
  LEFT JOIN option_value as v ON v.option_id = op.option_id
	`
	rows, err := or.db.Query(context.Background(), query)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	defer rows.Close()

	optionsMap := make(map[int]Option)
	valuesMap := make(map[int]OptionValue)

	var founded bool = false

	for rows.Next() {
		opt := Option{}
		v := OptionValue{}

		err := rows.Scan(&opt.Id, &opt.Title, &opt.Slug, &v.Id, &v.Value, &v.Info, &v.OptionId)

		if err != nil {
			return nil, exception.ServerError(err.Error())
		}

		if v.Id != nil {
			valuesMap[*v.Id] = v
		}

		optionsMap[opt.Id] = opt
		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, exception.NewErr(optionNotFound, exception.STATUS_NOT_FOUND)
	}

	for _, v := range valuesMap {
		opt := optionsMap[*v.OptionId]
		opt.Values = append(opt.Values, v)
		optionsMap[opt.Id] = opt
	}

	result := make([]Option, 0, len(optionsMap))

	for _, v := range optionsMap {
		result = append(result, v)
	}

	return result, nil

}

func (or *OptionRepository) GetById(id int) (*Option, exception.Error) {
	query := `
  SELECT op.option_id as op_id, op.title as op_title, op.slug as op_slug,
  v.option_value_id as v_id, v.value as v_value, v.info as v_info, v.option_id as v_option_id
  FROM option as op
  LEFT JOIN option_value as v ON v.option_id = op.option_id
  WHERE op.option_id = $1;
  `

	rows, err := or.db.Query(context.Background(), query, id)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	defer rows.Close()

	opt := Option{}

	valMap := make(map[int]OptionValue)

	founded := false

	for rows.Next() {
		v := OptionValue{}

		err := rows.Scan(&opt.Id, &opt.Title, &opt.Slug, &v.Id, &v.Value, &v.Info, &v.OptionId)

		if err != nil {
			return nil, exception.ServerError(err.Error())
		}

		if v.Id != nil {
			valMap[*v.Id] = v
		}
		if !founded {
			founded = true
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, exception.NewErr(optionNotFound, exception.STATUS_NOT_FOUND)
	}

	values := make([]OptionValue, 0, len(valMap))

	for _, val := range valMap {
		values = append(values, val)
	}

	opt.Values = values

	return &opt, nil
}

func (or *OptionRepository) CreateValue(dto CreateOptionValueDto) exception.Error {
	query := `
  INSERT INTO option_value (value, info, option_id) VALUES ($1,$2,$3);
  `
	_, err := or.db.Exec(context.Background(), query, dto.Value, dto.Info, dto.OptionId)

	if err != nil {
		return exception.NewErr(optionCreateError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (or *OptionRepository) DeleteOption(id int) exception.Error {
	query := `
  DELETE FROM option WHERE option_id = $1;
  `
	_, err := or.db.Exec(context.Background(), query, id)

	if err != nil {
		return exception.NewErr(optionDeleteError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (or *OptionRepository) DeleteValue(id int) exception.Error {
	query := `
  DELETE FROM option_value WHERE option_value_id = $1;
  `
	_, err := or.db.Exec(context.Background(), query, id)

	if err != nil {
		return exception.NewErr(valueDeleteError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (or *OptionRepository) AddToProductModel(dto AddOptionToProductModelDto) exception.Error {
	query := `
	INSERT INTO product_model_option (product_model_id, option_id, option_value_id) VALUES ($1,$2,$3);
	`

	_, err := or.db.Exec(context.Background(), query, dto.ProductModelId, dto.OptionId, dto.ValueId)
	if err != nil {
		return exception.NewErr(addProductError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (or *OptionRepository) AddSizeToProductModel(dto AddSizeToProductModelDto) exception.Error {
	query := `INSERT INTO model_sizes (product_model_id, size_id, in_stock) VALUES ($1,$2, $3);`

	_, err := or.db.Exec(context.Background(), query, dto.ProductModelId, dto.SizeId, dto.InStock)
	if err != nil {
		return exception.NewErr(addProductError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (or *OptionRepository) CreateSize(dto CreateSizeDto) exception.Error {
	query := `INSERT INTO sizes (numeric_size, literal_size) VALUES ($1,$2);`

	_, err := or.db.Exec(context.Background(), query, dto.Numeric, dto.Literal)
	if err != nil {
		return exception.NewErr(sizeCreateError, exception.STATUS_INTERNAL_ERROR)
	}

	return nil
}

func (or *OptionRepository) GetCatalogFilters(categorySlug string) (*CatalogFilters, exception.Error) {
	query := `
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
	sz.size_value as size_value, sz.size_id as size_id
	FROM product as p
	inner join category_tree as ct on p.category_id = ct.category_id
	inner join brand as b on p.brand_id = b.brand_id
	inner join product_model as pm on pm.product_id = p.product_id
	inner join product_model_option as pmop on pmop.product_model_id = pm.product_model_id
	inner join option as op on op.option_id = pmop.option_id
	inner join option_value as v on v.option_value_id = pmop.option_value_id
	inner join model_sizes as ms on ms.product_model_id = pm.product_model_id
	inner join sizes as sz on ms.size_id = sz.size_id;
	`

	rows, err := or.db.Query(context.Background(), query, categorySlug)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	defer rows.Close()

	optionsMap := make(map[int]CatalogOption)
	valuesMap := make(map[int]CatalogValue)
	sizeMap := make(map[int]CatalogSize)

	for rows.Next() {
		opt := CatalogOption{}
		val := CatalogValue{}
		sz := CatalogSize{}

		err := rows.Scan(&opt.Id, &opt.Title, &opt.Slug, &val.Id, &val.Value, &val.OptionId, &sz.Value, &sz.Id)

		if err != nil {
			return nil, exception.ServerError(err.Error())
		}

		_, ok := optionsMap[opt.Id]
		if !ok {
			optionsMap[opt.Id] = opt
		}

		_, ok = sizeMap[sz.Id]
		if !ok {
			sizeMap[sz.Id] = sz
		}

		_, ok = valuesMap[val.Id]
		if !ok {
			valuesMap[val.Id] = val
		}
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	options := make([]CatalogOption, 0, len(optionsMap))
	sizes := make([]CatalogSize, 0, len(sizeMap))

	for _, v := range valuesMap {
		opt := optionsMap[v.OptionId]
		opt.Values = append(opt.Values, v)
		optionsMap[opt.Id] = opt
	}

	for _, v := range sizeMap {
		sizes = append(sizes, v)
	}

	for _, v := range optionsMap {
		options = append(options, v)
	}

	sort.Sort(CatalogSizeSorter(sizes))

	return &CatalogFilters{
		Options: options,
		Sizes:   sizes,
	}, nil

}
