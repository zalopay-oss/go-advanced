// hello.c

#include "hello.h" // Đảm bảo việc hiện thực hàm thỏa mãn interface của module.
#include <stdio.h>

void SayHello(const char *s)
{
  puts(s);
}