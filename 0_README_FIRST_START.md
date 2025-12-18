# CDN Load Testing Platform

┌──────────────┐
│ Admin EC2    │
│              │
│ deploy_all.sh│
│  ├ Terraform │
│  ├ Grafana   │
│  └ Prometheus│
└──────┬───────┘
│
▼
┌──────────────────────┐
│ Controller (EC2)     │
│ :8080                │
│                      │
│ - REST API           │
│ - JWT auth           │
│ - Chaos control      │
│ - Test lifecycle     │
│ - Cost estimation    │
└──────┬───────────────┘
│ DynamoDB
▼
┌──────────────────────┐
│ State Store          │
│ (Test state, chaos)  │
└──────┬───────────────┘
│
▼
┌──────────────────────┐
│ Load Agents (EC2 ASG)│
│                      │
│ - Geo simulation     │
│ - Chaos injection    │
│ - Adaptive RPS       │
│ - Prometheus metrics │
└─────────┬────────────┘
▼
CDN (Akamai / CF)

┌────────────┐
│   User     │
└─────┬──────┘
│ REST API
┌─────▼──────┐
│ Controller │
│  (Gin)    │
└─────┬──────┘
│ State / Chaos / RPS
┌─────▼──────────────┐
│ DynamoDB TestState │
└─────┬──────────────┘
│ polling
┌─────▼──────┐
│   Agents   │  ← autoscaled EC2
└─────┬──────┘
│ HTTP
┌─────▼─────────────────────┐
│ Broadpeak / CloudFront CDN│
└────────┬──────────────────┘
│ metrics
┌─────▼──────┐
│ Prometheus │
└─────┬──────┘
▼
Grafana

## 🚀 Quick start (10 minutes)

### 1. Prerequisites

* Docker + Docker Compose
* Terraform >= 1.5
* AWS credentials (IAM)

### 2. Clone

```bash
git clone https://github.com/OCherep/cdn-load-testing-platform.git
cd cdn-load-testing-platform
```

### 3. Deploy everything

```bash
chmod +x deploy_all.sh
./deploy_all.sh
```

### 4. Open Grafana

* [http://localhost:3000](http://localhost:3000)
* login: **admin / admin**

### 5. Create test

```bash
curl -X POST http://controller/api/tests \
  -H "Authorization: Bearer <JWT>" \
  -d '{"profile_key":"example.json","nodes":10,"sessions":500}'
```

### 6. Start test

```bash
curl -X POST http://controller/api/tests/<id>/start
```

## 📊 What you get

* Multi-CDN comparison (Akamai vs Cloudflare)
* SLA latency (p95)
* Error rate per region
* Edge stickiness
* Chaos & autoscaling

## 📦 Output

* Grafana dashboards
* Prometheus metrics
* PDF SLA report (optional)

## 🧠 Typical use cases

* CDN vendor comparison
* SLA validation
* Load & chaos testing
* Cost vs QoE analysis
