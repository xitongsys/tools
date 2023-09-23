#pragma once

#include <cstdint>
#include <tuple>
#include <cstddef>

namespace util
{
    template <int BUFFER_SIZE>
    struct Buffer
    {
        char buffer[BUFFER_SIZE];
        int read_pos, write_pos, end_pos;

        Buffer()
        {
            read_pos = 0;
            write_pos = 0;
            end_pos = BUFFER_SIZE;
        }

        size_t space()
        {
            if (read_pos > write_pos)
            {
                return (read_pos - write_pos - 1);
            }

            int size_right = BUFFER_SIZE - write_pos - 1;
            int size_left = read_pos - 1;

            if (size_left > size_right)
            {
                end_pos = write_pos;
                write_pos = 0;
                return size_left;
            }

            return size_right;
        }

        std::tuple<char *, size_t> write_buffer()
        {
            size_t space_size = space();
            return {buffer + write_pos, space_size};
        }

        int write_consume(size_t size)
        {
            if (write_pos + size >= BUFFER_SIZE)
            {
                return -1;
            }
            write_pos += size;

            return 0;
        }

        size_t avail()
        {
            if (read_pos <= write_pos)
            {
                return write_pos - read_pos;
            }

            if (read_pos < end_pos)
            {
                return end_pos - read_pos;
            }

            size_t size = write_pos;
            read_pos = 0;
            return size;
        }

        std::tuple<char *, size_t> read_buffer()
        {
            size_t avail_size = avail();
            return {buffer + read_pos, avail_size};
        }

        int read_consume(size_t size)
        {
            if (read_pos + size >= BUFFER_SIZE)
            {
                return -1;
            }

            read_pos += size;
            return 0;
        }
    };

}
