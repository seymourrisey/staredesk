#pragma once

#include <stdint.h>

class MovingAverage {
public:
    MovingAverage(uint8_t windowSize);
    ~MovingAverage();

    void add(float value);
    float get() const;
    bool isReady() const;

private:
    float*  _buffer;
    uint8_t _windowSize;
    uint8_t _count;
    uint8_t _index;
    float   _sum;
};
