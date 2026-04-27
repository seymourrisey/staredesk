#include <Arduino.h>
#include "sensors/pir.h"
#include "sensors/ultrasonic.h"
#include "sensors/ldr.h"
#include "sensors/condition.h"
#include "utils/moving_average.h"
#include "mqtt/client.h"
#include "config/config.h"

MovingAverage distanceAvg(MOVING_AVG_WINDOW);

Condition     lastCondition    = CONDITION_AWAY;
unsigned long lastHeartbeat    = 0;
unsigned long lastSampleTime   = 0;

void setup() {
    Serial.begin(115200);
    pirSetup();
    ultrasonicSetup();
    ldrSetup();
    mqttSetup();
    Serial.println("[StareDesk] Boot OK");
}

void loop() {
    mqttLoop();

    unsigned long now = millis();

    // Baca sensor setiap SENSOR_SAMPLE_INTERVAL_MS
    if (now - lastSampleTime < SENSOR_SAMPLE_INTERVAL_MS) return;
    lastSampleTime = now;

    bool  pir      = pirRead();
    float rawDist  = ultrasonicRead();
    int   ldr      = ldrRead();

    distanceAvg.add(rawDist);
    float distance = distanceAvg.get();

    Condition condition = evaluateCondition(
        pir, distance, ldr,
        mqttDistanceMin, mqttDistanceMax, mqttLdrThreshold
    );

    Serial.println("─────────────────────────");
    Serial.printf("PIR      : %s\n", pir ? "DETECTED" : "EMPTY");
    Serial.printf("Distance : %.1f cm\n", distance);
    Serial.printf("LDR      : %d\n", ldr);
    Serial.printf("Condition: %s\n", conditionToString(condition));

    bool conditionChanged = (condition != lastCondition);
    bool heartbeatDue     = (now - lastHeartbeat >= HEARTBEAT_INTERVAL_MS);

    if (conditionChanged || heartbeatDue) {
        mqttPublishTelemetry(pir, distance, ldr, conditionToString(condition));
        lastHeartbeat = now;
    }

    lastCondition = condition;
}
