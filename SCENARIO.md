# ПОВНИЙ КАТАЛОГ API ЗАПИТІВ

## 🔐 AUTH


| Method | URL           | Effect               |
| ------ | ------------- | -------------------- |
| POST   | `/auth/login` | Отримати JWT |

---

## 🧪 TEST LIFECYCLE


| Method | URL                  | Що робить         |
| ------ | -------------------- | ------------------------- |
| POST   | `/tests`             | Створює тест   |
| GET    | `/tests`             | Список тестів |
| POST   | `/tests/{id}/start`  | Запуск              |
| POST   | `/tests/{id}/stop`   | Зупинка            |
| POST   | `/tests/{id}/extend` | +TTL                      |

---

## ⚡ RUNTIME CONTROL


| Method | URL                  | Наслідок     |
| ------ | -------------------- | -------------------- |
| POST   | `/tests/{id}/rps`    | Міняє RPS       |
| POST   | `/tests/{id}/pause`  | Пауза           |
| POST   | `/tests/{id}/resume` | Продовжити |

---

## 🔥 CHAOS


| Method | URL                          | Дія         |
| ------ | ---------------------------- | -------------- |
| POST   | `/tests/{id}/chaos`          | Chaos now      |
| POST   | `/tests/{id}/chaos/schedule` | Chaos schedule |

---

## 📡 STREAMING


| Method | URL              |
| ------ | ---------------- |
| WS     | `/ws/tests/{id}` |

---

# 6️⃣ СЦЕНАРІЇ ВИКОРИСТАННЯ

### ✔️ CDN comparison

Broadpeak vs CloudFront

### ✔️ SLA evidence

PDF + Grafana

### ✔️ Geo-edge

Region → latency

### ✔️ Chaos

Packet loss, delay, DNS

### ✔️ Cost vs QoE

CloudFront baseline
