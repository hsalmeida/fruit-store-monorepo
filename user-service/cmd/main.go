package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/hsalmeida/fruit-store-monorepo/user-service/internal/consumer"
)

func main() {
	// Conecta no RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("rabbitmq dial error: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("rabbitmq channel error: %v", err)
	}
	defer ch.Close()

	// Garante que a fila existe
	q, err := ch.QueueDeclare("user.queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	// 2) Conecta no Postgres isolado (users_db)
	dbURL := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	// 3) Inicia o consumer
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}
	c := consumer.New(pool)
	log.Println("user-service consumer started")
	for d := range msgs {
		if err := c.Process(d); err != nil {
			log.Printf("process error: %v", err)
			d.Nack(false, true)
			continue
		}
		d.Ack(false)
	}
}
