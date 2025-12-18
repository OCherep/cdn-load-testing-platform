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



# ПОВНИЙ КАТАЛОГ API

## 🧭 БАЗОВІ ХОСТИ ТА ПОРТИ


| Компонент | Host                     | Port   | Призначення            |
| ------------------ | ------------------------ | ------ | --------------------------------- |
| Controller API     | `http://<controller-ip>` | `8080` | Керування тестами |
| Controller Metrics | `http://<controller-ip>` | `2112` | Prometheus scrape                 |
| Agent Metrics      | `http://<agent-ip>`      | `2112` | Edge / latency                    |
| Grafana            | `http://<controller-ip>` | `3000` | Dashboards                        |
| Prometheus         | `http://<controller-ip>` | `9090` | Metrics DB                        |

> **Auth:** всі `/api/*` → `Authorization: Bearer <JWT>`

---

## 🔐 AUTH

### POST `/api/auth/login`

**Host:** Controller:8080
**Що робить:** видає JWT

**Request**

<pre class="overflow-visible! px-0!" data-start="1000" data-end="1060"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="1000" data-end="1060"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"username"</span><span>:</span><span> </span><span>"admin"</span><span>,</span><span>
  </span><span>"password"</span><span>:</span><span> </span><span>"admin"</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

**Effect**

* створює JWT
* без нього всі інші API → 401

---

## 🧪 TEST LIFECYCLE

### POST `/api/tests`

**Старт нового тесту**

**Request**

<pre class="overflow-visible! px-0!" data-start="1204" data-end="1300"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="1204" data-end="1300"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"profile"</span><span>:</span><span> </span><span>"profiles/cdn-demo.json"</span><span>,</span><span>
  </span><span>"instances"</span><span>:</span><span> </span><span>10</span><span>,</span><span>
  </span><span>"duration_sec"</span><span>:</span><span> </span><span>1800</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

**Що відбувається**

* створюється `test_id`
* Terraform scale EC2
* запис у state store
* агенти стартують

---

### GET `/api/tests`

**Список тестів**

**Вплив:** НІ
**Виклик:** read-only

---

### GET `/api/tests/{id}`

**Стан тесту**

Повертає:

* статус (`running|paused|stopped`)
* rps
* chaos
* SLA config
* TTL

---

### POST `/api/tests/{id}/stop`

**Зупинка тесту**

**Що викликає**

* агенти припиняють load
* cost-guard → scale down
* статус = stopped

---

### POST `/api/tests/{id}/extend`

**Продовження тесту**

<pre class="overflow-visible! px-0!" data-start="1826" data-end="1858"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="1826" data-end="1858"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"seconds"</span><span>:</span><span> </span><span>600</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

---

## ⚙ LIVE CONTROL

### POST `/api/tests/{id}/rps`

**Live зміна навантаження**

<pre class="overflow-visible! px-0!" data-start="1944" data-end="1974"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="1944" data-end="1974"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"rps"</span><span>:</span><span> </span><span>15000</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

**Ефект**

* миттєво передається агентам
* adaptive engine оновлює rate

---

### POST `/api/tests/{id}/scale`

**Зміна кількості EC2**

<pre class="overflow-visible! px-0!" data-start="2111" data-end="2144"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="2111" data-end="2144"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"instances"</span><span>:</span><span> </span><span>25</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

**Ефект**

* Terraform / ASG scale
* нові агенти auto-join

---

## 🧨 CHAOS

### POST `/api/tests/{id}/chaos`

<pre class="overflow-visible! px-0!" data-start="2257" data-end="2356"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="2257" data-end="2356"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"enabled"</span><span>:</span><span> </span><span>true</span><span></span><span>,</span><span>
  </span><span>"latency_ms"</span><span>:</span><span> </span><span>300</span><span>,</span><span>
  </span><span>"error_rate"</span><span>:</span><span> </span><span>0.05</span><span>,</span><span>
  </span><span>"burst_pause"</span><span>:</span><span> </span><span>1000</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

**Ефект**

* агенти починають chaos middleware
* впливає **перед HTTP request**

---

### POST `/api/tests/{id}/chaos/schedule`

<pre class="overflow-visible! px-0!" data-start="2486" data-end="2579"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="2486" data-end="2579"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"enable_at"</span><span>:</span><span> </span><span>"2025-01-01T10:00:00Z"</span><span>,</span><span>
  </span><span>"disable_at"</span><span>:</span><span> </span><span>"2025-01-01T10:10:00Z"</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

---

## 🧲 STICKINESS / EDGE AFFINITY

### GET `/api/tests/{id}/stickiness`

**Повертає**

<pre class="overflow-visible! px-0!" data-start="2671" data-end="2792"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="2671" data-end="2792"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"global_ratio"</span><span>:</span><span> </span><span>0.93</span><span>,</span><span>
  </span><span>"per_edge"</span><span>:</span><span> </span><span>{</span><span>
    </span><span>"edge-1"</span><span>:</span><span> </span><span>0.98</span><span>,</span><span>
    </span><span>"edge-2"</span><span>:</span><span> </span><span>0.87</span><span>
  </span><span>}</span><span>,</span><span>
  </span><span>"sla_breach"</span><span>:</span><span> </span><span>false</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

---

## 🌍 GEO

### POST `/api/tests/{id}/geo`

<pre class="overflow-visible! px-0!" data-start="2842" data-end="2918"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="2842" data-end="2918"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"regions"</span><span>:</span><span> </span><span>[</span><span>"eu"</span><span>,</span><span> </span><span>"us"</span><span>,</span><span> </span><span>"asia"</span><span>]</span><span>,</span><span>
  </span><span>"weights"</span><span>:</span><span> </span><span>[</span><span>50</span><span>,</span><span> </span><span>30</span><span>,</span><span> </span><span>20</span><span>]</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

---

## 📄 REPORTS

### GET `/api/tests/{id}/report/pdf`

➡️ **SLA Evidence PDF**

---

### GET `/api/tests/{id}/report/csv`

➡️ Raw metrics export


## Коли ти викликаєш:

<pre class="overflow-visible! px-0!" data-start="4496" data-end="4618"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="4496" data-end="4618"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-bash"><span><span>curl -H </span><span>"Authorization: Bearer $JWT</span><span>" \
http://controller:8080/api/tests/report/pdf?test_id=abc123 \
-o sla.pdf
</span></span></code></div></div></pre>

➡️ Controller:

1. бере **реальний state**
2. бере **реальні метрики**
3. порівнює з **SLA з тесту**
4. генерує **юридично валідний PDF**
5. віддає файл

END OF DOCUMENT
