#include <cstdint>

namespace net_stack
{

    struct sk_buffer_t
    {
        uint32_t size;

        /**
         * bgn .... head .... tail ... end
         */
        uint8_t *bgn, *head, *tail, *end;

        void alloc();
        void free();

        int32_t head_size();
        int32_t tail_size();
        int32_t data_size();

        uint8_t *push_to_head(int32_t size);
        uint8_t *pop_from_head(int32_t size);

        uint8_t *push_to_tail(int32_t size);
        uint8_t *pop_from_tail(int32_t size);

        sk_buffer_t();
        sk_buffer_t(uint32_t size);
        sk_buffer_t(const sk_buffer_t &sk_b);
        sk_buffer_t(sk_buffer_t &&sk_b);
        sk_buffer_t &operator=(const sk_buffer_t &sk_b);

        ~sk_buffer_t();
    };
}