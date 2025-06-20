# fruit-store-monorepo

Desafio de Codificação: CRUD de Loja de Frutas em Go (com Microserviço Adicional)

[![CI](https://github.com/hsalmeida/fruit-store-monorepo/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/hsalmeida/fruit-store-monorepo/actions/workflows/ci.yml)

---

## Visão Geral

Este monorepo contém:

- **API** em Go para gerenciar frutas com CRUD, cache Redis e autorização baseada em JWT/roles.  
- **Microserviço** em Go (consumer) que consome eventos de criação/atualização/exclusão de usuários via RabbitMQ e persiste em um banco isolado.  
- **Infraestrutura**: PostgreSQL (fruits_db), PostgreSQL separado (users_db), Redis e RabbitMQ, orquestrados via Docker Compose.  
    - Optei por PostgreSQL por oferecer transações ACID, relações fortes entre entidades e consultas SQL poderosas num ecossistema maduro — garantindo consistência e facilidade de evolução do sistema.
- **Migrações** com `golang-migrate`, incluindo seed de um usuário admin padrão.

---

## Pré-requisitos

- Docker & Docker Compose (v2+)
- Make
- (Opcional) jq, psql ou DBeaver para inspecionar bancos

---

## Configuração e Execução

1. Clone este repositório e entre na pasta raiz:
   ```bash
   git clone https://github.com/hsalmeida/fruit-store-monorepo.git
   cd fruit-store-monorepo

2. Ajuste variáveis de ambiente (você pode copiar .env.example para .env):
    ```env
    DATABASE_URL=postgres://postgres:senha@postgres:5432/fruits_db?sslmode=disable
    REDIS_URL=redis://redis:6379/0
    RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
    JWT_SECRET=uma_chave_secreta

3. Suba toda a stack e aplique migrações com um único comando:
    ```bash
    make start

Isso vai:
- Subir Postgres (fruits_db), Postgres-users (users_db), Redis e RabbitMQ
- Executar migrações (criando tabelas e seed do admin)
- Buildar e subir a API (localhost:8080)
- Buildar e subir o consumer de usuários

4. Para parar e remover tudo:
    ```bash
    make down

5. Confira logs em tempo real:
    ```bash
    make logs

## Endpoints Principais

### 1. Autenticação
- Envio
    ```curl
    curl --location 'localhost:8080/auth/login' \
    --header 'Content-Type: application/json' \
    --data '{
        "username": "admin",
        "password": "adminpass"
    }'

- Retorno
    ```json
    { "token": "<JWT_TOKEN>" }

### 2. Usuários (somente admin)
- Criar usuários
    ```curl
    curl --location 'localhost:8080/users' \
    --header 'Authorization: Bearer $TOKEN' \
    --header 'Content-Type: application/json' \
    --data '{
        "username": "hiltinho",
        "password": "hiltinho",
        "role": "user"
    }'

- Listar usuários
    ```curl
    curl -X GET http://localhost:8080/users \
    --header 'Authorization: Bearer $TOKEN'

### 3. Frutas
- Listar todas (admin & user)
    ```curl
    curl -X GET http://localhost:8080/fruits \
    -H "Authorization: Bearer $TOKEN"

- Obter por ID (admin & user)
    ```curl
    curl -X GET http://localhost:8080/fruits/{id} \
    -H "Authorization: Bearer $TOKEN"

- Criar nova (admin)
    ```curl
    curl -X POST http://localhost:8080/fruits \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Banana","price":2.50,"quantity":100}'

- Atualizar (admin)
    ```curl
    curl -X PUT http://localhost:8080/fruits/{id} \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Banana Prata","price":3.00,"quantity":120}'

- Deletar (admin)
    ```curl
    curl -X DELETE http://localhost:8080/fruits/{id} \
    -H "Authorization: Bearer $TOKEN"

### Ferramentas Adicionais
- Swagger UI

    Acesse a documentação interativa em
    ```bash
    http://localhost:8080/swagger/index.html

- Acesso ao Banco (fruits_db)
    ```bash
    psql "postgres://postgres:senha@localhost:5432/fruits_db?sslmode=disable"

- Acesso ao Banco (users_db)
    ```bash
    psql "postgres://postgres:senha@localhost:5433/users_db?sslmode=disable"

### Estrutura
```csharp
├── api/
│   ├── cmd/           # main.go
│   ├── internal/
│   │   ├── handler/   # HTTP handlers
│   │   ├── service/   # regras de negócio
│   │   ├── repository/  
│   │   ├── model/     
│   │   ├── auth/      # JWT, middleware
│   │   └── server/    # server.go
│   ├── migrations/    # SQL migrations + seed admin
│   ├── docs/          # Swagger gerado
│   └── Dockerfile
├── user-service/
│   ├── cmd/           # consumer main.go
│   ├── internal/
│   │   ├── consumer/  # lógica de fila
│   │   ├── repository/
│   │   └── model/
│   ├── migrations/
│   └── Dockerfile
├── docker-compose.yml
└── Makefile

