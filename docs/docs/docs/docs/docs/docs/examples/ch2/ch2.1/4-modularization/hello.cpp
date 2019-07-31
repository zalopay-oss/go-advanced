// hello.cpp

#include <iostream>

extern "C"
{
#include "hello.h"
}

void SayHello(const char *s)
{
  std::cout << s;
}