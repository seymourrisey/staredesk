#include "ultrasonic.h"
#include "../config/config.h"
#include <Arduino.h>

void ultrasonicSetup() {
    pinMode(TRIG_PIN, OUTPUT);
    pinMode(ECHO_PIN, INPUT);
}

float ultrasonicRead() {
    digitalWrite(TRIG_PIN, LOW);
    delayMicroseconds(2);
    digitalWrite(TRIG_PIN, HIGH);
    delayMicroseconds(10);
    digitalWrite(TRIG_PIN, LOW);

    long duration = pulseIn(ECHO_PIN, HIGH, 30000); // timeout 30ms
    if (duration == 0) return -1.0f; // no echo = away

    return (duration * 0.0343f) / 2.0f;
}
