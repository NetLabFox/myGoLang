package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"golang.org/x/sys/windows"
)

type error interface {
	Error() string
}

func main() {

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go uploadFile(i, &wg)

	}

	wg.Wait()
}
func uploadFile(i int, wg *sync.WaitGroup) {
	c, err := ftp.Dial("localhost:21", ftp.DialWithTimeout(5*time.Second), ftp.DialWithDebugOutput(os.Stdout))

	if err != nil {
		fmt.Println("執行序", windows.GetCurrentThreadId(), err.Error())
	}
	fmt.Printf("執行序%d:連線成功\n", windows.GetCurrentThreadId())
	err = c.Login("fox", "fox")
	if err != nil {
		fmt.Println("執行序", windows.GetCurrentThreadId(), err.Error())
	}
	fmt.Printf("執行序%d:Login成功\n", windows.GetCurrentThreadId())
	data := bytes.NewBufferString("Hello World")
	err = c.Stor("test-file.txt"+strconv.Itoa(i), data)
	if err != nil {
		fmt.Println("執行序", windows.GetCurrentThreadId(), err.Error())
	}
	fmt.Printf("執行序%d:上傳成功\n", windows.GetCurrentThreadId())
	if err := c.Quit(); err != nil {
		fmt.Println("執行序", windows.GetCurrentThreadId(), err.Error())
	}
	fmt.Printf("執行序%d:連線結束\n", windows.GetCurrentThreadId())
	wg.Done()
}
