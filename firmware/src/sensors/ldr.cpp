#include "ldr.h"
#include "../config/config.h"
#include <Arduino.h>

void ldrSetup() {
    pinMode(LDR_PIN, INPUT);
}

int ldrRead() {
    return 4095 - analogRead(LDR_PIN);
}
