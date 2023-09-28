#pragma once

#include <cstdint>

#include "env.h"

namespace net_stack
{
    namespace network_layer
    {
        namespace ip
        {

            enum Protocol : uint8_t
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
                    uint8_t ihl : 4;
                    uint8_t ver : 4;
                };

                struct
                {
                    uint8_t ecn : 2;
                    uint8_t dscp : 6;
                } tos;

                uint16_t total_length;
                uint16_t id;

                struct
                {
                    uint16_t offset : 13; // 13 bits fragment-offset
                    uint16_t flags : 3;   // 3 bits flags and
                };

                uint8_t ttl;
                uint8_t protocol;
                uint16_t checksum;
                uint32_t src_addr;
                uint16_t dst_addr;
                uint8_t opt[0];
            };
#pragma pack(pop)




            struct ip_socket_t
            {
                char buffer;
            };

        

        }
    }
}