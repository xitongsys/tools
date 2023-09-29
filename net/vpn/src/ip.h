#pragma once

#include <cstdint>
#include <list>
#include <memory>

#include "env.h"
#include "sk_buffer.h"
#include "dev.h"
#include "sock.h"

namespace net_stack
{
    namespace ip
    {

        enum PROTOCOL : uint8_t
        {
            ICMP = 1,
            TCP = 6,
            UDP = 17,

        };

#pragma pack(push)
#pragma pack(1)
        struct ip_header_t
        {
            struct
            {
#if (IS_BIG_ENDIAN)
                uint8_t ver : 4;
                uint8_t ihl : 4;
#else
                uint8_t ihl : 4;
                uint8_t ver : 4;
#endif
            };

            uint8_t tos;
            uint16_t total_len;
            uint16_t id;
            uint16_t frag_off;
            uint8_t ttl;
            uint8_t protocol;
            uint16_t check;
            uint32_t saddr;
            uint16_t daddr;
            uint8_t opt[0];
        };
#pragma pack(pop)

        struct ip_sock_t : public sock_t
        {
            uint8_t ver;
            uint32_t src, dst;
            PROTOCOL protocol;
        };

    }
}