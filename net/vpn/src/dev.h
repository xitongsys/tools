#include <cstdint>
#include <memory>

#include "sk_buffer.h"

namespace net_stack
{
    struct dev_t
    {
        virtual int32_t mtu() = 0;
        virtual int write(std::shared_ptr<sk_buffer_t> p_sk_buffer) = 0;
        virtual int read(std::shared_ptr<sk_buffer_t> &psk_buffer) = 0;

        virtual int procee_input() = 0;
        virtual int process_output() = 0;
    };

}