#include "client.h"
#include "topics.h"
#include "../config/config.h"
#include <ArduinoJson.h>

// Threshold yang bisa di-override via MQTT config
int mqttDistanceMin  = DEFAULT_DISTANCE_MIN_CM;
int mqttDistanceMax  = DEFAULT_DISTANCE_MAX_CM;
int mqttLdrThreshold = DEFAULT_LDR_THRESHOLD;
int mqttAwayTimeout  = DEFAULT_AWAY_TIMEOUT_MIN;

static WiFiClientSecure _wifiClient;
static PubSubClient     _mqtt(_wifiClient);

// ─── Config callback ───────────────────────────────────────────────
static void onConfigReceived(const char* payload) {
    JsonDocument doc;
    DeserializationError err = deserializeJson(doc, payload);
    if (err) {
        Serial.printf("[MQTT] Config parse error: %s\n", err.c_str());
        return;
    }

    if (doc["distance_min_cm"].is<int>())      mqttDistanceMin  = doc["distance_min_cm"];
    if (doc["distance_max_cm"].is<int>())      mqttDistanceMax  = doc["distance_max_cm"];
    if (doc["ldr_threshold"].is<int>())        mqttLdrThreshold = doc["ldr_threshold"];
    if (doc["away_timeout_minutes"].is<int>()) mqttAwayTimeout  = doc["away_timeout_minutes"];

    Serial.printf("[MQTT] Config updated — min:%d max:%d ldr:%d away:%d\n",
        mqttDistanceMin, mqttDistanceMax, mqttLdrThreshold, mqttAwayTimeout);

    mqttPublishConfigAck();
}

// ─── Callback ──────────────────────────────────────────────────────
static void mqttCallback(char* topic, byte* payload, unsigned int length) {
    String t = String(topic);
    String p = "";
    for (unsigned int i = 0; i < length; i++) p += (char)payload[i];

    if (t == topicConfig()) {
        onConfigReceived(p.c_str());
    }
}

// ─── WiFi ──────────────────────────────────────────────────────────
static void connectWifi() {
    if (WiFi.status() == WL_CONNECTED) return;

    Serial.printf("[WiFi] Connecting to %s", WIFI_SSID);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }
    Serial.printf("\n[WiFi] Connected — IP: %s\n", WiFi.localIP().toString().c_str());
}

// ─── MQTT Connect ──────────────────────────────────────────────────
static void connectMqtt() {
    while (!_mqtt.connected()) {
        Serial.print("[MQTT] Connecting...");

        // LWT payload
        String lwtPayload = "{\"is_online\":false}";

        bool connected = _mqtt.connect(
            MQTT_CLIENT_ID,
            MQTT_USER,
            MQTT_PASSWORD,
            topicStatus().c_str(),   // LWT topic
            1,                        // LWT QoS
            true,                     // LWT retain
            lwtPayload.c_str()        // LWT payload
        );

        if (connected) {
            Serial.println(" connected.");
            _mqtt.subscribe(topicConfig().c_str(), 1);
            mqttPublishStatus(true);
        } else {
            Serial.printf(" failed (rc=%d), retry in 5s\n", _mqtt.state());
            delay(5000);
        }
    }
}

// ─── Setup ─────────────────────────────────────────────────────────
void mqttSetup() {
    connectWifi();

    _wifiClient.setInsecure(); // HiveMQ Cloud TLS tanpa verify cert
    _mqtt.setServer(MQTT_BROKER, MQTT_PORT);
    _mqtt.setCallback(mqttCallback);
    _mqtt.setBufferSize(512);

    connectMqtt();
}

// ─── Loop ──────────────────────────────────────────────────────────
void mqttLoop() {
    if (WiFi.status() != WL_CONNECTED) connectWifi();
    if (!_mqtt.connected()) connectMqtt();
    _mqtt.loop();
}

// ─── Publish Telemetry ─────────────────────────────────────────────
void mqttPublishTelemetry(bool pir, float distance, int ldr, const char* condition) {
    JsonDocument doc;
    doc["pir_detected"] = pir;
    doc["distance_cm"]  = round(distance * 10.0f) / 10.0f;
    doc["ldr_value"]    = ldr;
    doc["condition"]    = condition;

    char buf[256];
    serializeJson(doc, buf);
    _mqtt.publish(topicTelemetry().c_str(), buf, false);

    Serial.printf("[MQTT] Telemetry published: %s\n", buf);
}

// ─── Publish Status ────────────────────────────────────────────────
void mqttPublishStatus(bool isOnline) {
    JsonDocument doc;
    doc["is_online"] = isOnline;

    char buf[64];
    serializeJson(doc, buf);
    _mqtt.publish(topicStatus().c_str(), buf, true); // retain = true
    Serial.printf("[MQTT] Status published: %s\n", buf);
}

// ─── Publish Config Ack ────────────────────────────────────────────
void mqttPublishConfigAck() {
    JsonDocument doc;
    doc["ack"] = true;

    char buf[32];
    serializeJson(doc, buf);
    _mqtt.publish(topicConfigAck().c_str(), buf, false);
    Serial.println("[MQTT] Config ACK published");
}
