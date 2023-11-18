package option

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// dto add errs
type OptionRepository struct {
	db *pgxpool.Pool
}

func NewOptionRepository(db *pgxpool.Pool) *OptionRepository {
	return &OptionRepository{db: db}
}

func (or *OptionRepository) CreateOption(dto CreateOptionDto) error {
	query := `
  INSERT INTO option (title, slug) VALUES ($1,$2);
  `

	_, err := or.db.Exec(context.Background(), query, dto.Title, dto.Slug)

	if err != nil {
		return err
	}

	return nil
}

func (or *OptionRepository) UpdateOption(dto UpdateOptionDto, id int) error {
	query := `
  UPDATE option SET title = $1
  WHERE option_id = $2;
  `
	_, err := or.db.Exec(context.Background(), query, dto.Title, id)

	if err != nil {
		return err
	}
	return nil
}

func (or *OptionRepository) GetById(id int) (*Option, error) {
	query := `
  SELECT op.option_id as op_id, op.title as op_title, op.slug as op_slug,
  v.option_value_id as v_id, v.value as v_value, v.info as v_info, v.option_id as v_option_id
  FROM option as op
  LEFT JOIN option_value as v ON v.option_id = op.option_id
  WHERE op.option_id = $1;
  `

	rows, err := or.db.Query(context.Background(), query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	opt := Option{}

	valMap := make(map[int]OptionValue)

	processedRows := 0

	for rows.Next() {
		v := OptionValue{}

		err := rows.Scan(&opt.Id, &opt.Title, &opt.Slug, &v.Id, &v.Value, &v.Info, &v.OptionId)

		if err != nil {
			return nil, err
		}

		if v.Id != nil {
			valMap[*v.Id] = v
		}
		processedRows++
	}

	if rows.Err() != nil {
		return nil, err
	}

	if processedRows == 0 {
		return nil, errors.New("Option not found!")
	}

	values := make([]OptionValue, 0, len(valMap))

	for _, val := range valMap {
		values = append(values, val)
	}

	opt.Values = values

	return &opt, nil
}

func (or *OptionRepository) CreateValue(dto CreateOptionValueDto) error {
	query := `
  INSERT INTO option_value (value, info, option_id) VALUES ($1,$2,$3);
  `
	_, err := or.db.Exec(context.Background(), query, dto.Value, dto.Info, dto.OptionId)

	if err != nil {
		return err
	}

	return nil
}

func (or *OptionRepository) DeleteOption(id int) error {
	query := `
  DELETE FROM option WHERE option_id = $1;
  `
	_, err := or.db.Exec(context.Background(), query, id)

	if err != nil {
		return err
	}

	return nil
}

func (or *OptionRepository) DeleteValue(id int) error {
	query := `
  DELETE FROM option_value WHERE option_value_id = $1;
  `
	_, err := or.db.Exec(context.Background(), query, id)

	if err != nil {
		return err
	}

	return nil
}

func (or *OptionRepository) AddToProductModel(dto AddOptionToProductModelDto) error {
	query := `
	INSERT INTO product_model_option (product_model_id, option_id, option_value_id) VALUES ($1,$2,$3);
	`

	_, err := or.db.Exec(context.Background(), query, dto.ProductModelId, dto.OptionId, dto.ValueId)
	if err != nil {
		return err
	}

	return nil
}
