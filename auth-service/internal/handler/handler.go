package handler

import (
	"auth-service/internal/model"
	"auth-service/internal/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AuthService interface {
	Register(context.Context, model.RegisterInput) error
	Login(context.Context, model.LoginInput) (model.TokenPair, error)
	Validate(context.Context, string) (uint64, error)
	Refresh(context.Context, string) (model.TokenPair, error)
	Logout(context.Context, string) error
	DeleteUser(context.Context, uint64) error
}

type Handler struct{ svc AuthService }

func New(svc AuthService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/register", h.register)
	mux.HandleFunc("/login", h.login)
	mux.HandleFunc("/validate", h.validate)
	mux.HandleFunc("/refresh", h.refresh)
	mux.HandleFunc("/logout", h.logout)
	mux.HandleFunc("/delete", h.deleteUser)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok", "service": "auth-service"})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed"); return
	}
	var in model.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "invalid body"); return
	}
	if err := h.svc.Register(r.Context(), in); err != nil {
		writeError(w, 400, err.Error()); return
	}
	writeJSON(w, 201, map[string]string{"status": "registered"})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed"); return
	}
	var in model.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "invalid body"); return
	}
	pair, err := h.svc.Login(r.Context(), in)
	if errors.Is(err, service.ErrInvalidCredentials) {
		writeError(w, 401, "invalid credentials"); return
	}
	if err != nil {
		writeError(w, 500, err.Error()); return
	}
	writeJSON(w, 200, pair)
}

func (h *Handler) validate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed"); return
	}
	var body struct{ AccessToken string `json:"access_token"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "invalid body"); return
	}
	uid, err := h.svc.Validate(r.Context(), body.AccessToken)
	if errors.Is(err, service.ErrInvalidToken) {
		writeError(w, 401, "invalid token"); return
	}
	if err != nil {
		writeError(w, 500, err.Error()); return
	}
	writeJSON(w, 200, map[string]uint64{"user_id": uid})
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed"); return
	}
	var body struct{ RefreshToken string `json:"refresh_token"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "invalid body"); return
	}
	pair, err := h.svc.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		writeError(w, 401, err.Error()); return
	}
	writeJSON(w, 200, pair)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed"); return
	}
	var body struct{ RefreshToken string `json:"refresh_token"` }
	json.NewDecoder(r.Body).Decode(&body)
	h.svc.Logout(r.Context(), body.RefreshToken)
	writeJSON(w, 200, map[string]string{"status": "logged out"})
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		writeError(w, 405, "method not allowed"); return
	}
	var body struct{ ID uint64 `json:"id"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "invalid body"); return
	}
	if err := h.svc.DeleteUser(r.Context(), body.ID); err != nil {
		writeError(w, 500, err.Error()); return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
