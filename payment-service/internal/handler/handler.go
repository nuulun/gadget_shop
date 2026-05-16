package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"payment-service/internal/model"
	"payment-service/internal/service"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PaymentService interface {
	Process(context.Context, model.CreatePaymentInput) (model.Payment, error)
	GetByID(context.Context, uint64) (model.Payment, error)
	GetByOrderID(context.Context, uint64) (model.Payment, error)
	List(context.Context) ([]model.Payment, error)
}

type Handler struct{ svc PaymentService }

func New(svc PaymentService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok", "service": "payment-service"})
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/payments/", h.paymentByID)
	mux.HandleFunc("/payments", h.payments)
}

func (h *Handler) payments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var in model.CreatePaymentInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, 400, "invalid body")
			return
		}
		p, err := h.svc.Process(r.Context(), in)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, 201, p)
	case http.MethodGet:
		orderIDStr := r.URL.Query().Get("order_id")
		if orderIDStr != "" {
			orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
			if err != nil {
				writeError(w, 400, "invalid order_id")
				return
			}
			p, err := h.svc.GetByOrderID(r.Context(), orderID)
			if errors.Is(err, service.ErrNotFound) {
				writeError(w, 404, "payment not found")
				return
			}
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			writeJSON(w, 200, p)
			return
		}
		payments, err := h.svc.List(r.Context())
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, 200, payments)
	default:
		writeError(w, 405, "method not allowed")
	}
}

func (h *Handler) paymentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, 405, "method not allowed")
		return
	}
	id, err := parseID(r.URL.Path, "/payments/")
	if err != nil {
		writeError(w, 400, "invalid payment id")
		return
	}
	p, err := h.svc.GetByID(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, 404, "payment not found")
		return
	}
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, p)
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
	if s == "" {
		return 0, fmt.Errorf("no id")
	}
	return strconv.ParseUint(s, 10, 64)
}
