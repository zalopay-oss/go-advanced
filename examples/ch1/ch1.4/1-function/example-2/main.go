package main

func main() {
    for i := 0; i < 3; i++ {
        i := i // Xác định một biến cục bộ i trong vòng lặp
        defer func(){ println(i) } ()
    }
}