#include <bits/stdc++.h>

namespace dev
{
    struct tun_t
    {
        std::string path; // "/dev/net/tun"
        std::string name; // "tun0","tun1" ...
        int fd;

        tun_t(const std::string &path, const std::string &name);
        
        int tun_open();
        int tun_close();
        int tun_read(char *buffer, int size);
        int tun_write(char *buffer, int size);
    };
}