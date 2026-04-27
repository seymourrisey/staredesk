#include "pir.h"
#include "../config/config.h"
#include <Arduino.h>

void pirSetup() {
    pinMode(PIR_PIN, INPUT);
}

bool pirRead() {
    return digitalRead(PIR_PIN) == HIGH;
}
