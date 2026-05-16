# Gadget Shop — End Term Project Guide

## Architecture Overview

```
Internet
  └─► Frontend  (Nginx :80 / NodePort 30080)
        └─► /api/* ──► Gateway :8080 (NodePort 30800)
                          ├──► auth-service        :8081  ──► auth-db
                          ├──► account-service     :8082  ──► account-db
                          ├──► product-service     :8083  ──► product-db
                          ├──► order-service       :8084  ──► order-db
                          ├──► payment-service     :8085  ──► payment-db
                          └──► notification-service:8086  ──► notification-db

Prometheus :9090 (NodePort 30090) ◄── scrapes /metrics on all services
Grafana    :3000 (NodePort 30300) ◄── reads Prometheus + Loki
```

**Stack:**
- 6 Go microservices + API gateway + React frontend
- PostgreSQL (one DB per service)
- Docker Compose (local dev) → Docker Swarm (clustering) → Kubernetes/k3s (prod)
- Terraform (GCP VM provisioning)
- Ansible (automated deployment)
- Prometheus + Grafana (monitoring)

---

## Step 1 — Local Development with Docker Compose

### 1.1 Prerequisites

```bash
docker --version      # Docker 24+
docker compose version  # Docker Compose v2+
```

### 1.2 Start everything locally

```bash
# Clone the repo
git clone https://github.com/nuulun/gadget_shop.git
cd gadget_shop

# .env is already configured for local dev
# Build and start all services
docker compose up --build
```

### 1.3 Verify services are running

```bash
docker compose ps
```

Expected: all containers status `healthy` or `running`.

### 1.4 Access locally

| Service    | URL                        |
|------------|----------------------------|
| Frontend   | http://localhost            |
| Grafana    | http://localhost:3000  (admin/admin) |
| Prometheus | http://localhost:9090       |

### 1.5 Test the API

```bash
# Register a user
curl -s -X POST http://localhost/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","email":"test@test.com","password":"123456","first_name":"Test","last_name":"User"}' | jq

# Login
TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"123456"}' | jq -r '.access_token')

# Create an order
curl -s -X POST http://localhost/api/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":2}]}' | jq

# Process a payment
curl -s -X POST http://localhost/api/payments \
  -H "Content-Type: application/json" \
  -d '{"order_id":1,"amount":99.99,"method":"card"}' | jq

# Send a notification
curl -s -X POST http://localhost/api/notifications/send \
  -H "Content-Type: application/json" \
  -d '{"recipient":"test@test.com","type":"email","message":"Order confirmed!"}' | jq
```

---

## Step 2 — Terraform: Provision GCP VM

### 2.1 Prerequisites

```bash
# Install Terraform
# Install Google Cloud CLI
gcloud auth login
gcloud config set project sre-assignment-494811
```

### 2.2 Configure terraform.tfvars

File `terraform/terraform.tfvars` is already configured:
```
project_id     = "sre-assignment-494811"
region         = "us-central1"
zone           = "us-central1-a"
app_name       = "gadget-shop"
machine_type   = "e2-standard-4"
ssh_user       = "ubuntu"
ssh_public_key = "ssh-ed25519 ..."   # your public key
```

### 2.3 Apply Terraform

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

**What Terraform creates:**
- Static public IP address
- Firewall rules (ports: 22, 80, 3000, 9090, 30080, 30090, 30300, 30800)
- GCP VM (e2-standard-4, Ubuntu 22.04, 30GB SSD)
- Startup script installs: Docker, Docker Compose, k3s

### 2.4 Get VM IP

```bash
terraform output public_ip
# Output: 35.232.58.179
```

### 2.5 Wait for VM to be ready (~3-5 minutes)

```bash
# SSH into VM
ssh ubuntu@35.232.58.179

# Check Docker is installed
docker --version

# Check k3s is installed
sudo k3s kubectl get nodes
```

---

## Step 3 — Ansible: Automated Deployment

### 3.1 Prerequisites

```bash
pip install ansible
```

### 3.2 Configure inventory

File `ansible/inventory.ini` is already configured:
```ini
[gcp_vm]
gadget-shop-vm ansible_host=35.232.58.179 ansible_user=ubuntu ansible_ssh_private_key_file=~/.ssh/id_rsa
```

### 3.3 Setup VM (first time only)

```bash
ansible-playbook -i ansible/inventory.ini ansible/setup.yml
```

**What setup.yml does:**
- Installs Docker, k3s, git
- Clones the repository to `/home/ubuntu/gadget_shop`
- Configures kubeconfig permissions

### 3.4 Deploy application

```bash
ansible-playbook -i ansible/inventory.ini ansible/deploy.yml
```

**What deploy.yml does:**
- Pulls latest code from git
- Builds Docker images for all 8 services
- Imports images into k3s containerd
- Applies all Kubernetes manifests in correct order
- Waits for deployments to be ready

---

## Step 4 — Docker Swarm (Clustering)

Docker Swarm is the alternative to Kubernetes for simpler clustering.

### 4.1 SSH into VM

```bash
ssh ubuntu@35.232.58.179
cd /home/ubuntu/gadget_shop
```

### 4.2 Initialize Swarm

```bash
docker swarm init --advertise-addr 35.232.58.179
```

### 4.3 Deploy the stack

```bash
# Load environment variables
export $(cat .env | grep -v '#' | xargs)

docker stack deploy -c docker-swarm.yml gadget-shop
```

### 4.4 Verify

```bash
docker stack services gadget-shop
docker stack ps gadget-shop
```

### 4.5 Scale a service

```bash
# Scale order-service to 3 replicas
docker service scale gadget-shop_order-service=3
```

### 4.6 Remove stack

```bash
docker stack rm gadget-shop
```

---

## Step 5 — Kubernetes (k3s)

### 5.1 SSH into VM

```bash
ssh ubuntu@35.232.58.179
cd /home/ubuntu/gadget_shop
```

### 5.2 Create secrets (IMPORTANT — do this first)

```bash
# Edit the secrets file with real values
cp k8s/secrets.yml.example k8s/secrets.yml  # or edit directly
nano k8s/secrets.yml
# Fill in real DB passwords
```

The secrets file has this structure:
```yaml
stringData:
  JWT_SECRET: "your-actual-secret"
  AUTH_DB_USER: "postgres"
  AUTH_DB_PASSWORD: "your-password"
  AUTH_DB_NAME: "auth"
  AUTH_DB_DSN: "host=auth-db port=5432 user=postgres password=your-password dbname=auth sslmode=disable"
  # ... same for account, product, order, payment, notification
```

### 5.3 Build and import Docker images

```bash
# Build all images
for svc in auth-service account-service product-service order-service payment-service notification-service gateway frontend; do
  docker build -t $svc:latest ./$svc
  docker save $svc:latest | sudo k3s ctr images import -
done
```

### 5.4 Apply Kubernetes manifests

```bash
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# Apply in correct order
kubectl apply -f k8s/namespace.yml
kubectl apply -f k8s/secrets.yml
kubectl apply -f k8s/prometheus-rbac.yml

# Apply all remaining manifests
kubectl apply -f k8s/
```

### 5.5 Verify pods are running

```bash
kubectl get pods -n gadget-shop
kubectl get services -n gadget-shop
kubectl get hpa -n gadget-shop
```

Expected output:
```
NAME                       READY   STATUS    RESTARTS
auth-db-xxx                1/1     Running   0
auth-service-xxx           1/1     Running   0
account-db-xxx             1/1     Running   0
account-service-xxx        1/1     Running   0
product-db-xxx             1/1     Running   0
product-service-xxx        1/1     Running   0
order-db-xxx               1/1     Running   0
order-service-xxx          1/1     Running   0
payment-db-xxx             1/1     Running   0
payment-service-xxx        1/1     Running   0
notification-db-xxx        1/1     Running   0
notification-service-xxx   1/1     Running   0
gateway-xxx                1/1     Running   0
frontend-xxx               1/1     Running   0
prometheus-xxx             1/1     Running   0
grafana-xxx                1/1     Running   0
node-exporter-xxx          1/1     Running   0
```

### 5.6 Access on GCP VM

| Service    | URL                                |
|------------|------------------------------------|
| Frontend   | http://35.232.58.179:30080          |
| Gateway    | http://35.232.58.179:30800          |
| Grafana    | http://35.232.58.179:30300 (admin/admin) |
| Prometheus | http://35.232.58.179:30090          |

---

## Step 6 — Monitoring (Prometheus + Grafana)

### 6.1 Open Grafana

Go to http://35.232.58.179:30300 → login: `admin` / `admin`

Dashboards are auto-provisioned:
- **Microservices Dashboard** — request rate, error rate, latency, system CPU, system memory
- **Logs Dashboard** — container logs from all services

### 6.2 Verify Prometheus targets

Go to http://35.232.58.179:30090/targets

All targets should show **State: UP**:
- gateway, auth-service, account-service, product-service
- order-service, payment-service, notification-service
- node-exporter, prometheus

### 6.3 Alerts configured

| Alert | Trigger | Severity |
|-------|---------|----------|
| ServiceDown | Any service unreachable for 30s | critical |
| HighErrorRate | HTTP 5xx > 5% for 1m | warning |
| HighCPUUsage | System CPU > 80% for 2m | warning |
| HighMemoryUsage | System Memory > 85% for 2m | warning |

---

## Step 7 — Auto Scaling (HPA)

The `order-service` has Horizontal Pod Autoscaler configured.

### 7.1 Check HPA status

```bash
kubectl get hpa -n gadget-shop
```

```
NAME               REFERENCE            TARGETS   MINPODS   MAXPODS   REPLICAS
order-service-hpa  Deployment/order-service  2%/80%   1         5         1
```

### 7.2 Trigger auto-scaling with load test

```bash
# From your local machine
TOKEN=$(curl -s -X POST http://35.232.58.179:30800/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"123456"}' | jq -r '.access_token')

# Run load test (install apache bench: sudo apt install apache2-utils)
ab -n 5000 -c 50 -H "Authorization: Bearer $TOKEN" \
   http://35.232.58.179:30800/api/orders/my-orders
```

### 7.3 Watch scaling in real time

```bash
kubectl get hpa -n gadget-shop -w
kubectl get pods -n gadget-shop -w
```

---

## Step 8 — Incident Simulation

### 8.1 Simulate a service crash

```bash
# Kill order-service pod
kubectl delete pod -l app=order-service -n gadget-shop

# Watch Kubernetes auto-restart it
kubectl get pods -n gadget-shop -w
```

Kubernetes will restart the pod automatically in ~10 seconds. Prometheus will fire `ServiceDown` alert.

### 8.2 Simulate high error rate

```bash
# Send requests to a bad endpoint to generate 404s
ab -n 1000 -c 10 http://35.232.58.179:30800/api/nonexistent
```

Check Grafana → Error Rate panel.

### 8.3 Simulate CPU stress

```bash
# SSH into VM and run stress
ssh ubuntu@35.232.58.179
stress --cpu 4 --timeout 120s &
```

Watch Grafana → System CPU Usage panel spike above 80% and trigger the alert.

### 8.4 Postmortem template

After incident, document:
```
Incident: [Service/Resource affected]
Duration: [Start time] → [End time]
Impact: [What was unavailable]
Root Cause: [What caused it]
Detection: [How was it detected — Prometheus alert / manual]
Resolution: [What fixed it]
Prevention: [What to do so it doesn't happen again]
```

---

## Step 9 — Capacity Planning (Load Test)

### 9.1 Run load test from local machine

```bash
# Get auth token
TOKEN=$(curl -s -X POST http://35.232.58.179:30800/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"123456"}' | jq -r '.access_token')

# Load test: 1000 requests, 20 concurrent
ab -n 1000 -c 20 \
   -H "Authorization: Bearer $TOKEN" \
   http://35.232.58.179:30800/api/orders/my-orders

# Load test: products (no auth needed)
ab -n 2000 -c 50 http://35.232.58.179:30800/api/products
```

### 9.2 Interpret results

Key metrics from `ab` output:
- **Requests per second** — throughput of the service
- **Time per request** — average latency
- **Failed requests** — error rate
- **95th percentile** — tail latency (find in Transfer rate section)

### 9.3 Record baseline capacity

| Metric | Value |
|--------|-------|
| Max sustainable RPS | ~100 req/s |
| Avg latency at 100 RPS | ~50ms |
| CPU at 100 RPS | ~20-30% |
| Memory at 100 RPS | ~40% |

Compare metrics in Grafana during and after load test.

---

## Project Structure

```
assignment/
├── docker-compose.yml          ← Local development
├── docker-swarm.yml            ← Docker Swarm clustering
│
├── terraform/                  ← GCP VM provisioning
│   ├── main.tf                 ← VM, firewall, static IP
│   ├── variables.tf
│   ├── outputs.tf
│   └── terraform.tfvars        ← Your GCP config
│
├── ansible/                    ← Automated deployment
│   ├── inventory.ini           ← VM host (35.232.58.179)
│   ├── setup.yml               ← Install Docker + k3s
│   └── deploy.yml              ← Build images + apply k8s
│
├── k8s/                        ← Kubernetes manifests
│   ├── namespace.yml
│   ├── secrets.yml             ← DB credentials (not in git)
│   ├── *-db.yml                ← Database deployments (6 total)
│   ├── *-service.yml           ← Service deployments (6 total)
│   ├── gateway.yml
│   ├── frontend.yml
│   ├── prometheus.yml          ← ConfigMap + Deployment + Service
│   ├── prometheus-rbac.yml     ← ServiceAccount + ClusterRole
│   ├── grafana.yml
│   ├── node-exporter.yml       ← DaemonSet for system metrics
│   └── order-service.yml       ← Includes HPA (auto-scaling)
│
├── auth-service/               ← Go: JWT auth (port 8081)
├── account-service/            ← Go: user profiles (port 8082)
├── product-service/            ← Go: product catalog (port 8083)
├── order-service/              ← Go: orders + HPA (port 8084)
├── payment-service/            ← Go: payment simulation (port 8085)
├── notification-service/       ← Go: notification simulation (port 8086)
├── gateway/                    ← Go: API gateway + JWT (port 8080)
├── frontend/                   ← React + Nginx (port 80)
│
└── monitoring/
    ├── prometheus/
    │   ├── prometheus.yml      ← Scrape config (all 6 services)
    │   └── alerts.yml          ← ServiceDown, HighCPU, HighMemory, HighErrorRate
    └── grafana/provisioning/
        ├── datasources/        ← Prometheus auto-configured
        └── dashboards/         ← Microservices + Logs dashboards
```

---

## Quick Reference: All URLs

### Local (Docker Compose)
| Service    | URL |
|------------|-----|
| Frontend   | http://localhost |
| Grafana    | http://localhost:3000 |
| Prometheus | http://localhost:9090 |

### GCP VM (Kubernetes)
| Service    | URL |
|------------|-----|
| Frontend   | http://35.232.58.179:30080 |
| Gateway API | http://35.232.58.179:30800 |
| Grafana    | http://35.232.58.179:30300 |
| Prometheus | http://35.232.58.179:30090 |

### SSH
```bash
ssh ubuntu@35.232.58.179
```

### Useful kubectl commands
```bash
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

kubectl get pods -n gadget-shop
kubectl get services -n gadget-shop
kubectl get hpa -n gadget-shop
kubectl logs -f deployment/order-service -n gadget-shop
kubectl describe pod <pod-name> -n gadget-shop
kubectl rollout restart deployment/gateway -n gadget-shop
```
