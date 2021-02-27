package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	oriTotal := 0
	oriHash := ""
	path := "C:\\Users\\Fox-PC\\Desktop\\raweb_SIT.txt"
	for true {

		hash, err := chechHash(path)
		if err != nil {
			log.Fatal(err)
		}
		s := fmt.Sprintf("%x", hash)
		if s != oriHash {
			total, err := readLine(path, oriTotal)
			if err != nil {
				oriTotal = total
				fmt.Println(total)
			} else {
				fmt.Println("解析錯誤")
			}
			oriHash = s
		}
		time.Sleep(1 * time.Minute)
	}
	/*	scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}*/
}
func chechHash(path string) (hash []byte, err error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), err
}
func readLine(path string, lineNum int) (lastLine int, err error) {
	r, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {

		if lastLine > lineNum {
			// you can return sc.Bytes() if you need output in []bytes

			fmt.Println(sc.Text())
		} else {
			lastLine++
		}
	}
	return lastLine, io.EOF
}
