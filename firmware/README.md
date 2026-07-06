# StareDesk - Firmware

Firmware ESP32 untuk StareDesk. Sistem monitoring belajar berbasis IoT. Membaca tiga sensor secara berkala, mengevaluasi kondisi belajar, dan mengirimkan telemetry ke backend via MQTT over TLS.

---

## Hardware

| Komponen | Fungsi |
|---|---|
| ESP32 DOIT DevKit V1 (30 Pin, USB-C) | Mikrokontroler utama |
| PIR HC-SR501 | Deteksi kehadiran/gerakan |
| Ultrasonic HC-SR04 | Pengukuran jarak ke pengguna |
| LDR Module 4 Pin | Deteksi intensitas cahaya |

### Wiring

| Sensor | Pin ESP32 |
|---|---|
| HC-SR04 TRIG | GPIO 5 |
| HC-SR04 ECHO | GPIO 18 |
| PIR Signal | GPIO 19 |
| LDR Analog Out | GPIO 34 |

> HC-SR04 dipasang di monitor menghadap wajah/dada pengguna.

---

## Project Structure

```
firmware/
├── src/
│   ├── main.cpp                    # Entry point — setup & loop utama
│   ├── config/
│   │   ├── config.h                # Pin, threshold default, timing constant
│   │   ├── credentials.h           # WiFi & MQTT credentials (gitignored)
│   │   └── credentials-example.h  # Template credentials
│   ├── sensors/
│   │   ├── pir.cpp / pir.h         # Baca PIR HC-SR501
│   │   ├── ultrasonic.cpp / .h     # Baca HC-SR04
│   │   ├── ldr.cpp / ldr.h         # Baca LDR
│   │   └── condition.cpp / .h      # Evaluasi kondisi dari nilai sensor
│   ├── mqtt/
│   │   ├── client.cpp / client.h   # Koneksi MQTT, publish, subscribe
│   │   └── topics.h                # Topic string builder
│   └── utils/
│       └── moving_average.cpp / .h # Circular buffer moving average
├── platformio.ini
└── README.md
```

---

## Sensor Logic & Kondisi

Evaluasi kondisi dilakukan **di ESP32** sebelum dikirim ke backend, dengan urutan prioritas:

```
PIR = false
  → away

PIR = true + jarak < distance_min_cm
  → posture_risk

PIR = true + jarak > distance_max_cm
  → distracted

PIR = true + jarak normal + LDR < ldr_threshold
  → eye_strain_risk

PIR = true + jarak normal + LDR >= ldr_threshold
  → optimal
```

### Threshold Default

| Parameter | Default | Keterangan |
|---|---|---|
| `distance_min_cm` | 40 cm | Batas minimum — terlalu dekat |
| `distance_max_cm` | 90 cm | Batas maksimum — terlalu jauh |
| `ldr_threshold` | 500 | Di bawah nilai ini = cahaya kurang |
| `away_timeout_minutes` | 5 menit | Digunakan backend untuk end session |

Semua threshold dapat di-override secara real-time via MQTT config topic (tidak perlu flash ulang).

### Moving Average

Nilai jarak dari HC-SR04 diproses dengan moving average window 5 sampel sebelum dievaluasi, untuk meredam noise dari refleksi baju atau gerakan minor.

- Sample interval: `2000ms`
- Window size: `5`
- Evaluasi cycle efektif: ~10 detik untuk window penuh

---

## MQTT

### Broker

HiveMQ Cloud — koneksi TLS pada port `8883`.

### Topic Structure

```
ESP32 → Broker (Publish)
  study/{user_id}/device/telemetry      # Sensor data + kondisi
  study/{user_id}/device/status         # Online/offline (retain=true)
  study/{user_id}/device/config/ack     # Konfirmasi config diterima

Broker → ESP32 (Subscribe)
  study/{user_id}/device/config         # Threshold baru dari backend
```

### QoS

| Topic | QoS |
|---|---|
| `telemetry` | 0 — data terus mengalir, loss tidak kritis |
| `status` | 1 — harus sampai |
| `config` | 1 — harus sampai ke ESP32 |
| `config/ack` | 1 — harus sampai ke backend |

### LWT (Last Will & Testament)

Didaftarkan saat boot. Jika koneksi terputus mendadak (power cut, crash), broker otomatis publish:

```json
Topic   : study/{user_id}/device/status
Payload : { "is_online": false }
QoS     : 1
Retain  : true
```

### Telemetry Payload

```json
{
  "pir_detected": true,
  "distance_cm": 65.3,
  "ldr_value": 720,
  "condition": "optimal"
}
```

Dikirim saat:
- Kondisi berubah dari sebelumnya → **publish segera**
- Heartbeat interval 30 detik → **publish rutin**

---

## Boot Flow

```
Boot
├── Serial.begin(115200)
├── Setup semua sensor (PIR, Ultrasonic, LDR)
├── mqttSetup()
│   ├── Connect WiFi
│   ├── Connect MQTT Broker (TLS)
│   ├── Register LWT
│   ├── Subscribe study/{user_id}/device/config
│   └── Publish device/status → { is_online: true }
└── Loop
    ├── mqttLoop() — maintain koneksi & proses incoming
    ├── Baca sensor setiap SENSOR_SAMPLE_INTERVAL_MS
    ├── Update moving average distance
    ├── Evaluasi kondisi
    └── Kondisi berubah ATAU heartbeat due? → Publish telemetry
```

---

## Setup & Flash

### 1. Clone & buka di PlatformIO

```bash
git clone <repo-url>
cd firmware
```

### 2. Buat credentials.h

Salin dari template:

```bash
cp src/config/credentials-example.h src/config/credentials.h
```

Isi dengan nilai aktual:

```cpp
#define WIFI_SSID       "your_ssid"
#define WIFI_PASSWORD   "your_password"

#define MQTT_BROKER     "xxxx.hivemq.cloud"
#define MQTT_PORT       8883
#define MQTT_USER       "your_mqtt_user"
#define MQTT_PASSWORD   "your_mqtt_password"
#define MQTT_CLIENT_ID  "staredesk-esp32"
```

> `credentials.h` sudah ada di `.gitignore`. Jangan commit file ini.

### 3. Set USER_ID

Di `src/config/config.h`, sesuaikan USER_ID dengan ID user di database:

```cpp
#define USER_ID "USR-YYYYMMDD-XXXXXXXX"
```

### 4. Flash

```bash
pio run --target upload
```

Monitor serial output:

```bash
pio device monitor
```

---

## Dependencies

Didefinisikan di `platformio.ini`:

```ini
lib_deps =
    knolleary/PubSubClient @ ^2.8
    bblanchon/ArduinoJson @ ^7.4.1
```

| Library | Fungsi |
|---|---|
| PubSubClient | MQTT client untuk ESP32 |
| ArduinoJson | Serialize/deserialize payload JSON |

---

## Catatan Penting

- **`credentials.h` tidak boleh di-commit** — selalu gunakan `credentials-example.h` sebagai referensi.
- **`MQTT_CLIENT_ID` harus unik** dari Client ID yang digunakan Golang backend untuk menghindari kick conflict di broker.
- Jika HC-SR04 tidak mendapat echo (no object / terlalu jauh), `ultrasonicRead()` mengembalikan `-1.0f`. Nilai ini diabaikan oleh moving average dan kondisi fallback ke `away` via PIR check.
- LDR membaca `4095 - analogRead()` — nilai tinggi berarti cahaya terang, nilai rendah berarti gelap.
