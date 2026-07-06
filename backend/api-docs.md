# StareDesk API Documentation
> For frontend developers — Backend by seymourrisey

**Base URL (local):** `http://localhost:8080`  
**Base URL (production):** TBD

---

## Table of Contents
1. [Authentication](#authentication)
2. [Auth Endpoints](#auth-endpoints)
3. [Device Endpoints](#device-endpoints)
4. [Session Endpoints](#session-endpoints)
5. [Analytics Endpoints](#analytics-endpoints)
6. [Sensor Log Endpoints](#sensor-log-endpoints)
7. [WebSocket](#websocket)
8. [Data Types Reference](#data-types-reference)
9. [Error Responses](#error-responses)

---

## Authentication

Semua endpoint kecuali `/auth/login` memerlukan autentikasi.

JWT disimpan sebagai **httpOnly cookie** bernama `token`. Browser otomatis mengirim cookie di setiap request — tidak perlu set header manual.

**Untuk fetch requests, wajib tambahkan:**
```typescript
fetch(url, {
  credentials: 'include', // wajib untuk httpOnly cookie
})
```

**Untuk axios, set global config:**
```typescript
axios.defaults.withCredentials = true;
```

**WebSocket** menggunakan query param `?token=` karena browser tidak support custom header di WebSocket. Baca token dari endpoint `/auth/me` atau simpan sementara di memory saat login.

---

## Auth Endpoints

### POST `/auth/login`
Login dan set httpOnly cookie.

**Request:**
```json
{
  "email": "admin@staredesk.com",
  "password": "password"
}
```

**Response 200:**
```json
{
  "user": {
    "id": "USR-20260424-TESTUSER",
    "email": "admin@staredesk.com"
  }
}
```
Cookie `token` otomatis di-set oleh browser. JWT berlaku **24 jam**.

**Response 401:**
```json
{ "error": "invalid email or password" }
```

---

### POST `/auth/logout`
Clear cookie `token`.

**Response 200:**
```json
{ "message": "logged out" }
```

---

### GET `/auth/me`
Validasi token dan return info user yang sedang login.

**Response 200:**
```json
{ "user_id": "USR-20260424-TESTUSER" }
```

**Response 401:** Token tidak valid atau expired.

> Gunakan endpoint ini untuk: cek apakah user sudah login, get `user_id` untuk keperluan lain.

---

## Device Endpoints

### GET `/device/status`
Status online/offline device ESP32.

**Response 200:**
```json
{
  "id": "DST-20260424-XXXXXXXX",
  "user_id": "USR-20260424-TESTUSER",
  "is_online": true,
  "last_seen": "2026-04-30T03:50:25Z"
}
```

> `last_seen` adalah UTC timestamp. Konversi ke WIB di frontend:
> ```typescript
> new Date(last_seen).toLocaleString('id-ID', { timeZone: 'Asia/Jakarta' })
> ```

---

### GET `/device/config`
Ambil threshold konfigurasi aktif + status ack.

**Response 200:**
```json
{
  "id": "CFG-20260424-XXXXXXXX",
  "user_id": "USR-20260424-TESTUSER",
  "distance_min_cm": 40,
  "distance_max_cm": 90,
  "ldr_threshold": 500,
  "away_timeout_minutes": 3,
  "config_ack": true,
  "updated_at": "2026-04-30T03:50:25Z"
}
```

**Field `config_ack`:**
- `false` → config sudah disimpan di DB, menunggu ESP32 konfirmasi → tampilkan **"⏳ Pending..."**
- `true` → ESP32 sudah menerima dan apply config → tampilkan **"✅ Confirmed"**

---

### PUT `/device/config`
Update threshold konfigurasi. Otomatis publish ke ESP32 via MQTT.

**Request:**
```json
{
  "distance_min_cm": 40,
  "distance_max_cm": 90,
  "ldr_threshold": 500,
  "away_timeout_minutes": 3
}
```

**Response 200:**
```json
{
  "id": "CFG-20260424-XXXXXXXX",
  "distance_min_cm": 40,
  "distance_max_cm": 90,
  "ldr_threshold": 500,
  "away_timeout_minutes": 3,
  "config_ack": false,
  "updated_at": "2026-04-30T04:00:00Z"
}
```

> Setelah PUT, `config_ack` akan `false`. Tunggu WebSocket event `config_ack` untuk update UI menjadi "Confirmed".

---

## Session Endpoints

### GET `/sessions`
List semua sesi selesai dengan pagination.

**Query params:**
| Param | Type | Default | Keterangan |
|---|---|---|---|
| `limit` | int | 20 | Max 100 |
| `offset` | int | 0 | Untuk pagination |

**Response 200:**
```json
{
  "sessions": [
    {
      "id": "SSN-20260430-XXXXXXXX",
      "user_id": "USR-20260424-TESTUSER",
      "started_at": "2026-04-30T03:23:02Z",
      "ended_at": "2026-04-30T03:28:09Z",
      "duration_sec": 306,
      "dominant_condition": "optimal"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

> Sesi yang masih aktif (`ended_at = null`) tidak muncul di sini.

---

### GET `/sessions/summary`
Ringkasan sesi untuk range tertentu.

**Query params:**
| Param | Value | Keterangan |
|---|---|---|
| `range` | `today` \| `week` \| `month` | Default: `today` |

**Response 200:**
```json
{
  "range": "today",
  "total_sec": 306,
  "session_count": 1,
  "sessions": [
    {
      "id": "SSN-20260430-XXXXXXXX",
      "started_at": "2026-04-30T03:23:02Z",
      "ended_at": "2026-04-30T03:28:09Z",
      "duration_sec": 306,
      "dominant_condition": "optimal"
    }
  ]
}
```

**Cara hitung VS Last Period di frontend:**
```typescript
// Fetch dua periode, hitung selisih
const today = await fetch('/sessions/summary?range=today')
const yesterday = await fetch('/sessions/summary?range=today', {
  // tambahkan date param untuk hari kemarin — lihat catatan di bawah
})
const diff = today.total_sec - yesterday.total_sec
const pct = ((diff / yesterday.total_sec) * 100).toFixed(0)
```

> ⚠️ Endpoint summary saat ini tidak support `date` param — untuk VS Last Period, bisa derive dari `/sessions?limit=100` dan filter by date di frontend.

---

### GET `/sessions/:id`
Detail satu sesi.

**Response 200:**
```json
{
  "id": "SSN-20260430-XXXXXXXX",
  "user_id": "USR-20260424-TESTUSER",
  "started_at": "2026-04-30T03:23:02Z",
  "ended_at": "2026-04-30T03:28:09Z",
  "duration_sec": 306,
  "dominant_condition": "optimal"
}
```

**Response 404:**
```json
{ "error": "session not found" }
```

---

## Analytics Endpoints

### GET `/analytics/peak-hours`
Jam paling produktif berdasarkan total durasi sesi.

**Query params:**
| Param | Value | Keterangan |
|---|---|---|
| `range` | `week` \| `month` | **Wajib** — tidak support `today` |

**Response 200:**
```json
{
  "range": "week",
  "entries": [
    { "Hour": 10, "TotalSec": 3600, "SessionCount": 2 },
    { "Hour": 14, "TotalSec": 2700, "SessionCount": 1 },
    { "Hour": 9,  "TotalSec": 1800, "SessionCount": 1 }
  ]
}
```

> `entries` diurutkan by `TotalSec` descending — entry pertama adalah jam paling produktif.
> Untuk tampilkan "Top 3 Peak Hours", ambil 3 entry pertama.
> Untuk Daily view, **jangan tampilkan card ini** — endpoint tidak support daily.

---

### GET `/analytics/condition-breakdown`
Proporsi tiap kondisi dalam range tertentu.

**Query params:**
| Param | Value | Keterangan |
|---|---|---|
| `range` | `today` \| `week` \| `month` | Default: `today` |

**Response 200:**
```json
{
  "range": "today",
  "entries": [
    { "Condition": "optimal",         "Count": 120 },
    { "Condition": "eye_strain_risk", "Count": 30  },
    { "Condition": "posture_risk",    "Count": 25  },
    { "Condition": "distracted",      "Count": 15  },
    { "Condition": "away",            "Count": 40  }
  ]
}
```

> `Count` = jumlah sensor log records dengan kondisi tersebut.
> **Jangan tampilkan `away` di Condition Breakdown chart** — filter di frontend sebelum render.
> Untuk hitung persentase: `(entry.Count / totalCount) * 100` (exclude away dari total).

---

### GET `/analytics/timeline`
Data timeline per jam untuk satu hari.

**Query params:**
| Param | Format | Keterangan |
|---|---|---|
| `date` | `YYYY-MM-DD` | Default: hari ini (UTC) |

**Response 200:**
```json
{
  "date": "2026-04-30",
  "entries": [
    { "Hour": 0,  "DominantCondition": "",        "TotalSec": 0    },
    { "Hour": 1,  "DominantCondition": "",        "TotalSec": 0    },
    { "Hour": 10, "DominantCondition": "optimal", "TotalSec": 3600 },
    { "Hour": 11, "DominantCondition": "optimal", "TotalSec": 2400 },
    { "Hour": 14, "DominantCondition": "distracted", "TotalSec": 1800 }
  ]
}
```

> Response selalu 24 entries (jam 0–23).
> Jam dengan `TotalSec: 0` dan `DominantCondition: ""` = tidak ada aktivitas di jam tersebut.
> Gunakan `DominantCondition` untuk warna bar chart.

---

## Sensor Log Endpoints

### GET `/sensor-logs`
Raw sensor data — untuk debug atau advanced analytics.

**Query params:**
| Param | Format | Keterangan |
|---|---|---|
| `from` | RFC3339 | Contoh: `2026-04-30T00:00:00Z` |
| `to` | RFC3339 | Contoh: `2026-04-30T23:59:59Z` |
| `limit` | int | Default: 100, max: 1000 |

**Response 200:**
```json
{
  "from": "2026-04-30T00:00:00Z",
  "to": "2026-04-30T23:59:59Z",
  "limit": 100,
  "logs": [
    {
      "id": "SLG-20260430-XXXXXXXX",
      "user_id": "USR-20260424-TESTUSER",
      "distance_cm": 41.9,
      "ldr_value": 1270,
      "pir_detected": true,
      "condition": "optimal",
      "log_type": "heartbeat",
      "recorded_at": "2026-04-30T03:23:02Z"
    }
  ]
}
```

---

## WebSocket

### GET `/ws?token=<JWT>`

WebSocket connection untuk realtime data stream.

**URL:** `ws://localhost:8080/ws?token=<JWT_TOKEN>`

Token didapat dari: simpan JWT di memory saat response login, atau fetch ulang via `/auth/me` — **tapi `/auth/me` tidak return token**, hanya `user_id`. Simpan token di memory (variabel JavaScript) saat login response sebelum cookie di-set.

**Cara connect di frontend:**
```typescript
// Simpan token saat login
let wsToken = ''

// Di login handler — ambil token dari response sebelum cookie handling
// Token masih tersedia di response body Login sebelum implementasi cookie

const ws = new WebSocket(`ws://localhost:8080/ws?token=${wsToken}`)

ws.onmessage = (event) => {
  const data = JSON.parse(event.data)
  // handle berdasarkan data.type
}
```

---

### Message Types

#### `telemetry`
Dikirim setiap 30 detik atau saat kondisi berubah.
```json
{
  "type": "telemetry",
  "timestamp": "2026-04-30T03:50:25Z",
  "device": {
    "is_online": true,
    "last_seen": "2026-04-30T03:50:25Z"
  },
  "sensors": {
    "distance_cm": 41.9,
    "ldr_value": 1270,
    "pir_detected": true
  },
  "condition": "optimal",
  "session": {
    "is_active": true
  }
}
```

#### `condition_change`
Struktur sama dengan `telemetry`, dikirim **segera** saat kondisi berubah (tidak tunggu 30 detik).
```json
{
  "type": "condition_change",
  ...sama dengan telemetry...
}
```

#### `session_start`
```json
{
  "type": "session_start",
  "timestamp": "2026-04-30T03:50:25Z",
  "started_at": "2026-04-30T03:50:25Z"
}
```

#### `session_end`
```json
{
  "type": "session_end",
  "timestamp": "2026-04-30T03:53:49Z"
}
```

#### `device_status`
```json
{
  "type": "device_status",
  "timestamp": "2026-04-30T03:45:05Z",
  "device": {
    "is_online": true,
    "last_seen": "2026-04-30T03:45:05Z"
  }
}
```

#### `config_ack`
```json
{
  "type": "config_ack",
  "timestamp": "2026-04-30T03:45:05Z"
}
```

---

## Data Types Reference

### Condition Values
| Value | Display Label | Warna UI |
|---|---|---|
| `optimal` | Optimal | Kuning/Amber |
| `eye_strain_risk` | Eye Strain Risk | Kuning tua |
| `posture_risk` | Posture Risk | Oranye |
| `distracted` | Distracted | Abu-abu gelap |
| `away` | Away | Abu-abu terang |

### Timestamp Format
Semua timestamp dari backend adalah **UTC dalam format RFC3339**: `2026-04-30T03:50:25Z`

Konversi ke WIB di frontend:
```typescript
const toWIB = (utc: string) =>
  new Date(utc).toLocaleString('id-ID', { timeZone: 'Asia/Jakarta' })
```

### Duration Format
`duration_sec` adalah integer dalam detik. Konversi ke display:
```typescript
const formatDuration = (sec: number) => {
  const h = Math.floor(sec / 3600)
  const m = Math.floor((sec % 3600) / 60)
  if (h > 0) return `${h}j ${m}m`
  return `${m}m`
}
```

---

## Error Responses

### 400 Bad Request
```json
{ "error": "Key: 'loginRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag" }
```

### 401 Unauthorized
```json
{ "error": "invalid email or password" }
```
atau saat token expired/invalid:
```json
{ "error": "unauthorized" }
```

### 404 Not Found
```json
{ "error": "session not found" }
```

### 500 Internal Server Error
```json
{ "error": "internal server error" }
```

---

## Quick Start untuk Frontend Dev

```typescript
// 1. Login
const res = await fetch('http://localhost:8080/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({ email: 'admin@staredesk.com', password: 'password' })
})
const { user } = await res.json()

// 2. Fetch data (cookie otomatis dikirim)
const status = await fetch('http://localhost:8080/device/status', {
  credentials: 'include'
})

// 3. Connect WebSocket (butuh token — simpan dari login response jika tersedia)
const ws = new WebSocket('ws://localhost:8080/ws?token=YOUR_TOKEN')
ws.onmessage = (e) => {
  const msg = JSON.parse(e.data)
  switch (msg.type) {
    case 'telemetry': // update live cards
    case 'session_start': // start timer
    case 'session_end': // stop timer, refresh summary
    case 'device_status': // update status bar
    case 'config_ack': // update settings "Confirmed"
  }
}

// 4. Logout
await fetch('http://localhost:8080/auth/logout', {
  method: 'POST',
  credentials: 'include'
})
```

---

*StareDesk API — Backend: Golang (Gin) + Supabase PostgreSQL + HiveMQ MQTT*  
*Last updated: 30 April 2026*
