package scheduler

import (
	"context"
	"log"

	"github.com/go-co-op/gocron"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
)

type ActionScheduler struct {
	cron *gocron.Scheduler
	db   db.PostgresClient
}

func NewActionScheduler(cron *gocron.Scheduler, db db.PostgresClient) *ActionScheduler {
	return &ActionScheduler{cron: cron, db: db}
}

func (s *ActionScheduler) Start() {

	ctx := context.Background()

	go s.discountCancellation(ctx)
}

func (s *ActionScheduler) discountCancellation(ctx context.Context) {
	s.cron.Every(15).Second().Do(func() {
		q := `
    select a.action_id, am.product_model_id from action as a
    inner join action_model as am on a.action_id = am.action_id
    where current_date >= end_date and is_activated = true;
    `

		rows, err := s.db.Query(ctx, q)
		if err != nil {
			return
		}
		defer rows.Close()

		modelMap := make(map[int]int)
		actionMap := make(map[string]string)

		for rows.Next() {
			var modelId int
			var actionId string

			err := rows.Scan(&actionId, &modelId)
			if err != nil {
				return
			}

			_, ok := modelMap[modelId]
			if !ok {
				modelMap[modelId] = modelId
			}
			_, ok = actionMap[actionId]
			if !ok {
				actionMap[actionId] = actionId
			}
		}

		if err := rows.Err(); err != nil {
			return
		}

		modelIds := make([]int, 0, len(modelMap))
		actionIds := make([]string, 0, len(actionMap))

		for k := range modelMap {
			modelIds = append(modelIds, k)
		}
		for k := range actionMap {
			actionIds = append(actionIds, k)
		}

		if len(modelIds) > 0 && len(actionIds) > 0 {

			var txErr error = nil

			tx, err := s.db.Begin(ctx)
			if err != nil {
				txErr = err
				return
			}

			defer func() {
				if txErr != nil {
					tx.Rollback(ctx)
				} else {
					tx.Commit(ctx)
				}
			}()

			updateAction := "update action set is_activated = false where action_id = any ($1);"

			_, err = tx.Exec(ctx, updateAction, actionIds)
			if err != nil {
				txErr = err
				return
			}

			updateModels := "update product_model set discount = null where product_model_id = any ($1);"

			_, err = tx.Exec(ctx, updateModels, modelIds)
			if err != nil {
				txErr = err
				return
			}

			log.Printf("Action scheduler successfully execute update discount operation!, modelIds:%v; actionIds: %v", modelIds, actionIds)
			return
		}
		log.Printf("Action scheduler successfully execute operation without updates!")

	})
}
