#pragma once

enum Condition {
    CONDITION_OPTIMAL,
    CONDITION_EYE_STRAIN_RISK,
    CONDITION_POSTURE_RISK,
    CONDITION_DISTRACTED,
    CONDITION_AWAY
};

const char* conditionToString(Condition c);

Condition evaluateCondition(
    bool  pirDetected,
    float distance,
    int   ldrValue,
    int   distanceMin,
    int   distanceMax,
    int   ldrThreshold
);
