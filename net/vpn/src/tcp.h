#pragma once

#include <cstdint>

#include "ip.h"

namespace net_stack
{
    namespace tcp
    {

        enum TCP_STATE : uint8_t
        {
            CLOSED,
            LISTEN,
            SYN_RECEIVED,
            SYN_SENT,
            ESTABLISHED,
            FIN_WAIT_1,
            FIN_WAIT_2,
            CLOSING,
            TIME_WAIT,
            CLOSE_WAIT,
            LAST_ACK,
        };

#pragma pack(push)
#pragma pack(1)
        struct tcp_header_t
        {
            uint16_t src_port;
            uint16_t dst_port;
            uint32_t seq;
            uint32_t ack;
            uint8_t data_offset;
            uint8_t flags;
            uint16_t window_size;
            uint16_t checksum;
            uint16_t urgent_p;
            uint8_t opt[0];
        };
#pragma pack(pop)

        struct tcp_sock_t : public net_stack::ip::ip_sock_t
        {
            uint16_t src_port, dst_port;
            TCP_STATE state;
        };

    }

}