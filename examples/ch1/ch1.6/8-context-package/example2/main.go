package main 

import (
    "fmt"
    "context"
)

// Trả về channel có chuỗi số: 2, 3, 4, ...
func GenerateNatural(ctx context.Context) chan int {
    ch := make(chan int)
    go func() {
        for i := 2; ; i++ {
            select {
            case <- ctx.Done():
                return
            case ch <- i:
            }
        }
    }()
    return ch
}

// Bộ lọc: xóa các số có thể chia hết cho số nguyên tố
func PrimeFilter(ctx context.Context, in <-chan int, prime int) chan int {
    out := make(chan int)
    go func() {
        for {
            if i := <-in; i%prime != 0 {
                select {
                case <- ctx.Done():
                    return
                case out <- i:
                }
            }
        }
    }()
    return out
}

func main() {
    // Kiểm soát trạng thái Goroutine nền thông qua context
    ctx, cancel := context.WithCancel(context.Background())

    ch := GenerateNatural(ctx) // chuỗi số: 2, 3, 4, ...
    for i := 0; i < 100; i++ {
        prime := <-ch // số nguyên tố mới
        fmt.Printf("%v: %v\n", i+1, prime)
        ch = PrimeFilter(ctx, ch, prime) // Bộ lọc dựa trên số nguyên tố mới
    }

    cancel()
}