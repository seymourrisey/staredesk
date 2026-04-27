#include <Arduino.h>
#include "sensors/pir.h"
#include "sensors/ultrasonic.h"
#include "sensors/ldr.h"
#include "sensors/condition.h"
#include "utils/moving_average.h"
#include "config/config.h"

MovingAverage distanceAvg(MOVING_AVG_WINDOW);

// Threshold aktif (akan di-override via MQTT config nanti)
int distanceMin = DEFAULT_DISTANCE_MIN_CM;
int distanceMax = DEFAULT_DISTANCE_MAX_CM;
int ldrThreshold = DEFAULT_LDR_THRESHOLD;

void setup() {
    Serial.begin(115200);
    pirSetup();
    ultrasonicSetup();
    ldrSetup();
    Serial.println("[StareDesk] Boot OK");
}

void loop() {
    bool  pir     = pirRead();
    float rawDist = ultrasonicRead();
    int   ldr     = ldrRead();

    distanceAvg.add(rawDist);
    float distance = distanceAvg.get();

    Condition condition = evaluateCondition(
        pir, distance, ldr,
        distanceMin, distanceMax, ldrThreshold
    );

    Serial.println("─────────────────────────");
    Serial.printf("PIR      : %s\n", pir ? "DETECTED" : "EMPTY");
    Serial.printf("Distance : %.1f cm\n", distance);
    Serial.printf("LDR      : %d\n", ldr);
    Serial.printf("Condition: %s\n", conditionToString(condition));

    delay(SENSOR_SAMPLE_INTERVAL_MS);
}
