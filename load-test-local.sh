#!/usr/bin/env bash
# load-test-local.sh — run load simulation from local machine
set -euo pipefail

HOST="${1:-http://34.122.34.46}"
DURATION="${2:-60}"
CONCURRENCY="${3:-50}"

echo "============================================"
echo " Load Simulation (from local machine)"
echo " Host:        $HOST"
echo " Duration:    ${DURATION}s"
echo " Concurrency: $CONCURRENCY users"
echo "============================================"

if ! command -v ab &>/dev/null; then
  echo "ERROR: 'ab' not found."
  echo "Install: https://www.apachelounge.com/download/"
  echo "Or run via WSL: sudo apt-get install apache2-utils"
  exit 1
fi

echo ""
echo "--- Test 1: Concurrent requests to /api/products (${DURATION}s, ${CONCURRENCY} users) ---"
ab -t "$DURATION" -c "$CONCURRENCY" -k \
  -H "Accept: application/json" \
  "${HOST}/api/products" 2>&1 | grep -E "Requests per second|Time per request|Failed requests|Transfer rate"

echo ""
echo "--- Test 2: Repeated requests to /api/products (30s, 20 users) ---"
ab -t 30 -c 20 -k \
  -H "Accept: application/json" \
  "${HOST}/api/products" 2>&1 | grep -E "Requests per second|Failed requests"

echo ""
echo "--- Test 3: Health endpoints (15s, 10 users) ---"
ab -t 15 -c 10 -k "${HOST}/health" 2>&1 | grep -E "Requests per second|Failed requests"

echo ""
echo "============================================"
echo " Load simulation complete."
echo " Check Grafana: ${HOST/80/}:3000"
echo "============================================"
