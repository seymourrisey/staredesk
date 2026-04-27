#pragma once

#include <Arduino.h>
#include "../config/config.h"

// Format: study/{user_id}/device/...
inline String topicTelemetry() {
    return String("study/") + USER_ID + "/device/telemetry";
}

inline String topicStatus() {
    return String("study/") + USER_ID + "/device/status";
}

inline String topicConfig() {
    return String("study/") + USER_ID + "/device/config";
}

inline String topicConfigAck() {
    return String("study/") + USER_ID + "/device/config/ack";
}
