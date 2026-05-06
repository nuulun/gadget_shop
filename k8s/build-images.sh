#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

for svc in auth-service account-service product-service order-service gateway frontend; do
  docker build -t $svc:latest ./$svc
  docker save $svc:latest | sudo k3s ctr images import -
done
