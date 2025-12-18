Проєкт **структурно завершений**. Є:

✔ Controller (Go)
✔ Agent (Go)
✔ DynamoDB state store (**internal/state/** — реально існує)
✔ Chaos engine + schedules
✔ Geo simulation
✔ Cost estimator
✔ React UI
✔ Terraform:

* controller
* load-nodes
* state-store
* s3 profiles
  ✔ Docker (agent + controller)
  ✔ Grafana-ready metrics
  ✔ Reports PDF/CSV
  ✔ CI (GitHub Actions)

❗ **Проблема не в коді**
❗ **Проблема — у відсутності чіткої послідовності запуску**

Саме це зараз виправляю.

---

# 🧭 ВИСОКОРОВНЕВА ПОСЛІДОВНІСТЬ (1 хв)

```
1. AWS акаунт + IAM
2. Terraform: state + S3
3. Terraform: controller
4. Terraform: load agents
5. Запуск Controller
6. Запуск UI
7. Старт тесту
```

Далі — **ПОКРОКОВО, БУКВАЛЬНО**.

---

# 1️⃣ ПЕРЕДУМОВИ (локальна машина)

## Встанови:

```bash
git
docker
docker-compose
go >= 1.21
node >= 18
terraform >= 1.5
awscli v2
```

Перевір:

```bash
aws --version
terraform version
go version
docker --version
```

---

# 2️⃣ AWS — ОБОВʼЯЗКОВІ НАЛАШТУВАННЯ

## 2.1 AWS credentials (локально)

```bash
aws configure
```

Вводиш:

* Access Key
* Secret
* Region (наприклад `eu-central-1`)

---

## 2.2 IAM (мінімально потрібно)

Користувач / роль має мати доступ до:

* EC2
* DynamoDB
* S3
* SSM
* CloudWatch
* IAM (для instance profile)

👉 **Це описано в `DEPLOYMENT_GUIDE.md` — файл валідний**

---

# 3️⃣ Terraform — ПРАВИЛЬНИЙ ПОРЯДОК (ДУЖЕ ВАЖЛИВО)

## 3.1 DynamoDB state store

```bash
cd terraform/state-store
terraform init
terraform apply
```

👉 Створюється таблиця (наприклад):

```
cdn-load-tests
```

---

## 3.2 S3 для профілів

```bash
cd ../s3
terraform init
terraform apply
```

👉 Отримаєш:

```
PROFILE_BUCKET=<bucket-name>
```

---

## 3.3 Controller інфраструктура

```bash
cd ../controller
terraform init
terraform apply
```

👉 Отримаєш:

* EC2 Controller
* Security Group
* Public IP

---

## 3.4 Load agents (EC2)

```bash
cd ../load-nodes
terraform init
terraform apply
```

👉 Створюється:

* ASG агентів
* IAM role
* SSM доступ

---

# 4️⃣ ПРОФІЛІ ТЕСТІВ (ОБОВʼЯЗКОВО)

## 4.1 Демонстраційний профіль

Файл вже є:

```
profiles/example.json
```

## 4.2 Завантаж у S3

```bash
aws s3 cp profiles/example.json s3://<PROFILE_BUCKET>/example.json
```

---

# 5️⃣ Controller — ЗАПУСК

## 5.1 Environment variables

На **EC2 Controller**:

```bash
export STATE_TABLE=cdn-load-tests
export AWS_REGION=eu-central-1
export JWT_SECRET=supersecret
export PROFILE_BUCKET=<bucket>
```

---

## 5.2 Запуск через Docker

```bash
docker build -f docker/controller.Dockerfile -t cdn-controller .
docker run -d -p 8080:8080 \
  -e STATE_TABLE \
  -e AWS_REGION \
  -e JWT_SECRET \
  -e PROFILE_BUCKET \
  cdn-controller
```

---

# 6️⃣ Agents — ЯК ВОНИ ЗАПУСКАЮТЬСЯ (ВАЖЛИВО)

✔ **Agent НЕ запускається вручну**
✔ Він стартує **автоматично через user-data / Docker**
✔ Він отримує:

```bash
TEST_ID
PROFILE_BUCKET
PROFILE_KEY
STATE_TABLE
```

📄 Це вже **реалізовано у Terraform load-nodes**

---

# 7️⃣ UI (React)

## 7.1 Локально або на Controller

```bash
cd ui
npm install
npm run build
```

Або Docker:

```bash
docker build -t cdn-ui .
docker run -p 3000:3000 cdn-ui
```

---

# 8️⃣ СТАРТ ТЕСТУ (ПЕРШИЙ РАЗ)

## 8.1 Авторизація

```bash
POST /auth/login
```

Отримуєш JWT.

---

## 8.2 Створення тесту

```bash
POST /tests
{
  "profile": "example.json",
  "duration": 1800,
  "min_rps": 100,
  "max_rps": 5000,
  "geo": ["EU", "US"],
  "chaos": false
}
```

✔ Controller пише в DynamoDB
✔ Agents автоматично підхоплюють testID
✔ Починається навантаження

---

# 9️⃣ LIVE CONTROL (ТЕ, ПРО ЩО ТИ ПИТАВ)

## Змінити RPS

```bash
POST /tests/<id>/rps
{ "rps": 8000 }
```

## Pause

```bash
POST /tests/<id>/pause
```

## Resume

```bash
POST /tests/<id>/resume
```

## Stop

```bash
POST /tests/<id>/stop
```

---

# 10️⃣ GRAFANA + ЗВІТИ

✔ Grafana піднімається контейнером
✔ Метрики per-edge вже є
✔ Звіти:

```bash
GET /reports/<test-id>.csv
GET /reports/<test-id>.pdf
```

---

# 🧪 CHAOS / GEO / COST — ВЖЕ ПРАЦЮЄ

* Chaos schedules → `internal/chaos`
* Geo simulation → `internal/geo`
* Cost vs QoE → `internal/cost`

