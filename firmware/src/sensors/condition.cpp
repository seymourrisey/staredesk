#include "condition.h"

const char* conditionToString(Condition c) {
    switch (c) {
        case CONDITION_OPTIMAL:         return "optimal";
        case CONDITION_EYE_STRAIN_RISK: return "eye_strain_risk";
        case CONDITION_POSTURE_RISK:    return "posture_risk";
        case CONDITION_DISTRACTED:      return "distracted";
        case CONDITION_AWAY:            return "away";
        default:                        return "away";
    }
}

Condition evaluateCondition(
    bool  pirDetected,
    float distance,
    int   ldrValue,
    int   distanceMin,
    int   distanceMax,
    int   ldrThreshold
) {
    if (!pirDetected)                return CONDITION_AWAY;
    if (distance < distanceMin)      return CONDITION_POSTURE_RISK;
    if (distance > distanceMax)      return CONDITION_DISTRACTED;
    if (ldrValue < ldrThreshold)     return CONDITION_EYE_STRAIN_RISK;
    return CONDITION_OPTIMAL;
}
