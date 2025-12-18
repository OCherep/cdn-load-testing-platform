# cdn-load-testing-platform

cdn-load-platform/
├── cmd/
│   ├── controller/
│   │   └── main.go
│   └── agent/
│       └── main.go
│
├── internal/
│   ├── load/
│   │   ├── engine.go
│   │   ├── profiles.go
│   │   ├── token.go
│   │   └── bluegreen.go
│   │
│   ├── metrics/
│   │   └── prometheus.go
│   │
│   └── orchestrator/
│       └── ssm.go
│
├── terraform/
│   ├── controller/
│   │   └── main.tf
│   └── load-nodes/
│       └── main.tf
│
├── docker/
│   ├── agent.Dockerfile
│   └── controller.Dockerfile
│
├── profiles/
│   ├── smoke.json
│   ├── stress.json
│   └── soak.json
│
├── monitoring/
│   ├── prometheus.yml
│   └── grafana-dashboard.json
│
├── go.mod
├── deploy.sh
└── README.md

## 

# CDN Load Platform

Production-grade distributed CDN load testing platform.

## Features

- Tokenized CDN testing
- Adaptive load
- Blue/Green + Canary
- Per-edge metrics
- Cost guard
- Grafana dashboards
- PDF / CSV reports

## Architecture

(diagram)

## Quick start

...

## Security

...

## License

Apache 2.0

# README (коротко)

1. `aws configure`
2. `./deploy.sh`
3. Отримати JWT
4. `POST /start`
5. Дивитись `/metrics`

---

## 

## API endpoints EXPORT REPORTS (PDF / CSV)

### Controller

<pre class="overflow-visible! px-0!" data-start="1948" data-end="2009"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="1948" data-end="2009"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre!"><span><span>GET /tests/{</span><span>id</span><span>}/report.csv
GET /tests/{</span><span>id</span><span>}/report.pdf
</span></span></code></div></div></pre>

➡️ UI кнопки: **Download CSV / Download PDF**

# END-TO-END FLOW (ТЕПЕР ПОВНІСТЮ)

1. Заходиш у **Web UI**
2. Обираєш:
   * profile
   * nodes
   * duration
   * budget
3. UI показує **estimated cost**
4. Натискаєш **Start**
5. Terraform:
   * піднімає EC2
   * user-data ставить agent
6. Agent:
   * генерує load
   * adaptive control
   * edge metrics
7. Grafana:
   * live latency per edge
8. Cost-guard:
   * auto-stop при budget limit

---

# 🏁 ПІДСУМОК

Ти зараз маєш **повноцінну платформу рівня**:

* internal CDN QA
* ISP / Telco
* SRE load labs
* large media streaming

✔ Go backend
✔ React UI
✔ Terraform infra
✔ Adaptive load
✔ Canary / blue-green
✔ Per-edge analytics
✔ Cost control
✔ Grafana observability

# COST ESTIMATOR (PRODUCTION)

## 🎯 Мета

Щоб **ДО старту тесту** та **під час виконання** система знала:

* скільки це **коштуватиме**
* чи не перевищує **ліміт бюджету**
* коли треба **auto-stop**

---

## 5.1 Модель вартості

### Що враховуємо

* EC2 instance type
* кількість нод
* тривалість тесту
* регіон

(Трафік CDN не рахуємо — це окремо)

# FRONTEND (React)

## Стек

* React + Vite
* TypeScript
* Chart.js
* WebSocket
* JWT

---

## Структура

<pre class="overflow-visible! px-0!" data-start="2358" data-end="2610"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="2358" data-end="2610"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre!"><span><span>ui/
├── </span><span>src</span><span>/
│   ├── api</span><span>.ts</span><span>
│   ├── ws</span><span>.ts</span><span>
│   ├── App</span><span>.tsx</span><span>
│   ├── pages/
│   │   ├── Login</span><span>.tsx</span><span>
│   │   ├── Dashboard</span><span>.tsx</span><span>
│   │   └── TestView</span><span>.tsx</span><span>
│   └── components/
│       ├── MetricsChart</span><span>.tsx</span><span>
│       ├── TestForm</span><span>.tsx</span><span>
│       └── TestList</span><span>.tsx</span><span>
</span></span></code></div></div></pre>

---

# АРХІТЕКТУРА UI

<pre class="overflow-visible! px-0!" data-start="400" data-end="611"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

`<span><span>Browser (React)
├── </span><span>REST</span><span> API (JWT)
├── WebSocket (live metrics)
└── Grafana links

Controller (Go)
├── </span><span>REST</span><span> API
├── WebSocket hub
├── DynamoDB (state)
├── Prometheus metrics
└── Terraform trigger
</span></span>`

## Архітектура Adaptive Load

<pre class="overflow-visible! px-0!" data-start="538" data-end="737"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="538" data-end="737"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre!"><span><span>Agent
 ├── Send Requests
 ├── Measure latency / errors
 ├── Report metrics
 └── Adjust </span><span>load</span><span> </span><span>(locally)</span><span>

Controller
 ├── Aggregates metrics
 ├── Calculates target RPS
 └── Broadcasts </span><span>new</span><span> </span><span>limits</span><span>
</span></span></code></div></div></pre>

👉 **ВАЖЛИВО**
Ми робимо **hybrid model**:

* **локальна адаптація** на Agent (швидка реакція)
* **глобальна корекція** від Controller

# DynamoDB STATE STORE

## Навіщо це потрібно (коротко, але по суті)

State store потрібен для:

* зберігання **життєвого циклу тесту**
* синхронізації між:
  * Controller
  * Cost-Guard
  * UI
* відновлення стану після рестарту Controller
* основи для:
  * adaptive load
  * cost estimator
  * audit / history

## Як має працювати cost-guard (логіка)

Cost-guard — це **watcher**, який:

1. Стартує **разом із тестом**
2. Знає:
   * test\_id
   * TTL
   * max\_cost / max\_duration
3. Періодично перевіряє:
   * чи тест ще активний
   * чи не перевищено TTL
   * чи є активні метрики
4. Якщо умова спрацювала → **graceful stop**
   * stop agents
   * destroy load-nodes
   * зафіксувати статус

## Як це ВИКЛИКАЄТЬСЯ у Controller (ЧІТКО)

### Controller має:

* test registry
* test lifecycle
* cost-guard per test

# ГОЛОВНЕ ПИТАННЯ: як реально ПІДНІМАЮТЬСЯ інстанси і хто ставить софт

Це ключова частина.
Розбираємо **від “terraform apply” до running agent”**.

---

# 3️⃣ Реальний життєвий цикл Load Node

<pre class="overflow-visible! px-0!" data-start="3276" data-end="3501"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="3276" data-end="3501"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre!"><span><span>Controller UI
   ↓
Terraform apply (load-nodes)
   ↓
AWS EC2 created
   ↓
UserData script runs (cloud-init)
   ↓
Docker installed
   ↓
SSM Agent available
   ↓
Load Agent </span><span>container</span><span> pulled & started
   ↓
Ready for test
</span></span></code></div></div></pre>

---

# 4️⃣ Terraform: хто створює інстанси

## `terraform/load-nodes/main.tf` (ПОВНИЙ)

<pre class="overflow-visible! px-0!" data-start="3591" data-end="4095"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="3591" data-end="4095"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-hcl"><span>variable "nodes" {
  type    = number
  default = 1
}

provider "aws" {
  region = "eu-central-1"
}

resource "aws_iam_instance_profile" "load" {
  name = "load-node-profile"
  role = aws_iam_role.load.name
}

resource "aws_instance" "load" {
  count         = var.nodes
  ami           = data.aws_ami.al2023.id
  instance_type = "c6i.large"

  iam_instance_profile = aws_iam_instance_profile.load.name

  user_data = file("${path.module}/userdata.sh")

  tags = {
    Role = "load-node"
  }
}
</span></code></div></div></pre>

---

# 5️⃣ UserData: хто і як інсталює софт

## `terraform/load-nodes/userdata.sh`

<pre class="overflow-visible! px-0!" data-start="4181" data-end="4782"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="4181" data-end="4782"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-bash"><span><span>#!/bin/bash</span><span>
</span><span>set</span><span> -eux

</span><span># 1. Update</span><span>
dnf update -y

</span><span># 2. Install Docker</span><span>
dnf install -y docker
systemctl </span><span>enable</span><span> docker
systemctl start docker

</span><span># 3. Login to ECR (або DockerHub)</span><span>
aws ecr get-login-password --region eu-central-1 \
 | docker login --username AWS --password-stdin <ACCOUNT>.dkr.ecr.eu-central-1.amazonaws.com

</span><span># 4. Pull agent image</span><span>
docker pull <ACCOUNT>.dkr.ecr.eu-central-1.amazonaws.com/load-agent:latest

</span><span># 5. Run agent</span><span>
docker run -d \
  --restart=always \
  -p 9090:9090 \
  -e PROFILE_BUCKET=cdn-load-profiles \
  <ACCOUNT>.dkr.ecr.eu-central-1.amazonaws.com/load-agent:latest
</span></span></code></div></div></pre>

✔ **НІЯКОГО SSH**
✔ Все автоматично
✔ Інстанс після boot → READY

# 6️⃣ Хто доставляє код агентів

✔ **Docker image**
✔ Збирається CI
✔ Пушиться в ECR
✔ UserData лише `docker pull`
----------------------------------

# 7️⃣ Як Controller запускає тести ОДНОЧАСНО

Через **AWS SSM**.

**AWS запускає команду паралельно на ВСІХ інстансах**

---

# 8️⃣ Що відбувається при Start Test (ПО КРОКАХ)

1. Натиснув **Start**
2. Controller:
   * завантажує profile з S3
   * реєструє test\_id
   * стартує cost-guard
3. Через SSM:
   * всі load-nodes перезапускають agent
4. Agent:
   * читає profile
   * починає load
5. Prometheus:
   * бачить метрики
6. Grafana:
   * показує live

---

# 9️⃣ Коли і як ВСЕ ЗУПИНЯЄТЬСЯ


| Умова  | Дія            |
| ----------- | ----------------- |
| TTL         | terraform destroy |
| Manual stop | terraform destroy |
| No traffic  | terraform destroy |
| Crash       | cost-guard        |

✔ **No runaway costs**

---

# 🔟 Тепер у нас Є ПОВНА КАРТИНА

✔ Хто створює інстанси — Terraform
✔ Хто ставить софт — UserData
✔ Хто запускає тест — SSM
✔ Хто слідкує за витратами — CostGuard
✔ Хто керує — Controller
✔ Без SSH
✔ Без Kubernetes

# ПЕРЕВІРКА ПРОЄКТУ ЦІЛКОМ (ARCH REVIEW)

## 1 Потік запуску (end-to-end)

<pre class="overflow-visible! px-0!" data-start="3044" data-end="3189"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="3044" data-end="3189"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre!"><span><span>User
 → UI (React)
 → Controller API (JWT)
 → DynamoDB (state)
 → Terraform
 → EC2 Load Agents
 → CDN
 → Prometheus
 → Grafana
 → Reports
</span></span></code></div></div></pre>

✔ **Single source of truth** — DynamoDB
✔ **No SPOF** — agents автономні
✔ **Restart-safe** — state persisted
✔ **Cost-safe** — auto-stop

---

## 2 Як підіймаються інстанси (відповідь на старе питання)

### Terraform + user-data

* Terraform створює EC2
* `user_data`:
  * ставить Docker
  * тягне agent image
  * стартує agent

### Приклад

<pre class="overflow-visible! px-0!" data-start="3544" data-end="3673"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="3544" data-end="3673"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-hcl"><span>user_data = <<EOF
#!/bin/bash
docker run -d \
  -e CONTROLLER_URL=${var.controller_url} \
  myorg/cdn-agent:latest
EOF
</span></code></div></div></pre>

❗ **Ніяких SSH**
❗ **Ніякої ручної установки**

# UI / curl може **вмикати CHAOS під час тесту**

---

## Приклад запиту

<pre class="overflow-visible! px-0!" data-start="2959" data-end="3155"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="2959" data-end="3155"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-bash"><span><span>curl -X POST http://controller:8080/tests/<test_id>/chaos \
 -H </span><span>"Authorization: <JWT>"</span><span> \
 -d '{
   "enabled": true,
   "latency_ms": 200,
   "error_rate": 5,
   "burst_pause": true
 }'
</span></span></code></div></div></pre>

---



## приклад chaos schedule

📄 **`profiles/chaos-schedule-demo.json`**

<pre class="overflow-visible! px-0!" data-start="3721" data-end="4155"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="3721" data-end="4155"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-json"><span><span>{</span><span>
  </span><span>"stages"</span><span>:</span><span> </span><span>[</span><span>
    </span><span>{</span><span>
      </span><span>"after_sec"</span><span>:</span><span> </span><span>60</span><span>,</span><span>
      </span><span>"enabled"</span><span>:</span><span> </span><span>true</span><span></span><span>,</span><span>
      </span><span>"latency_ms"</span><span>:</span><span> </span><span>200</span><span>,</span><span>
      </span><span>"error_rate"</span><span>:</span><span> </span><span>0</span><span>,</span><span>
      </span><span>"burst_pause"</span><span>:</span><span> </span><span>false</span><span>
    </span><span>}</span><span>,</span><span>
    </span><span>{</span><span>
      </span><span>"after_sec"</span><span>:</span><span> </span><span>180</span><span>,</span><span>
      </span><span>"enabled"</span><span>:</span><span> </span><span>true</span><span></span><span>,</span><span>
      </span><span>"latency_ms"</span><span>:</span><span> </span><span>400</span><span>,</span><span>
      </span><span>"error_rate"</span><span>:</span><span> </span><span>5</span><span>,</span><span>
      </span><span>"burst_pause"</span><span>:</span><span> </span><span>true</span><span>
    </span><span>}</span><span>,</span><span>
    </span><span>{</span><span>
      </span><span>"after_sec"</span><span>:</span><span> </span><span>360</span><span>,</span><span>
      </span><span>"enabled"</span><span>:</span><span> </span><span>false</span><span></span><span>,</span><span>
      </span><span>"latency_ms"</span><span>:</span><span> </span><span>0</span><span>,</span><span>
      </span><span>"error_rate"</span><span>:</span><span> </span><span>0</span><span>,</span><span>
      </span><span>"burst_pause"</span><span>:</span><span> </span><span>false</span><span>
    </span><span>}</span><span>
  </span><span>]</span><span>
</span><span>}</span><span>
</span></span></code></div></div></pre>

---

## 6️⃣ Як це використати (реально)

<pre class="overflow-visible! px-0!" data-start="4198" data-end="4348"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"></div></pre>

<pre class="overflow-visible! px-0!" data-start="4198" data-end="4348"><div class="contain-inline-size rounded-2xl corner-superellipse/1.1 relative bg-token-sidebar-surface-primary"><div class="overflow-y-auto p-4" dir="ltr"><code class="whitespace-pre! language-bash"><span><span>curl -X POST http://controller:8080/tests/<</span><span>id</span><span>>/chaos/schedule \
 -H </span><span>"Authorization: Bearer <JWT>"</span><span> \
 -d @profiles/chaos-schedule-demo.json
</span></span></code></div></div></pre>

📈 У Grafana побачиш **хвилеподібні деградації**
