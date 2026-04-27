#include "moving_average.h"

MovingAverage::MovingAverage(uint8_t windowSize) {
    _windowSize = windowSize;
    _buffer     = new float[windowSize];
    _count      = 0;
    _index      = 0;
    _sum        = 0.0f;

    for (uint8_t i = 0; i < windowSize; i++) {
        _buffer[i] = 0.0f;
    }
}

MovingAverage::~MovingAverage() {
    delete[] _buffer;
}

void MovingAverage::add(float value) {
    if (value < 0) return; // abaikan no-echo (-1)

    _sum -= _buffer[_index];
    _buffer[_index] = value;
    _sum += value;

    _index = (_index + 1) % _windowSize;
    if (_count < _windowSize) _count++;
}

float MovingAverage::get() const {
    if (_count == 0) return -1.0f;
    return _sum / _count;
}

bool MovingAverage::isReady() const {
    return _count == _windowSize;
}
