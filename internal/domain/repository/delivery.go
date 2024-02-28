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

type DeliveryRepository struct {
	db db.PostgresClient
}

func NewDeliveryRepository(db db.PostgresClient) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

func (r *DeliveryRepository) Create(ctx context.Context, dto model.CreateDeliveryPointDto) fall.Error {
	query := `INSERT INTO delivery_point
  (title,city,address,coords,with_fitting,work_schedule,info)
  VALUES ($1,$2,$3,$4,$5,$6,$7);
  `
	_, err := r.db.Exec(ctx, query, dto.Title, dto.City, dto.Address, dto.Coords, dto.WithFitting, dto.WorkSchedule, dto.Info)
	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}

func (r *DeliveryRepository) SearchPoints(ctx context.Context, text string, withFitting bool) ([]model.DeliveryPoint, fall.Error) {

	filter := ""

	if withFitting {
		filter = "AND with_fitting = true"
	}

	query := fmt.Sprintf(`
	SELECT delivery_point_id,title,city,address,coords,with_fitting,work_schedule,info
	FROM delivery_point
	WHERE CONCAT(city, ' ', address, ' ', title) ILIKE $1 %s ORDER BY city, delivery_point_id; 
	`, filter)
	rows, err := r.db.Query(ctx, query, "%"+text+"%")
	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var points []model.DeliveryPoint

	for rows.Next() {
		p := model.DeliveryPoint{}

		err := rows.Scan(&p.Id, &p.Title, &p.City, &p.Address, &p.Coords, &p.WithFitting, &p.WorkSchedule, &p.Info)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		points = append(points, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return points, nil
}

func (r *DeliveryRepository) FindById(ctx context.Context, id int) (*model.DeliveryPoint, fall.Error) {
	query := `SELECT delivery_point_id,title,city,address,coords,with_fitting,work_schedule,info
	FROM delivery_point WHERE delivery_point_id=$1;`

	row := r.db.QueryRow(ctx, query, id)

	p := model.DeliveryPoint{}

	err := row.Scan(&p.Id, &p.Title, &p.City, &p.Address, &p.Coords, &p.WithFitting, &p.WorkSchedule, &p.Info)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.DeliveryPointNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}
	return &p, nil
}

func (r *DeliveryRepository) Update(ctx context.Context, dto model.UpdateDeliveryPointDto, id int) fall.Error {
	var queries []string

	if dto.Address != nil {
		queries = append(queries, fmt.Sprintf("address = '%s'", *dto.Address))
	}

	if dto.City != nil {
		queries = append(queries, fmt.Sprintf("city = '%s'", *dto.City))
	}

	if dto.Coords != nil {
		queries = append(queries, fmt.Sprintf("coords = '%s'", *dto.Coords))
	}

	if dto.Info != nil {
		queries = append(queries, fmt.Sprintf("info = '%s'", *dto.Info))
	}

	if dto.Title != nil {
		queries = append(queries, fmt.Sprintf("title = '%s'", *dto.Title))
	}

	if dto.WithFitting != nil {
		queries = append(queries, fmt.Sprintf("with_fitting = %t", *dto.WithFitting))
	}

	if dto.WorkSchedule != nil {
		queries = append(queries, fmt.Sprintf("work_schedule = '%s'", *dto.WorkSchedule))
	}

	if len(queries) > 0 {
		q := "UPDATE delivery_point SET " + strings.Join(queries, ",") + " WHERE delivery_point_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.DeliveryPointUpdateError, err.Error()))
		}
		return nil
	}
	return nil
}

func (r *DeliveryRepository) Delete(ctx context.Context, id int) fall.Error {
	q := "DELETE FROM delivery_point WHERE delivery_point_id = $1;"

	_, err := r.db.Exec(ctx, q, id)

	if err != nil {
		return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.DeliveryPointDeleteError, err.Error()))
	}

	return nil
}
