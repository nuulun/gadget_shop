package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/internal/model"
	"order-service/internal/service"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type OrderService interface {
	Create(context.Context, model.CreateOrderInput) (model.Order, error)
	GetByID(context.Context, uint64) (model.Order, error)
	ListByUserID(context.Context, uint64) ([]model.Order, error)
	List(context.Context) ([]model.Order, error)
}

type Handler struct{ svc OrderService }

func New(svc OrderService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(mux *http.ServeMux) {

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok", "service": "order-service"})
	})
	
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/orders/my-orders", h.myOrders)
	mux.HandleFunc("/orders/", h.orderByID)
	mux.HandleFunc("/orders", h.orders)
}

func (h *Handler) orders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var in model.CreateOrderInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, 400, "invalid body"); return
		}
		o, err := h.svc.Create(r.Context(), in)
		if err != nil { writeError(w, 400, err.Error()); return }
		writeJSON(w, 201, o)
	case http.MethodGet:
		orders, err := h.svc.List(r.Context())
		if err != nil { writeError(w, 500, err.Error()); return }
		writeJSON(w, 200, orders)
	default:
		writeError(w, 405, "method not allowed")
	}
}

func (h *Handler) myOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { writeError(w, 405, "method not allowed"); return }
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" { userIDStr = r.Header.Get("X-User-ID") }
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil { writeError(w, 400, "user_id required"); return }
	orders, err := h.svc.ListByUserID(r.Context(), userID)
	if err != nil { writeError(w, 500, err.Error()); return }
	writeJSON(w, 200, orders)
}

func (h *Handler) orderByID(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "my-orders") { return }
	id, err := parseID(r.URL.Path, "/orders/")
	if err != nil { writeError(w, 400, "invalid order id"); return }
	switch r.Method {
	case http.MethodGet:
		o, err := h.svc.GetByID(r.Context(), id)
		if errors.Is(err, service.ErrNotFound) { writeError(w, 404, "order not found"); return }
		if err != nil { writeError(w, 500, err.Error()); return }
		writeJSON(w, 200, o)
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
