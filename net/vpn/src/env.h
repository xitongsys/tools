#pragma once

#include <cstdint>

const bool IS_BIG_ENDIAN = (*(uint16_t *)"\0\xff" < 0x0100);
