#pragma once

#include <Arduino.h>
#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>

void mqttSetup();
void mqttLoop();
void mqttPublishTelemetry(bool pir, float distance, int ldr, const char* condition);
void mqttPublishStatus(bool isOnline);
void mqttPublishConfigAck();

extern int mqttDistanceMin;
extern int mqttDistanceMax;
extern int mqttLdrThreshold;
extern int mqttAwayTimeout;
