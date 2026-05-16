package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"notification-service/internal/model"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type NotificationService interface {
	Send(context.Context, model.SendNotificationInput) (model.Notification, error)
	List(context.Context) ([]model.Notification, error)
}

type Handler struct{ svc NotificationService }

func New(svc NotificationService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok", "service": "notification-service"})
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/notifications/send", h.send)
	mux.HandleFunc("/notifications", h.list)
}

func (h *Handler) send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed")
		return
	}
	var in model.SendNotificationInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "invalid body")
		return
	}
	if in.Recipient == "" || in.Message == "" {
		writeError(w, 400, "recipient and message are required")
		return
	}
	n, err := h.svc.Send(r.Context(), in)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, n)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, 405, "method not allowed")
		return
	}
	notifications, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, notifications)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
