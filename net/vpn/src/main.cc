#include <bits/stdc++.h>

#include "tcp.h"
#include "tun.h"

int main()
{
    std::cout << sizeof(net_stack::tcp_header_t) << std::endl;

    dev::tun_t t("/dev/net/tun", "tun1");
    std::cout << t.tun_open() << std::endl;
    t.tun_close();

    return 0;
}