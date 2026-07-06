# StareDesk - Backend

Backend service untuk project IoT StareDesk.  
Project ini berfungsi sebagai pusat komunikasi antara perangkat IoT (ESP32), database, MQTT Broker, dan frontend dashboard.

Backend dibangun menggunakan Golang dengan arsitektur modular berbasis clean architecture agar mudah dikembangkan, dirawat, dan scalable.

[API Documentation](backend/api-docs.md)

---

# Fitur Utama

- REST API menggunakan Gin
- MQTT subscriber & publisher untuk komunikasi realtime dengan ESP32
- WebSocket realtime untuk dashboard monitoring
- JWT Authentication
- Session tracking
- Sensor logging
- Analytics monitoring
- PostgreSQL database
- Docker support

---

# Teknologi yang Digunakan

- Golang
- Gin Framework
- PostgreSQL
- MQTT (HiveMQ / MQTT Broker)
- Gorilla WebSocket
- JWT
- Docker
- PGX PostgreSQL Driver

---

# Arsitektur Project

Project menggunakan pendekatan layered / clean architecture:

```text
Handler -> Usecase -> Repository -> Database
```
## Penjelasan

| Layer          | Fungsi                                   |
| -------------- | ---------------------------------------- |
| Handler        | Menerima request HTTP / MQTT / WebSocket |
| Usecase        | Berisi business logic                    |
| Repository     | Interface akses data                     |
| Infrastructure | Implementasi database dan broker         |
| Entity         | Struktur model data                      |

---

## Cara Menjalankan Project

### 1. Clone Repository
```bash
git clone github.com/seymourrisey/staredesk
cd staredesk/backend
```
### 2. Install Depedency
```bash
go mod tidy
```

### 3. Setup Environment
.env.example
```bash
DATABASE_URL=postgres://postgres:password@localhost:5432/staredesk

JWT_SECRET=your-secret-key

MQTT_BROKER=your-broker.s1.eu.hivemq.cloud
MQTT_PORT=8883
MQTT_USERNAME=your-username
MQTT_PASSWORD=your-password
MQTT_CLIENT_ID=staredesk-backend
MQTT_USER_ID=user-001

APP_PORT=8080
ALLOWED_ORIGINS=http://localhost:3000
```
### 4. Setup Database
```bash
db/staredesk.sql
```
Atau jalankan langsung ke PostgreSQL / Supabase SQL Editor.

### 5. Run
```bash
go run cmd/main.go
```
Server akan berjalan di:
```bash
http://localhost:8080
```
