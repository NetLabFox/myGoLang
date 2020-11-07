package main

import (
	"database/sql"

	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
)

type error interface {
	Error() string
}

func main() {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "root:!QAZ2wsx@tcp(127.0.0.1:3306)/sslverify")

	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	f, err := excelize.OpenFile("C:\\Users\\CHT-User\\Desktop\\0714.xlsx")

	if err != nil {

		println("Get Files Error：" + err.Error())

		return

	}
	rows, err := f.GetRows("Sheet1")

	if err != nil {

		println("Get Rows Error：" + err.Error())

		return

	}
	stmt, err := db.Prepare("INSERT INTO ssldomainlist (domain,taxID) VALUES (?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, colCell := range rows {
		_, err := stmt.Exec(colCell[0], colCell[1])
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
