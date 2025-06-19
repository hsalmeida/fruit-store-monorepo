package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/hsalmeida/fruit-store-monorepo/api/internal/server"
)

func main() {
	//Conexão com o banco
	dbURL := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Conexão Redis
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_URL")})

	// JWTAuth
	jwtAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	// RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	conn, _ := amqp.Dial(rabbitURL)
	ch, _ := conn.Channel()

	srv, err := server.New(
		server.WithDB(pool),
		server.WithRedis(rdb),
		server.WithJWTAuth(jwtAuth),
		server.WithRabbit(conn, ch),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", srv.Router)
}
