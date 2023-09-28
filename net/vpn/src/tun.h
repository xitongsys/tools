#pragma once

#include <list>
#include <string>

#include "dev.h"
#include "sk_buffer.h"

namespace net_stack
{
    struct tun_t : public dev_t
    {
        std::string path; // "/dev/net/tun"
        std::string name; // "tun0","tun1" ...
        int fd;

        std::list<std::shared_ptr<sk_buffer_t>> input_queue, output_queue;

        tun_t();

        int open(const std::string &path, const std::string &name);
        int close();

        int read_from_tun(char *data, size_t size);
        int write_to_tun(char *data, size_t size);

        int process_input();
        int process_output();
    };
}