# 🧭 Що ми в результаті отримаємо

✔ Контролер (Admin UI + API)
✔ Агентів навантаження (EC2)
✔ Реальні HTTP-запити до CDN
✔ Live-метрики (Prometheus + Grafana)
✔ Chaos / Canary / Blue-Green
✔ Auto-stop (cost-guard)

---

# 🔹 ЕТАП 0. Що потрібно перед стартом

## 0.1. Обліковий запис AWS

Потрібно:

* AWS account
* IAM user з **Programmatic access**
* Права:
  * EC2
  * DynamoDB
  * S3
  * IAM
  * SSM
  * CloudWatch

---

## 0.2. Локальна машина

На твоєму комп’ютері має бути:


| Що           | Команда перевірки |
| -------------- | --------------------------------- |
| Git            | `git --version`                   |
| Docker         | `docker --version`                |
| Docker Compose | `docker compose version`          |
| Terraform      | `terraform version`               |
| Go             | `go version`                      |
| AWS CLI        | `aws --version`                   |

👉 Якщо чогось нема — встанови перед продовженням.

---

# 🔹 ЕТАП 1. Клонування репозиторію

```bash
git clone https://github.com/YOUR_ORG/cdn-load-platform.git
cd cdn-load-platform
```

---

# 🔹 ЕТАП 2. AWS credentials (ДУЖЕ ВАЖЛИВО)

## 2.1. Створюємо профіль AWS

```bash
aws configure --profile cdn-load
```

Введи:

```
AWS Access Key ID     = ****************
AWS Secret Access Key = ****************
Default region name   = eu-central-1
Default output format = json
```

Перевір:

```bash
aws sts get-caller-identity --profile cdn-load
```

---

# 🔹 ЕТАП 3. Terraform — інфраструктура

## 3.1. Створюємо S3 (профілі тестів)

```bash
cd terraform/s3
terraform init
terraform apply
```

Запам’ятай:

* **bucket name** (виведе Terraform)

---

## 3.2. Завантажуємо тест-профіль

```bash
aws s3 cp profiles/example.json s3://<BUCKET_NAME>/profiles/example.json
```

---

## 3.3. Контролер (Admin node)

```bash
cd ../controller
terraform init
terraform apply
```

Terraform виведе:

* Public IP контролера

👉 Запиши його.

---

## 3.4. Load-nodes (агенти)

```bash
cd ../load-nodes
terraform init
terraform apply
```

👉 За замовчуванням створиться **2 EC2 агенти**

---

# 🔹 ЕТАП 4. Що відбувається автоматично (важливо)

Коли Terraform створює EC2:

✔ Ubuntu
✔ Docker
✔ Docker Compose
✔ AWS SSM agent
✔ Агент запускається в Docker

**НІЧОГО ВРУЧНУ НА EC2 СТАВИТИ НЕ ТРЕБА**

👉 Це робиться через `user_data` у Terraform.

---

# 🔹 ЕТАП 5. Контролер (Web UI)

## 5.1. Відкрий у браузері

```
http://<CONTROLLER_PUBLIC_IP>:8080
```

Ти побачиш:

* форму запуску тесту
* графіки
* кнопки Chaos / Canary

---

## 5.2. Логін

Файл:

```
docs/auth.md
```

За замовчуванням:

```
login: admin
password: admin123
```

(міняється через ENV)

---

# 🔹 ЕТАП 6. Запуск тесту (живі CDN-запити)

## 6.1. Створюємо тест

У Web UI:

1. **Profile**: `profiles/example.json`
2. **Instances**: `2`
3. **Sessions**: `1000`
4. **Mode**: Normal
5. Натисни **Apply**

⏱ Через \~2 хвилини:

* агенти готові
* state у DynamoDB

---

## 6.2. Старт тесту

Натисни:

```
▶ Start Test
```

👉 **Усі агенти стартують ОДНОЧАСНО**

---

# 🔹 ЕТАП 7. Chaos testing

У UI:

1. Chaos → Enable
2. Latency: 200 ms
3. Error rate: 5%
4. Apply

👉 Через 5 сек chaos активний

---

# 🔹 ЕТАП 8. Grafana (Live metrics)

## 8.1. Відкрий Grafana

```
http://<CONTROLLER_PUBLIC_IP>:3000
```

Логін:

```
admin / admin
```

Dashboards:

* Load overview
* Per-edge latency
* Chaos impact
* Cost guard

---

# 🔹 ЕТАП 9. Звіти (PDF / CSV)

У UI:

```
Export → PDF / CSV
```

Або CLI:

```bash
curl http://controller/api/report/<test_id> -o report.pdf
```

---

# 🔹 ЕТАП 10. Auto-stop (cost guard)

Якщо:

* RPS = 0
* Errors > 90%
* TTL минув

👉 Контролер:

* зупиняє тест
* термінує EC2
* зберігає звіт

---

# 🔹 ЕТАП 11. Очистка ресурсів (ВАЖЛИВО)

```bash
cd terraform/load-nodes
terraform destroy

cd ../controller
terraform destroy

cd ../s3
terraform destroy
```

---

# 🔐 Секрети та конфіги (де лежать)


| Що          | Де                 |
| ------------- | -------------------- |
| JWT secret    | ENV controller       |
| AWS creds     | `~/.aws/credentials` |
| Chaos config  | DynamoDB             |
| Profiles      | S3                   |
| Grafana admin | docker-compose       |

---

# 🧪 Перша “жива” перевірка CDN

Тестуй:

* FO vs Edge latency
* Canary деградацію
* Chaos resilience
* Cost vs performance

---

# 🔚 Якщо хочеш — наступні кроки

Можу:

* підготувати **README.md (1:1 copy)**
* додати **demo профіль CDN**
* зробити **single-command deploy**
* додати **RBAC**

Скажи:
👉 **хочеш це як один Markdown-файл документації чи як downloadable PDF?**
