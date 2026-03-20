package handler

//
//import (
//	"Practice5/internal/models"
//	"Practice5/internal/repository"
//	"encoding/json"
//	"net/http"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type UserHandler struct {
//	repo *repository.Repository
//}
//
//func NewUserHandler(repo *repository.Repository) *UserHandler {
//	return &UserHandler{repo: repo}
//}
//
//func (h *UserHandler) users(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	switch r.Method {
//	case http.MethodGet:
//		ctx := r.Context()
//		query := r.URL.Query()
//
//		filter := &models.UserFilter{}
//
//		if name := query.Get("name"); name != "" {
//			filter.Name.String = name
//			filter.Name.Valid = true
//		}
//
//		if offset := query.Get("offset"); offset != "" {
//			v, err := strconv.Atoi(offset)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusBadRequest)
//				return
//			}
//			filter.Offset = uint64(v)
//		}
//
//		if email := query.Get("email"); email != "" {
//			filter.Email.String = email
//			filter.Email.Valid = true
//		}
//
//		if gender := query.Get("gender"); gender != "" {
//			filter.Gender.String = gender
//			filter.Gender.Valid = true
//		}
//
//		if birthDate := query.Get("birth_date"); birthDate != "" {
//			t, err := time.Parse("2026-01-02", birthDate)
//			if err != nil {
//				http.Error(w, "invalid birth_date", http.StatusBadRequest)
//				return
//			}
//			filter.BirthDate.Time = t
//			filter.BirthDate.Valid = true
//		}
//
//		if limit := query.Get("limit"); limit != "" {
//			v, err := strconv.Atoi(limit)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusBadRequest)
//				return
//			}
//			filter.Limit = uint64(v)
//		}
//
//		if sort := query.Get("sort"); sort != "" {
//			filter.SortBy.String = sort
//			filter.SortBy.Valid = true
//		}
//
//		if order := query.Get("order"); order != "" {
//			filter.Order.String = order
//			filter.Order.Valid = true
//		}
//
//		if showDeleted := query.Get("show_deleted"); showDeleted == "true" {
//			filter.ShowDeleted = true
//		}
//
//		users, err := h.repo.ListByFilter(ctx, filter)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		json.NewEncoder(w).Encode(users)
//
//	case http.MethodPost:
//		var user models.User
//
//		err := json.NewDecoder(r.Body).Decode(&user)
//
//		if err != nil {
//			http.Error(w, "invalid body", http.StatusBadRequest)
//			return
//		}
//
//		id, err := h.repo.CreateUser(&user)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		json.NewEncoder(w).Encode(map[string]any{
//			"id": id,
//		})
//
//	default:
//		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
//	}
//}
//
//func (h *UserHandler) userByID(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
//
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		http.Error(w, "invalid user id", http.StatusBadRequest)
//		return
//	}
//
//	switch r.Method {
//	case http.MethodGet:
//		user, err := h.repo.GetUserByID(id)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		json.NewEncoder(w).Encode(user)
//
//	case http.MethodPut, http.MethodPatch:
//		var user models.User
//		err := json.NewDecoder(r.Body).Decode(&user)
//		if err != nil {
//			http.Error(w, "invalid body", http.StatusBadRequest)
//			return
//		}
//		user.ID = int64(id)
//		err = h.repo.UpdateUser(&user)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		json.NewEncoder(w).Encode(map[string]string{
//			"status": "updated",
//		})
//
//	case http.MethodDelete:
//		err := h.repo.DeleteUser(id)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusNotFound)
//			return
//		}
//		json.NewEncoder(w).Encode(map[string]string{
//			"status": "deleted",
//		})
//
//	default:
//		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
//
//	}
//}
//
//func (h *UserHandler) CommonFriends(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	ctx := r.Context()
//
//	if r.Method != http.MethodGet {
//		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	query := r.URL.Query()
//
//	user1 := query.Get("user1")
//	user2 := query.Get("user2")
//
//	user1ID, err := strconv.ParseInt(user1, 10, 64)
//	if err != nil {
//		http.Error(w, "invalid user1 id", http.StatusBadRequest)
//		return
//	}
//
//	user2ID, err := strconv.ParseInt(user2, 10, 64)
//	if err != nil {
//		http.Error(w, "invalid user2 id", http.StatusBadRequest)
//		return
//	}
//
//	friends, err := h.repo.GetCommonFriends(ctx, user1ID, user2ID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	json.NewEncoder(w).Encode(friends)
//}
