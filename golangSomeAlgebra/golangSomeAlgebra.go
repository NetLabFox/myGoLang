package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var count int

//Expresion 結果
type Expresion interface {
	result() int
}

//StringMultiplication 乘法
type StringMultiplication struct {
	expresion string
}

//StringAddition 加法
type StringAddition struct {
	expresion string
}

func (SM StringMultiplication) result() int {
	temp := strings.Split(SM.expresion, "*")
	num1, _ := strconv.Atoi(temp[0])
	num2, _ := strconv.Atoi(temp[1])
	return num1 * num2
}

func (SA StringAddition) result() int {
	temp := strings.Split(SA.expresion, "+")
	num1, _ := strconv.Atoi(temp[0])
	num2, _ := strconv.Atoi(temp[1])
	return num1 + num2
}

//待補

func getResult(e Expresion) int {
	return e.result()
}

func generateVal(channel chan int, query string) {

	//	time.Sleep(500 * time.Millisecond)
	val := 0
	if query[1] == '+' {
		val = getResult(StringAddition{query})
	} else {
		val = getResult(StringMultiplication{query})
	}
	//待補
	//	channel <- val case就已經執行一次了
	select {
	case channel <- val:
		{
			count++
			if count == 5 {
				close(channel) //不close range不能用
			}
		}
	}
}

func golangSomeAlgebra(queries []string) []int {

	result := []int{}
	channel := make(chan int)
	for i := 0; i < 5; i++ {
		if len(queries) == (i + 1) {

			go generateVal(channel, queries[i])

		} else {

			go generateVal(channel, queries[i])

		}
		queries[i] = ""
	}
	//待補

	for value := range channel { //range 要等close 才可以使用，所以可以當作wait來使用

		result = append(result, value)
	}

	sort.Ints(result)
	return result
}
func main() {
	//var wg sync.WaitGroup
	queries := []string{"1+3", "0+2", "9+8", "9*2", "1*3"}
	results := golangSomeAlgebra(queries)

	for _, result := range results {
		fmt.Println(result)
	}

}
