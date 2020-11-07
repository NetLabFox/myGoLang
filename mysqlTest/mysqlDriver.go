package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type error interface {
	Error() string
}

//Data struct is for table big5toutf8 with scan
type Data struct {
	CNameBig5 string `json:"cNameBig5"`
	CNameCNS  string `json:"cNameCNS"`
	CNameUTF8 string `json:"cNameUTF8"`
}

func main() {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "root:!QAZ2wsx@tcp(127.0.0.1:3306)/big5toutf8")

	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	results, err := db.Query("SELECT cNameBig5,cNameCNS,cNameUTF8 FROM big5toutf8")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var data Data
		// for each row, scan the result into our tag composite object
		err = results.Scan(&data.CNameBig5, &data.CNameCNS, &data.CNameUTF8)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		e, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(e))
	}

}
