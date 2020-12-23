package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"golang.org/x/sys/windows"
)

type error interface {
	Error() string
}

var ip = flag.String("ip", "localhost", "ip位置")
var acc = flag.String("acc", "fox", "帳號")
var pwd = flag.String("pwd", "fox", "密碼")
var epsv = flag.Bool("epsv", true, "EPSVmode")
var tn = flag.String("tn", "", "threadNum")

func main() {
	flag.Parse()
	fmt.Println("-ip:", *ip)
	fmt.Println("-acc:", *acc)
	fmt.Println("-pwd:", *pwd)
	fmt.Println("-epsv:", *epsv)
	fmt.Println("-tn:", *tn)
	threadNum, _ := strconv.Atoi(*tn)
	var wg sync.WaitGroup
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go uploadFile(i, *ip, *acc, *pwd, *epsv, &wg)

	}

	wg.Wait()
}
func uploadFile(i int, ip string, acc string, pwd string, epsv bool, wg *sync.WaitGroup) {
	//c, err := ftp.Dial(ip+":21", ftp.DialWithTimeout(5*time.Second), ftp.DialWithDebugOutput(os.Stdout), ftp.DialWithDisabledEPSV(epsv))
	c, err := ftp.Dial(ip+":21", ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(epsv))
	if err != nil {
		fmt.Println("執行序", windows.GetCurrentThreadId(), err.Error())
	}
	fmt.Printf("執行序%d:連線成功\n", windows.GetCurrentThreadId())
	err = c.Login(acc, pwd)
	if err != nil {
		fmt.Println("執行序", windows.GetCurrentThreadId(), err.Error())
	}
	fmt.Printf("執行序%d:Login成功\n", windows.GetCurrentThreadId())
	file, _ := ioutil.ReadFile("doc/Upload.pdf")
	data := bytes.NewBuffer(file)

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
