#include <cstdint>

#include "env.h"

namespace net_stack
{

    enum uint8_t
    {
        ICMP = 1,
        TCP = 6,
        UDP = 17,

    } IP_PROTOCOL;
    

#pragma pack(push)
#pragma pack(1)
    struct ip_header_t
    {
#if (IS_BIG_ENDIAN)
        struct
        {
            uint8_t ver : 4;
            uint8_t ihl : 4;
        };
#else
        struct
        {
            uint8_t ihl : 4;
            uint8_t ver : 4;
        };
#endif

        uint8_t tos;
        uint16_t total_length;
        uint16_t id;

#if (IS_BIG_ENDIAN)
        struct
        {
            uint16_t flags : 3;   // 3 bits flags and
            uint16_t offset : 13; // 13 bits fragment-offset
        };
#else
        struct
        {
            uint16_t offset : 13; // 13 bits fragment-offset
            uint16_t flags : 3;   // 3 bits flags and
        };
#endif

        uint8_t ttl;
        uint8_t protocol;
        uint16_t checksum;
        uint32_t src_addr;
        uint16_t dst_addr;
        uint8_t opt[0];
    };
#pragma pack(pop)
}