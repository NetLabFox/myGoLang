package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	f, err := excelize.OpenFile("C:\\Users\\Fox-PC\\Desktop\\活頁簿2.xlsx")

	if err != nil {

		println(err.Error())

		return

	}
	rows, err := f.GetRows("Sheet1")
	m := make(map[string]int)
	for _, row := range rows {

		for _, colCell := range row {

			m[strings.Split(colCell, ".")[0]]++

		}

	}
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	}
}
