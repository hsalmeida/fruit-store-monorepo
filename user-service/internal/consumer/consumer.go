package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/hsalmeida/fruit-store-monorepo/user-service/internal/model"
	"github.com/hsalmeida/fruit-store-monorepo/user-service/internal/repository"
)

type Consumer struct {
	repo repository.UserRepository
}

// New cria um Consumer usando apenas o DB pool
func New(db *pgxpool.Pool) *Consumer {
	return &Consumer{repo: repository.NewUserRepository(db)}
}

// Payload esperado na fila
type userEvent struct {
	Action string     `json:"action"`
	User   model.User `json:"user"`
}

// Process trata create/update/delete e persiste no DB
func (c *Consumer) Process(d amqp.Delivery) error {
	var evt userEvent
	if err := json.Unmarshal(d.Body, &evt); err != nil {
		return err
	}

	ctx := context.Background()
	switch evt.Action {
	case "create", "update":
		if err := c.repo.Save(ctx, evt.User); err != nil {
			return err
		}
		log.Printf("usuario %s: %s", evt.Action, evt.User.Username)
	default:
		log.Printf("acao desconhecida %q, pulando", evt.Action)
	}
	return nil
}
