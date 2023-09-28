#include <cstdio>
#include <cstdlib>
#include <unistd.h>
#include <memory.h>

#include <sys/stat.h>
#include <sys/types.h>
#include <sys/ioctl.h>
#include <fcntl.h>
#include <netinet/in.h>
#include <net/if.h>
#include <linux/if_tun.h>

#include "tun.h"

namespace net_stack
{
    tun_t::tun_t() : path(""), name(""), fd(-1)
    {
    }

    int tun_t::open(const std::string &path, const std::string &name)
    {
        this->path = path;
        this->name = name;

        if (name.size() > IFNAMSIZ)
        {
            return -1;
        }

        fd = ::open(path.c_str(), O_RDWR);
        if (fd == -1)
        {
            return -1;
        }

        struct ifreq ifr;
        memset(&ifr, 0, sizeof(ifr));
        ifr.ifr_flags = IFF_TUN | IFF_NO_PI;
        strncpy(ifr.ifr_name, name.c_str(), IFNAMSIZ);

        return ioctl(fd, TUNSETIFF, &ifr);
    }

    int tun_t::mtu()
    {
        return 2048;
    }

    int tun_t::close()
    {
        return ::close(fd);
    }

    int tun_t::read_from_tun(char *data, size_t size)
    {
        return ::read(fd, data, size);
    }

    int tun_t::write_to_tun(char *data, size_t size)
    {
        return ::write(fd, data, size);
    }

    int tun_t::write(std::shared_ptr<sk_buffer_t> p_sk_buffer)
    {
        output_queue.push_back(p_sk_buffer);
    }

    int tun_t::read(std::shared_ptr<sk_buffer_t> &p_sk_buffer)
    {
        if (input_queue.empty())
        {
            p_sk_buffer = nullptr;
            return -1;
        }

        p_sk_buffer = input_queue.front();
        input_queue.pop_front()

            return 0;
    }

    int tun_t::process_input()
    {
        std::shared_ptr<sk_buffer_t> p_sk_buff = std::make_shared<sk_buffer_t>(mtu());

    }

    int tun_t::process_output()
    {
    }

}
