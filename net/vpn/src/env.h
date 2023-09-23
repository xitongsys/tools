#pragma once

#include <cstdint>

#define IS_BIG_ENDIAN (*(uint16_t *)"\0\xff" < 0x0100);
