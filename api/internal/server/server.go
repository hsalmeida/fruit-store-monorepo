package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-redis/redis/v8"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/auth"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/handler"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/publisher"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	_ "github.com/hsalmeida/fruit-store-monorepo/api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	Router     *chi.Mux
	DB         *pgxpool.Pool
	Redis      *redis.Client
	JWTAuth    *jwtauth.JWTAuth
	RabbitConn *amqp.Connection
	RabbitChan *amqp.Channel
}

// Option configura um Server
type Option func(*Server) error

// New cria um Server aplicando todas as Option
func New(opts ...Option) (*Server, error) {
	s := &Server{Router: chi.NewRouter()}
	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}
	s.setupRoutes()
	return s, nil
}

// opções disponíveis
func WithDB(pool *pgxpool.Pool) Option {
	return func(s *Server) error {
		s.DB = pool
		return nil
	}
}

func WithRedis(client *redis.Client) Option {
	return func(s *Server) error {
		s.Redis = client
		return nil
	}
}

func WithJWTAuth(auth *jwtauth.JWTAuth) Option {
	return func(s *Server) error {
		s.JWTAuth = auth
		return nil
	}
}

func WithRabbit(conn *amqp.Connection, ch *amqp.Channel) Option {
	return func(s *Server) error {
		s.RabbitConn = conn
		s.RabbitChan = ch
		return nil
	}
}

func (s *Server) setupRoutes() {

	s.Router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Rota pública de login utilizando auth.GenerateToken
	s.Router.Route("/auth/login", func(r chi.Router) {
		handler := handler.NewLoginHandler(s.DB)
		r.Post("/", handler.Login)
	})

	// Adiciona middleware JWT às rotas protegidas
	s.Router.Route("/fruits", func(r chi.Router) {
		handler := handler.NewFruitHandler(s.DB, s.Redis)
		r.Use(jwtauth.Verifier(s.JWTAuth))
		r.Use(auth.MustAuth)
		//List & Get: só admin OU user
		r.With(auth.RoleAuth("admin", "user")).Get("/", handler.List)
		r.With(auth.RoleAuth("admin", "user")).Get("/{id}", handler.Get)

		//Create/Update/Delete: só admin
		r.With(auth.RoleAuth("admin")).Post("/", handler.Create)
		r.With(auth.RoleAuth("admin")).Put("/{id}", handler.Update)
		r.With(auth.RoleAuth("admin")).Delete("/{id}", handler.Delete)
	})

	s.Router.Route("/users", func(r chi.Router) {
		pub := publisher.NewRabbitPublisher(s.RabbitChan, "user.queue")
		handler := handler.NewUserHandler(s.DB, pub)
		r.Use(jwtauth.Verifier(s.JWTAuth))
		r.Use(auth.MustAuth, auth.RoleAuth("admin"))
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
	})
	
}
