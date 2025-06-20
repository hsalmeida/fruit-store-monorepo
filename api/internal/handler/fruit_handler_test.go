package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/handler"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/model"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/service"
)

// mockService implementa FruitService para os testes
type mockService struct {
	listFruits []model.Fruit
	listErr    error

	getFruit model.Fruit
	getErr   error

	createErr error
}

func (m *mockService) ListFruits(ctx context.Context) ([]model.Fruit, error) {
	return m.listFruits, m.listErr
}
func (m *mockService) GetFruit(ctx context.Context, id uuid.UUID) (model.Fruit, error) {
	return m.getFruit, m.getErr
}
func (m *mockService) CreateFruit(ctx context.Context, f *model.Fruit) error {
	// atribui um ID para verificar no teste
	f.ID = uuid.New()
	return m.createErr
}
func (m *mockService) UpdateFruit(ctx context.Context, f *model.Fruit) error { return nil }
func (m *mockService) DeleteFruit(ctx context.Context, id uuid.UUID) error   { return nil }

// newHandler monta um FruitHandler usando o mockService e um Redis que sempre falha (para pular cache)
func newHandler(ms service.FruitService) *handler.FruitHandler {
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	h := handler.NewFruitHandler(nil, rdb)
	// injetamos o mockService no handler (supondo que você tenha um método WithService)
	return h.WithService(ms)
}

func TestListFruits_Success(t *testing.T) {
	fruits := []model.Fruit{{ID: uuid.New(), Name: "Banana", Price: 1.23, Quantity: 10}}
	h := newHandler(&mockService{listFruits: fruits})

	req := httptest.NewRequest(http.MethodGet, "/fruits", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, recebeu %d", rec.Code)
	}
	var got []model.Fruit
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("falha ao decodificar JSON: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Banana" {
		t.Errorf("resposta inesperada: %#v", got)
	}
}

func TestListFruits_Error(t *testing.T) {
	h := newHandler(&mockService{listErr: errors.New("erro DB")})
	req := httptest.NewRequest(http.MethodGet, "/fruits", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("esperado 500, recebeu %d", rec.Code)
	}
}

func TestGetFruit_Success(t *testing.T) {
	id := uuid.New()
	fruit := model.Fruit{ID: id, Name: "Maçã", Price: 2.34, Quantity: 5}
	h := newHandler(&mockService{getFruit: fruit})

	req := httptest.NewRequest(http.MethodGet, "/fruits/"+id.String(), nil)
	// injeta variável de rota para o chi.URLParam funcionar
	rctx := chi.NewRouteContext()
	// 2) adiciona o parâmetro de rota
	rctx.URLParams.Add("id", id.String())
	// 3) injeta esse RouteContext no Context da request
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	h.Get(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, recebeu %d", rec.Code)
	}
	var got model.Fruit
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("falha ao decodificar: %v", err)
	}
	if got.ID != id {
		t.Errorf("esperado ID %s, recebido %s", id, got.ID)
	}
}

func TestGetFruit_BadID(t *testing.T) {
	h := newHandler(&mockService{})
	req := httptest.NewRequest(http.MethodGet, "/fruits/invalid-uuid", nil)
	// injeta variável de rota para o chi.URLParam funcionar
	rctx := chi.NewRouteContext()
	// 2) adiciona o parâmetro de rota
	rctx.URLParams.Add("id", "invalid-uuid")
	// 3) injeta esse RouteContext no Context da request
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	h.Get(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, recebeu %d", rec.Code)
	}
}

func TestCreateFruit_Success(t *testing.T) {
	body := `{"name":"Laranja","price":3.21,"quantity_in_stock":7}`
	h := newHandler(&mockService{})

	req := httptest.NewRequest(http.MethodPost, "/fruits", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Create(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("esperado 201, recebeu %d", rec.Code)
	}
	var got model.Fruit
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("falha ao decodificar: %v", err)
	}
	if got.Name != "Laranja" {
		t.Errorf("esperado name=Laranja, recebeu %s", got.Name)
	}
}

func TestCreateFruit_BadRequest(t *testing.T) {
	h := newHandler(&mockService{})
	req := httptest.NewRequest(http.MethodPost, "/fruits", bytes.NewBufferString("not-json"))
	rec := httptest.NewRecorder()

	h.Create(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, recebeu %d", rec.Code)
	}
}
