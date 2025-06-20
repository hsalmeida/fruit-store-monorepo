# Makefile para orquestrar infra, migrações e subida dos serviços

infra: ## Sobe infra: Postgres (fruits_db), Postgres-users (users_db), Redis e RabbitMQ
	docker-compose up -d postgres postgres-users redis rabbitmq

migrate: infra ## Executa migrações da API e do user-service
	docker-compose run --rm migrate-api
	docker-compose run --rm migrate-user-service

build: ## Builda as imagens da API e do user-service
	docker-compose build api user-service

up: migrate build ## Sobe API e user-service após infra e migrações
	docker-compose up -d api user-service

start: ## Sobe infra, executa migrate e depois sobe as aplicações
	$(MAKE) migrate
	$(MAKE) build
	$(MAKE) up

down: ## Para e remove todos os containers e redes criados
	docker-compose down

logs: ## Segue os logs de todos os serviços
	docker-compose logs -f

test: ## Roda testes unitários em API e user-service
	cd api && go test ./...
	cd user-service && go test ./...
