package handler

import (
	"account-service/internal/model"
	"account-service/internal/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AccountService interface {
	Create(context.Context, model.CreateUserInput) (model.User, error)
	GetByID(context.Context, uint64) (model.User, error)
	List(context.Context, int, int) ([]model.User, error)
	Update(context.Context, uint64, model.UpdateUserInput) (model.User, error)
	Delete(context.Context, uint64) error
}

type Handler struct{ svc AccountService }

func New(svc AccountService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok", "service": "account-service"})
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/users", h.users)
	mux.HandleFunc("/users/", h.userByID)
}

func (h *Handler) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit == 0 { limit = 20 }
		users, err := h.svc.List(r.Context(), limit, offset)
		if err != nil { writeError(w, 500, err.Error()); return }
		writeJSON(w, 200, users)
	case http.MethodPost:
		var in model.CreateUserInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, 400, "invalid body"); return
		}
		u, err := h.svc.Create(r.Context(), in)
		if err != nil { writeError(w, 400, err.Error()); return }
		writeJSON(w, 201, u)
	default:
		writeError(w, 405, "method not allowed")
	}
}

func (h *Handler) userByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/users/")
	if err != nil { writeError(w, 400, "invalid user id"); return }

	switch r.Method {
	case http.MethodGet:
		u, err := h.svc.GetByID(r.Context(), id)
		if errors.Is(err, service.ErrNotFound) { writeError(w, 404, "user not found"); return }
		if err != nil { writeError(w, 500, err.Error()); return }
		writeJSON(w, 200, u)
	case http.MethodPut:
		var in model.UpdateUserInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, 400, "invalid body"); return
		}
		u, err := h.svc.Update(r.Context(), id, in)
		if errors.Is(err, service.ErrNotFound) { writeError(w, 404, "user not found"); return }
		if err != nil { writeError(w, 400, err.Error()); return }
		writeJSON(w, 200, u)
	case http.MethodDelete:
		err := h.svc.Delete(r.Context(), id)
		if errors.Is(err, service.ErrNotFound) { writeError(w, 404, "user not found"); return }
		if err != nil { writeError(w, 500, err.Error()); return }
		writeJSON(w, 200, map[string]string{"status": "deleted"})
	default:
		writeError(w, 405, "method not allowed")
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
func parseID(path, prefix string) (uint64, error) {
	s := strings.TrimPrefix(path, prefix)
	s = strings.Split(s, "/")[0]
	if s == "" { return 0, fmt.Errorf("no id") }
	return strconv.ParseUint(s, 10, 64)
}
