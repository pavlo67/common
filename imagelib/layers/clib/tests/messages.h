#include <inttypes.h>

typedef struct {
    char*   text;
    int16_t len;
} Message;

void  msg_write(Message* msg, const char* str);
