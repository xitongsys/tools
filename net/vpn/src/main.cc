#include <iostream>

#include "tcp.h"
#include "tun.h"
#include "ip.h"

int main()
{
    std::cout << sizeof(net_stack::tcp::tcp_header_t) << std::endl;

    dev::tun_t t;
    std::cout << t.tun_open("/dev/net/tun", "tun0") << std::endl;

    const int BS = 4 * 1024;
    char buffer[BS];

    for (;;)
    {
        net_stack::ip::ip_header_t *ip_header = (net_stack::ip::ip_header_t *)buffer;
        int cnt = t.tun_read(buffer, BS);

        std::cout << cnt << " " << (int)ip_header->ver << " " << (int)ip_header->ihl*5 << std::endl;
    }

    t.tun_close();

    return 0;
}