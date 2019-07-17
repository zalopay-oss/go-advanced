package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// db là một đối tượng của kiểu sql.DB
	// đối tượng là một thread-safe chứa kết nối
	// Tùy chọn kết nối có thể được đặt trong phương thức sql.DB, ở đây bỏ qua
	db, err := sql.Open("mysql",
		"username:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		id   int
		name string
	)
	rows, err := db.Query("select id, name from users where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	// Đọc nội dung các rows rồi gọi Close()
	// kết nối sẽ không được giải phóng cho đến khi defer rows.Close() thực thi
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
