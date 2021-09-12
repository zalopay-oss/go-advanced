package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("SELECT * FROM accounts WHERE id = ?;")

	res, err := stmt.Exec(2)
	res, err = stmt.Exec(3)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}
