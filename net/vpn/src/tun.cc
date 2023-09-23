#include <sys/stat.h>
#include <sys/types.h>
#include <sys/ioctl.h>
#include <fcntl.h>
#include <netinet/in.h>
#include <net/if.h>
#include <linux/if_tun.h>

#include "tun.h"

namespace dev
{
    tun_t::tun_t(const std::string &path, const std::string &name) : path(path), name(name)
    {
    }

    int tun_t::tun_open()
    {
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
