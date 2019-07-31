package main

/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
<<<<<<< HEAD
    const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    const char* os = "linux";
=======
    static const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    static const char* os = "linux";
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376
#else
#    error(unknown os)
#endif
*/
import "C"

func main() {
	print(C.GoString(C.os))
}
