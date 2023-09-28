#include <cstdint>

#include "sk_buffer.h"

namespace net_stack
{
    sk_buffer_t::sk_buffer_t()
    {
        size = 0;
        bgn = head = tail = end = nullptr;
    }

    void sk_buffer_t::alloc()
    {
        bgn = new uint8_t[size];
        end = head + size;
        head = bgn;
        tail = bgn;
    }

    void sk_buffer_t::free()
    {
        if (bgn)
        {
            delete bgn;
        }
    }

    int32_t sk_buffer_t::data_size()
    {
        return tail - head;
    }

    sk_buffer_t::sk_buffer_t(uint32_t size)
    {
        this->size = size;
        alloc();
    }

    sk_buffer_t::sk_buffer_t(const sk_buffer_t &sk_b)
    {
        free();
        size = sk_b.size;
        alloc();

        bgn += (sk_b.bgn - sk_b.head);
        tail += (sk_b.tail - sk_b.head);
    }

    sk_buffer_t::sk_buffer_t(sk_buffer_t &&sk_b)
    {
        free();
        size = sk_b.size;
        bgn = sk_b.bgn;
        head = sk_b.head;
        end = sk_b.end;
        tail = sk_b.tail;

        sk_b.bgn = nullptr;
        sk_b.head = nullptr;
        sk_b.tail = nullptr;
        sk_b.end = nullptr;
    }

    sk_buffer_t &sk_buffer_t::operator=(const sk_buffer_t &sk_b)
    {
        free();
        size = sk_b.size;
        alloc();

        bgn += (sk_b.bgn - sk_b.head);
        tail += (sk_b.tail - sk_b.head);
    }

    sk_buffer_t::~sk_buffer_t()
    {
        free();
    }

    int32_t sk_buffer_t::head_size()
    {
        return head - bgn;
    }

    int32_t sk_buffer_t::tail_size()
    {
        return end - tail;
    }

    uint8_t *sk_buffer_t::push_to_head(int32_t push_size)
    {
        if (head_size() < push_size)
        {
            return nullptr;
        }

        head -= push_size;
        return head;
    }

    uint8_t *sk_buffer_t::pop_from_head(int32_t pop_size)
    {
        if (data_size() < pop_size)
        {
            return nullptr;
        }

        uint8_t *old_head = head;
        head += pop_size;
        return old_head;
    }

    uint8_t *sk_buffer_t::push_to_tail(int32_t push_size)
    {
        if (tail_size() < push_size)
        {
            return nullptr;
        }
        uint8_t *old_tail = tail;
        tail += push_size;
        return old_tail;
    }

    uint8_t *sk_buffer_t::pop_from_tail(int32_t pop_size)
    {
        if (data_size() < pop_size)
        {
            return nullptr;
        }

        tail -= pop_size;
        return tail;
    }

}