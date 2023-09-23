#include <string>

#include "ip.h"
#include "tcp.h"
#include "tun.h"
#include "buffer.h"

namespace agent
{
    constexpr size_t BUFFER_SIZE = 4 * 1024;

    struct agent_t
    {
        dev::tun_t tun;
        util::Buffer<BUFFER_SIZE> tun_buffer;

        int open_tun(const std::string &path, const std::string &name);
        int read_tun();
        int handle_tun_input();
    };

}
