package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"product-service/internal/model"
	"product-service/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ProductService is the interface handlers depend on.
type ProductService interface {
	List(ctx context.Context, f model.ListFilter) (service.ListResult, error)
	GetByID(ctx context.Context, id uint64) (model.Product, error)
	Create(ctx context.Context, in model.CreateProductInput) (model.Product, error)
	Update(ctx context.Context, id uint64, in model.UpdateProductInput) (model.Product, error)
	Delete(ctx context.Context, id uint64) error
}

// Handler bundles all HTTP handlers for the product service.
type Handler struct {
	svc ProductService
}

// New creates a Handler.
func New(svc ProductService) *Handler {
	return &Handler{svc: svc}
}

// Register mounts all routes on mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/products", h.products)   // GET list / POST create
	mux.HandleFunc("/products/", h.productID) // GET/PUT/DELETE by id
}

// ─── Route handlers ───────────────────────────────────────────────────────────

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "product-service"})
}

// products handles GET /products and POST /products.
func (h *Handler) products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listProducts(w, r)
	case http.MethodPost:
		h.createProduct(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// productID handles GET /products/{id}, PUT /products/{id}, DELETE /products/{id}.
func (h *Handler) productID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/products/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getProduct(w, r, id)
	case http.MethodPut:
		h.updateProduct(w, r, id)
	case http.MethodDelete:
		h.deleteProduct(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ─── Action handlers ──────────────────────────────────────────────────────────

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := model.ListFilter{
		Category: q.Get("category"),
		Limit:    queryInt(q.Get("limit"), 20),
		Offset:   queryInt(q.Get("offset"), 0),
	}
	if v := q.Get("min_price"); v != "" {
		filter.MinPrice, _ = strconv.ParseFloat(v, 64)
	}
	if v := q.Get("max_price"); v != "" {
		filter.MaxPrice, _ = strconv.ParseFloat(v, 64)
	}

	result, err := h.svc.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var in model.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p, err := h.svc.Create(r.Context(), in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request, id uint64) {
	p, err := h.svc.GetByID(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request, id uint64) {
	var in model.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p, err := h.svc.Update(r.Context(), id, in)
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request, id uint64) {
	err := h.svc.Delete(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseID(path, prefix string) (uint64, error) {
	suffix := strings.TrimPrefix(path, prefix)
	suffix = strings.Split(suffix, "/")[0]
	if suffix == "" {
		return 0, fmt.Errorf("no id")
	}
	return strconv.ParseUint(suffix, 10, 64)
}

func queryInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
