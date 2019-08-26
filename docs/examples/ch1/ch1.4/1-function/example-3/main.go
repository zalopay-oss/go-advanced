package main

func main() {
    for i := 0; i < 3; i++ {
        // truyền i vào hàm (pass by value)
        // câu lệnh defer sẽ lấy các tham số từ lời gọi
        defer func(i int){ println(i) } (i)
    }
}