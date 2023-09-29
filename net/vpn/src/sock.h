#include <cstdint>
#include <memory>
#include <list>

#include "dev.h"
#include "sk_buffer.h"

namespace net_stack
{
    struct sock_t
    {
        std::shared_ptr<dev_t> dev;
        std::list<sk_buffer_t> sk_buffers;
    };
}
