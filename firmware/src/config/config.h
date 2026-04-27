#pragma once

#include "credentials.h"

// ─── Pin Definitions ───────────────────────────────────────────────
#define TRIG_PIN    5
#define ECHO_PIN    18
#define PIR_PIN     19
#define LDR_PIN     34

// ─── Threshold Defaults ────────────────────────────────────────────
#define DEFAULT_DISTANCE_MIN_CM     40
#define DEFAULT_DISTANCE_MAX_CM     90
#define DEFAULT_LDR_THRESHOLD       500
#define DEFAULT_AWAY_TIMEOUT_MIN    5

// ─── Sampling & Heartbeat ──────────────────────────────────────────
#define SENSOR_SAMPLE_INTERVAL_MS   2000
#define HEARTBEAT_INTERVAL_MS       30000

// ─── Moving Average ────────────────────────────────────────────────
#define MOVING_AVG_WINDOW           5

// ─── User ──────────────────────────────────────────────────────────
#define USER_ID     "USR-20260424-TESTUSER"
