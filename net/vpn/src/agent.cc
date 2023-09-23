#include "agent.h"
#include "ip.h"

namespace agent
{
    int agent_t::open_tun(const std::string &path, const std::string &name)
    {
        return tun.tun_open(path, name);
    }

    int agent_t::read_tun()
    {
        char *buffer = nullptr;
        size_t size = 0;
        std::tie(buffer, size) = tun_buffer.write_buffer();

        if (size <= 0)
        {
            return 0;
        }

        int read_size = tun.tun_read(buffer, size);
        if (read_size > 0)
        {
            tun_buffer.write_consume(read_size);
        }

        return read_size;
    }

    int agent_t::handle_tun_input()
    {
        char *buffer = nullptr;
        size_t size = 0;
        std::tie(buffer, size) = tun_buffer.read_buffer();

        if (size <= sizeof(net_stack::ip_header_t))
        {
            return 0;
        }

        net_stack::ip_header_t *ip_header = (net_stack::ip_header_t *)buffer;

        if (ip_header->ver == 4) // ipv4
        {
            if (ip_header->protocol == net_stack::ICMP) // icmp
            {
            }
            else if (ip_header->protocol == net_stack::TCP) // tcp
            {
                
            }
            else if (ip_header->protocol == net_stack::UDP) // udp
            {
            }
            else
            {
            }

        }
        else if (ip_header->ver == 6) // ipv6
        {
        }
        else
        {
        }
    }

}