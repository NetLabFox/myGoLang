package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func myWeb(w http.ResponseWriter, r *http.Request) {

	t := template.New("index")

	t.Parse("<div id='templateTextDiv'>Hi,{{.name}},{{.someStr}}</div>")

	data := map[string]string{
		"name":    "zeta",
		"someStr": "這是一個開始",
	}

	t.Execute(w, data)

	// fmt.Fprintln(w, "這是一個開始")
}

func main() {
	http.HandleFunc("/", myWeb)

	fmt.Println("服務器即將開啓，訪問地址 http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服務器開啓錯誤: ", err)
	}
}
