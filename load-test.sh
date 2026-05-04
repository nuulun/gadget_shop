#!/usr/bin/env bash
# load-test.sh — simulate load for capacity planning analysis
set -euo pipefail

HOST="${1:-http://localhost}"
DURATION="${2:-60}"
CONCURRENCY="${3:-50}"

echo "============================================"
echo " Load Simulation"
echo " Host:        $HOST"
echo " Duration:    ${DURATION}s"
echo " Concurrency: $CONCURRENCY users"
echo "============================================"

# Install tools if missing
if ! command -v ab &>/dev/null; then
  echo "Installing apache2-utils..."
  sudo apt-get install -y apache2-utils &>/dev/null
fi
if ! command -v stress &>/dev/null; then
  echo "Installing stress..."
  sudo apt-get install -y stress &>/dev/null
fi

echo ""
echo "--- Test 1: Concurrent user requests (products) ---"
ab -t "$DURATION" -c "$CONCURRENCY" -k \
  -H "Accept: application/json" \
  "${HOST}/api/products" 2>&1 | grep -E "Requests per second|Time per request|Failed requests|Transfer rate"

echo ""
echo "--- Test 2: Repeated API calls (health endpoints) ---"
ab -t 30 -c 20 -k "${HOST}/api/products" 2>&1 | \
  grep -E "Requests per second|Failed requests"

echo ""
echo "--- Test 3: CPU stress (30s, 2 cores) ---"
echo "Stressing CPU for 30 seconds..."
stress --cpu 2 --timeout 30s &
STRESS_PID=$!

# While stressing, hit the API
ab -t 30 -c 10 -k "${HOST}/api/products" 2>&1 | \
  grep -E "Requests per second|Time per request|Failed requests"

wait $STRESS_PID 2>/dev/null || true

echo ""
echo "============================================"
echo " Load simulation complete."
echo " Check Grafana for metrics:"
echo "   http://$(hostname -I | awk '{print $1}'):3000"
echo "============================================"
