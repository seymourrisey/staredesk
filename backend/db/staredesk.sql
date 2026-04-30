-- StareDesk Database Schema
-- Run this in Supabase SQL Editor

-- 1. users
CREATE TABLE users (
    id          VARCHAR PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- 2. device_config
CREATE TABLE device_config (
    id                   VARCHAR PRIMARY KEY,
    user_id              VARCHAR NOT NULL REFERENCES users(id),
    distance_min_cm      INT DEFAULT 40,
    distance_max_cm      INT DEFAULT 90,
    ldr_threshold        INT DEFAULT 500,
    away_timeout_minutes INT DEFAULT 5,
    config_ack           BOOLEAN DEFAULT FALSE,
    updated_at           TIMESTAMP DEFAULT NOW()
);

-- 3. device_status
CREATE TABLE device_status (
    id        VARCHAR PRIMARY KEY,
    user_id   VARCHAR NOT NULL REFERENCES users(id),
    is_online BOOLEAN DEFAULT FALSE,
    last_seen TIMESTAMP
);

-- 4. sensor_logs
CREATE TABLE sensor_logs (
    id           VARCHAR PRIMARY KEY,
    user_id      VARCHAR NOT NULL REFERENCES users(id),
    distance_cm  FLOAT,
    ldr_value    INT,
    pir_detected BOOLEAN,
    condition    TEXT CHECK (condition IN ('optimal', 'eye_strain_risk', 'posture_risk', 'distracted', 'away')),
    log_type     TEXT CHECK (log_type IN ('heartbeat', 'condition_change')),
    recorded_at  TIMESTAMP DEFAULT NOW()
);

-- 5. sessions
CREATE TABLE sessions (
    id                 VARCHAR PRIMARY KEY,
    user_id            VARCHAR NOT NULL REFERENCES users(id),
    started_at         TIMESTAMP,
    ended_at           TIMESTAMP,
    duration_sec       INT,
    dominant_condition TEXT CHECK (dominant_condition IN ('optimal', 'eye_strain_risk', 'posture_risk', 'distracted', 'away'))
);
