package handlers

import (
	models_ "Practice2/internal/models"
	storage_ "Practice2/internal/storage"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	store *storage_.MemoryStore
}

func NewHandler(s *storage_.MemoryStore) *Handler {
	return &Handler{s}
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		h.getTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodPatch:
		h.updateTask(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (h *Handler) getTasks(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		json.NewEncoder(w).Encode(h.store.GetAll())
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}

	task, ok := h.store.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "task not found",
		})
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {

	var dto models_.CreateTaskDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if dto.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid title",
		})
		return
	}

	task := h.store.Create(dto.Title, dto.Description)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}

	var dto models_.UpdateTaskDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := h.store.Update(id, dto.Completed)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "task not found",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"updated": true,
	})
}
