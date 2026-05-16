package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/internal/config"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Gateway struct {
	cfg        *config.Config
	httpClient *http.Client
}

func New(cfg *config.Config) *Gateway {
	return &Gateway{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (g *Gateway) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok", "service": "gateway"})
	})
	mux.Handle("/metrics", promhttp.Handler())

	// Auth
	mux.HandleFunc("/api/auth/register", g.authRegister)
	mux.HandleFunc("/api/auth/login", g.authLogin)
	mux.HandleFunc("/api/auth/logout", g.proxy(g.cfg.AuthURL+"/logout"))
	mux.HandleFunc("/api/auth/refresh", g.proxy(g.cfg.AuthURL+"/refresh"))
	mux.HandleFunc("/api/auth/me", g.authMe)

	// Users
	mux.HandleFunc("/api/users", g.proxy(g.cfg.AccountURL+"/users"))
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		suffix := strings.TrimPrefix(r.URL.Path, "/api/users/")
		g.proxyTo(w, r, g.cfg.AccountURL+"/users/"+suffix)
	})

	// Products
	mux.HandleFunc("/api/products", g.proxy(g.cfg.ProductURL+"/products"))
	mux.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		suffix := strings.TrimPrefix(r.URL.Path, "/api/products/")
		g.proxyTo(w, r, g.cfg.ProductURL+"/products/"+suffix)
	})

	// Orders
	mux.HandleFunc("/api/orders/my-orders", g.ordersMyOrders)
	mux.HandleFunc("/api/orders/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "my-orders") { return }
		suffix := strings.TrimPrefix(r.URL.Path, "/api/orders/")
		g.proxyTo(w, r, g.cfg.OrderURL+"/orders/"+suffix)
	})
	mux.HandleFunc("/api/orders", g.orders)

	// Payments
	mux.HandleFunc("/api/payments/", func(w http.ResponseWriter, r *http.Request) {
		suffix := strings.TrimPrefix(r.URL.Path, "/api/payments/")
		g.proxyTo(w, r, g.cfg.PaymentURL+"/payments/"+suffix)
	})
	mux.HandleFunc("/api/payments", g.proxy(g.cfg.PaymentURL+"/payments"))

	// Notifications
	mux.HandleFunc("/api/notifications/send", g.proxy(g.cfg.NotificationURL+"/notifications/send"))
	mux.HandleFunc("/api/notifications", g.proxy(g.cfg.NotificationURL+"/notifications"))
}

// ─── Auth handlers ────────────────────────────────────────────────────────────

func (g *Gateway) authRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { writeError(w, 405, "method not allowed"); return }

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "invalid body"); return
	}

	// 1. Create account
	accountPayload := map[string]interface{}{
		"login": body["login"], "email": body["email"],
		"phone": body["phone"], "first_name": body["first_name"],
		"last_name": body["last_name"], "middle_name": body["middle_name"],
		"age": body["age"],
	}
	accountResp, err := g.post(r.Context(), g.cfg.AccountURL+"/users", accountPayload)
	if err != nil { writeError(w, 502, "account service error"); return }
	defer accountResp.Body.Close()
	if accountResp.StatusCode != http.StatusCreated {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(accountResp.StatusCode)
		io.Copy(w, accountResp.Body)
		return
	}
	var createdUser map[string]interface{}
	json.NewDecoder(accountResp.Body).Decode(&createdUser)

	userID := uint64(0)
	if id, ok := createdUser["id"].(float64); ok { userID = uint64(id) }

	// 2. Register auth credentials
	authPayload := map[string]interface{}{
		"id": userID, "login": body["login"],
		"email": body["email"], "password": body["password"],
	}
	authResp, err := g.post(r.Context(), g.cfg.AuthURL+"/register", authPayload)
	if err != nil { writeError(w, 502, "auth service error"); return }
	defer authResp.Body.Close()
	if authResp.StatusCode != http.StatusCreated {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(authResp.StatusCode)
		io.Copy(w, authResp.Body)
		return
	}

	// 3. Auto-login
	loginPayload := map[string]interface{}{
		"login_or_email": body["login"],
		"password":       body["password"],
	}
	loginResp, err := g.post(r.Context(), g.cfg.AuthURL+"/login", loginPayload)
	if err != nil { writeError(w, 502, "auth service error on login"); return }
	defer loginResp.Body.Close()

	var tokens map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&tokens)

	writeJSON(w, 201, map[string]interface{}{"user": createdUser, "tokens": tokens})
}

func (g *Gateway) authLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { writeError(w, 405, "method not allowed"); return }
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "invalid body"); return
	}
	// Support both "login" and "email" field names from frontend
	loginOrEmail := ""
	if v, ok := body["login"].(string); ok && v != "" { loginOrEmail = v }
	if v, ok := body["email"].(string); ok && v != "" { loginOrEmail = v }

	authPayload := map[string]interface{}{
		"login_or_email": loginOrEmail,
		"password":       body["password"],
	}
	resp, err := g.post(r.Context(), g.cfg.AuthURL+"/login", authPayload)
	if err != nil { writeError(w, 502, "auth service error"); return }
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (g *Gateway) authMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { writeError(w, 405, "method not allowed"); return }
	userID, err := g.extractUserID(r)
	if err != nil { writeError(w, 401, "unauthorized"); return }
	g.proxyTo(w, r, fmt.Sprintf("%s/users/%d", g.cfg.AccountURL, userID))
}

// ─── Order handlers ───────────────────────────────────────────────────────────

func (g *Gateway) orders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		userID, err := g.extractUserID(r)
		if err != nil { writeError(w, 401, "unauthorized"); return }
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, 400, "invalid body"); return
		}
		body["user_id"] = userID
		g.proxyWithBody(w, r, g.cfg.OrderURL+"/orders", body)
	case http.MethodGet:
		g.proxyTo(w, r, g.cfg.OrderURL+"/orders")
	default:
		writeError(w, 405, "method not allowed")
	}
}

func (g *Gateway) ordersMyOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { writeError(w, 405, "method not allowed"); return }
	userID, err := g.extractUserID(r)
	if err != nil { writeError(w, 401, "unauthorized"); return }
	g.proxyTo(w, r, fmt.Sprintf("%s/orders/my-orders?user_id=%d", g.cfg.OrderURL, userID))
}

// ─── Proxy helpers ────────────────────────────────────────────────────────────

// proxy returns a handler that forwards the request as-is to targetURL.
func (g *Gateway) proxy(targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		g.proxyTo(w, r, targetURL)
	}
}

func (g *Gateway) proxyTo(w http.ResponseWriter, r *http.Request, targetURL string) {
	var bodyReader io.Reader
	if r.Body != nil { bodyReader = r.Body }

	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, bodyReader)
	if err != nil { writeError(w, 500, "proxy error"); return }

	// Forward headers (except hop-by-hop)
	for k, vs := range r.Header {
		if isHopByHop(k) { continue }
		for _, v := range vs { req.Header.Add(k, v) }
	}
	req.URL.RawQuery = r.URL.RawQuery

	resp, err := g.httpClient.Do(req)
	if err != nil {
		log.Printf("[gateway] upstream error %s: %v", targetURL, err)
		writeError(w, 502, "upstream service unavailable")
		return
	}
	defer resp.Body.Close()

	for k, vs := range resp.Header {
		if isHopByHop(k) { continue }
		for _, v := range vs { w.Header().Add(k, v) }
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (g *Gateway) proxyWithBody(w http.ResponseWriter, r *http.Request, targetURL string, body interface{}) {
	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, bytes.NewReader(data))
	if err != nil { writeError(w, 500, "proxy error"); return }
	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = r.URL.RawQuery

	resp, err := g.httpClient.Do(req)
	if err != nil { writeError(w, 502, "upstream service unavailable"); return }
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (g *Gateway) post(ctx interface{ Done() <-chan struct{} }, url string, body interface{}) (*http.Response, error) {
	data, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil { return nil, err }
	req.Header.Set("Content-Type", "application/json")
	return g.httpClient.Do(req)
}

// ─── JWT ──────────────────────────────────────────────────────────────────────

func (g *Gateway) extractUserID(r *http.Request) (uint64, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" { return 0, fmt.Errorf("no auth header") }
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" { return 0, fmt.Errorf("invalid auth header") }

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(g.cfg.JWTSecret), nil
	})
	if err != nil { return 0, fmt.Errorf("invalid token: %w", err) }

	uid, ok := claims["sub"].(float64)
	if !ok { return 0, fmt.Errorf("invalid token subject") }
	return uint64(uid), nil
}

// ─── Middleware ───────────────────────────────────────────────────────────────

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[gateway] %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
	})
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

var hopByHop = map[string]bool{
	"Connection": true, "Keep-Alive": true, "Transfer-Encoding": true,
	"Te": true, "Trailers": true, "Upgrade": true, "Proxy-Authorization": true,
	"Proxy-Authenticate": true,
}

func isHopByHop(header string) bool { return hopByHop[http.CanonicalHeaderKey(header)] }
