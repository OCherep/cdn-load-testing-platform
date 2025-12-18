# CDN Load Platform — Deployment & Usage Guide

## 1. Overview

This platform provides production-grade CDN load testing with:

* Tokenized URLs
* Adaptive load
* Canary & Blue/Green testing
* Chaos testing
* Per-edge analytics
* Cost guard
* Grafana observability
* PDF / CSV reports

---

## 2. Prerequisites (Local Machine)

### Required tools

* Docker >= 24
* Docker Compose
* Terraform >= 1.5
* Go >= 1.22
* AWS CLI >= 2

---

## 3. AWS Preparation

### 3.1 AWS Account

* Separate AWS account for load testing
* IAM user or role with:
  * EC2
  * DynamoDB
  * S3
  * SSM
  * IAM PassRole

### 3.2 Configure AWS CLI

```bash
aws configure
AWS Access Key ID: XXXXX
AWS Secret Access Key: XXXXX
Region: eu-central-1
```

---

## 4. Secrets & Configuration

### 4.1 Environment variables

Create `.env`:

```env
AWS_REGION=eu-central-1
STATE_TABLE=cdn-load-tests
PROFILE_BUCKET=cdn-load-profiles
JWT_SECRET=change_me
```

### 4.2 Grafana credentials

Default (DEV ONLY):

```
User: admin
Password: admin
```

Change via environment variables in production.

---

## 5. Infrastructure Deployment

### 5.1 State Store

```bash
cd terraform/state-store
terraform init
terraform apply
```

### 5.2 S3 Profiles

```bash
cd terraform/s3
terraform apply
```

Upload test profiles:

```bash
aws s3 cp profiles/example.json s3://cdn-load-profiles/
```

---

## 6. Local Admin Stack

```bash
docker compose up -d
```

Services:

* Controller: [http://localhost:8080](http://localhost:8080)
* UI: [http://localhost:5173](http://localhost:5173)
* Grafana: [http://localhost:3000](http://localhost:3000)

---

## 7. Running a Test (UI)

1. Open UI
2. Login (admin/admin)
3. Create test:
   * profile: stress.json
   * nodes: 50
   * duration: 30m
   * budget: \$10
   * canary: 5%
4. Start test

---

## 8. Chaos Testing

Enable Chaos:

* latency injection
* error rate
* burst mode

Chaos is applied only to selected traffic:

* canary only
* specific edge groups

Safe by design.

---

## 9. Observability

### Grafana dashboards

* Edge latency heatmap
* p95/p99 per POP
* Error rate spikes
* Cost burn-rate

---

## 10. Reports

After test completion:

* Download CSV
* Download PDF executive report

---

## 11. Cost Control

* Cost estimated before start
* Live burn-rate
* Auto-stop on budget exceed

---

## 12. Shutdown & Cleanup

```bash
terraform destroy
docker compose down
```

---

## 13. Security Notes

* No SSH to load nodes
* IAM roles only
* Secrets never in code
* JWT tokens short-lived

---

## 14. When to Use

✔ CDN validation
✔ Canary rollout
✔ Pre-event testing
✔ Edge troubleshooting

❌ No limits
❌ No canary
❌ No budget

---

END OF DOCUMENT
