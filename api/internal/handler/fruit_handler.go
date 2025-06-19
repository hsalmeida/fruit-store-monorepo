package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hsalmeida/fruit-store-monorepo/api/internal/model"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/repository"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/service"
)

// FruitHandler agrupa as dependências
type FruitHandler struct {
	svc   service.FruitService
	cache *redis.Client
}

// NewFruitHandler injeta PostgreSQL e Redis
func NewFruitHandler(db *pgxpool.Pool, cache *redis.Client) *FruitHandler {
	repo := repository.NewFruitRepository(db)
	svc := service.NewFruitService(repo)
	return &FruitHandler{svc: svc, cache: cache}
}

func (h *FruitHandler) List(w http.ResponseWriter, r *http.Request) {

	if data, err := h.cache.Get(r.Context(), "fruits:all").Result(); err == nil {
		w.Write([]byte(data))
		return
	}

	fruits, err := h.svc.ListFruits(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, _ := json.Marshal(fruits)
	h.cache.Set(r.Context(), "fruits:all", jsonData, time.Minute*5)

	json.NewEncoder(w).Encode(fruits)
}

func (h *FruitHandler) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	fruit, err := h.svc.GetFruit(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(fruit)
}

func (h *FruitHandler) Create(w http.ResponseWriter, r *http.Request) {
	var f model.Fruit
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateFruit(r.Context(), &f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.cache.Del(r.Context(), "fruits:all")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

func (h *FruitHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var f model.Fruit
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	f.ID = id
	if err := h.svc.UpdateFruit(r.Context(), &f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.cache.Del(r.Context(), "fruits:all")
	json.NewEncoder(w).Encode(f)
}

func (h *FruitHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteFruit(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
