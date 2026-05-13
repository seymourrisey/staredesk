<p align="left">
  <img src="/backend/Logo-Title.png" width="600">
</p>

# StareDesk - Sit, Stay, Focus

StareDesk adalah sistem IoT personal untuk memantau dan menganalisis sesi belajar secara realtime menggunakan kombinasi sensor pada meja belajar.

Project ini menghubungkan ESP32, MQTT Broker, Backend Golang, Database PostgreSQL, dan Frontend NextJS untuk menghasilkan monitoring kondisi belajar berbasis data nyata.

---

## Fitur Utama

- Monitoring kondisi belajar realtime
- Deteksi kehadiran otomatis
- Tracking sesi belajar otomatis
- Analytics produktivitas
- Realtime dashboard menggunakan WebSocket
- Komunikasi device menggunakan MQTT
- Threshold sensor dapat diubah langsung dari website

---

## Stack

| Layer | Teknologi |
|---|---|
| Frontend | NextJS + TypeScript |
| Backend | Golang (Gin) |
| Database | PostgreSQL (Supabase) |
| Realtime | WebSocket |
| IoT Communication | MQTT |
| Firmware | ESP32 + PlatformIO |

---

## Sensor yang Digunakan

| Sensor | Fungsi |
|---|---|
| PIR HC-SR501 | Deteksi kehadiran |
| HC-SR04 | Deteksi jarak pengguna |
| LDR | Deteksi intensitas cahaya |

---

## Arsitektur Sistem

```text
ESP32
  │
  ▼
MQTT Broker
  │
  ▼
Golang Backend
  ├── PostgreSQL
  ├── WebSocket
  └── REST API
          │
          ▼
NextJS Frontend
