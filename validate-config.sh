#!/usr/bin/env bash
# validate-config.sh — validate environment and config before deployment
set -euo pipefail

ERRORS=0

check() {
  local name=$1
  local value=$2
  if [[ -z "$value" ]]; then
    echo "[ERROR] $name is not set"
    ERRORS=$((ERRORS + 1))
  else
    echo "[OK]    $name"
  fi
}

check_host() {
  local name=$1
  local host=$2
  local port=$3
  if nc -z -w3 "$host" "$port" 2>/dev/null; then
    echo "[OK]    $name reachable ($host:$port)"
  else
    echo "[WARN]  $name not reachable ($host:$port) — may not be started yet"
  fi
}

echo "============================================"
echo " Configuration Validation"
echo "============================================"

# Load .env
if [[ ! -f .env ]]; then
  echo "[ERROR] .env file not found"
  exit 1
fi
source .env

echo ""
echo "--- Required environment variables ---"
check "JWT_SECRET"          "${JWT_SECRET:-}"
check "AUTH_DB_USER"        "${AUTH_DB_USER:-}"
check "AUTH_DB_PASSWORD"    "${AUTH_DB_PASSWORD:-}"
check "AUTH_DB_NAME"        "${AUTH_DB_NAME:-}"
check "ACCOUNT_DB_USER"     "${ACCOUNT_DB_USER:-}"
check "ACCOUNT_DB_PASSWORD" "${ACCOUNT_DB_PASSWORD:-}"
check "ACCOUNT_DB_NAME"     "${ACCOUNT_DB_NAME:-}"
check "PRODUCT_DB_USER"     "${PRODUCT_DB_USER:-}"
check "PRODUCT_DB_PASSWORD" "${PRODUCT_DB_PASSWORD:-}"
check "PRODUCT_DB_NAME"     "${PRODUCT_DB_NAME:-}"
check "ORDER_DB_USER"       "${ORDER_DB_USER:-}"
check "ORDER_DB_PASSWORD"   "${ORDER_DB_PASSWORD:-}"
check "ORDER_DB_NAME"       "${ORDER_DB_NAME:-}"

echo ""
echo "--- docker-compose.yml ---"
if [[ -f docker-compose.yml ]]; then
  echo "[OK]    docker-compose.yml exists"
  if grep -q "order-db_broken" docker-compose.yml; then
    echo "[ERROR] docker-compose.yml contains broken DB host (order-db_broken)"
    ERRORS=$((ERRORS + 1))
  else
    echo "[OK]    DB hostnames look correct"
  fi
else
  echo "[ERROR] docker-compose.yml not found"
  ERRORS=$((ERRORS + 1))
fi

echo ""
echo "--- Monitoring config ---"
if [[ -f monitoring/prometheus/prometheus.yml ]]; then
  echo "[OK]    prometheus.yml exists"
else
  echo "[ERROR] prometheus.yml not found"
  ERRORS=$((ERRORS + 1))
fi

if [[ -f monitoring/prometheus/alerts.yml ]]; then
  echo "[OK]    alerts.yml exists"
else
  echo "[ERROR] alerts.yml not found"
  ERRORS=$((ERRORS + 1))
fi

echo ""
echo "============================================"
if [[ $ERRORS -eq 0 ]]; then
  echo " All checks passed. Ready to deploy."
  echo "============================================"
  exit 0
else
  echo " $ERRORS error(s) found. Fix before deploying."
  echo "============================================"
  exit 1
fi
