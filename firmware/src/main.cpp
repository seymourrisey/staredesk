#include <Arduino.h>
#include "sensors/pir.h"
#include "sensors/ultrasonic.h"
#include "sensors/ldr.h"
#include "utils/moving_average.h"
#include "config/config.h"

MovingAverage distanceAvg(MOVING_AVG_WINDOW);

void setup() {
    Serial.begin(115200);
    pirSetup();
    ultrasonicSetup();
    ldrSetup();
    Serial.println("[StareDesk] Boot OK");
}

void loop() {
    bool  pir      = pirRead();
    float rawDist  = ultrasonicRead();
    int   ldr      = ldrRead();

    distanceAvg.add(rawDist);
    float distance = distanceAvg.get();

    Serial.println("─────────────────────────");
    Serial.printf("PIR      : %s\n", pir ? "DETECTED" : "EMPTY");
    Serial.printf("Distance : %.1f cm (raw: %.1f)\n", distance, rawDist);
    Serial.printf("LDR      : %d\n", ldr);
    Serial.printf("MA Ready : %s\n", distanceAvg.isReady() ? "YES" : "warming up...");

    delay(SENSOR_SAMPLE_INTERVAL_MS);
}
