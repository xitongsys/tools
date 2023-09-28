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
    namespace dev_layer
    {
        tun_t::tun_t() : path(""), name(""), fd(-1)
        {
        }

        int tun_t::tun_open(const std::string &path, const std::string &name)
        {
            this->path = path;
            this->name = name;

            if (name.size() > IFNAMSIZ)
            {
                return -1;
            }

            fd = open(path.c_str(), O_RDWR);
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

        int tun_t::tun_close()
        {
            return close(fd);
        }

        int tun_t::tun_read(char *buffer, int size)
        {
            return read(fd, buffer, size);
        }

        int tun_t::tun_write(char *buffer, int size)
        {
            return write(fd, buffer, size);
        }
    }

}
