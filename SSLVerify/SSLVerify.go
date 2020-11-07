package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type error interface {
	Error() string
}

func main() {
	dic := [11]string{"", "www", "autodiscover", "app", "api", "shop", "eip", "uat", "portal", "vpn", "service"}
	f, err := excelize.OpenFile("C:\\Users\\Fox-PC\\Desktop\\0714.xlsx")

	if err != nil {

		println("Get Files Error：" + err.Error())

		return

	}
	rows, err := f.GetRows("Sheet1")

	if err != nil {

		println("Get Rows Error：" + err.Error())

		return

	}
	count := 0
	for _, colCell := range rows {

		for _, s := range dic {
			issuer, notAfter, certType, err := getSSLInfo(colCell[0], s)
			if err == nil {
				//fmt.Println(s, issuer, notAfter, certType)
				f.SetCellValue("Sheet2", fmt.Sprintf("A%v", count), s+"."+colCell[0])
				f.SetCellValue("Sheet2", fmt.Sprintf("B%v", count), colCell[1])
				f.SetCellValue("Sheet2", fmt.Sprintf("C%v", count), issuer)
				f.SetCellValue("Sheet2", fmt.Sprintf("D%v", count), notAfter)
				f.SetCellValue("Sheet2", fmt.Sprintf("E%v", count), certType)
				count++
				f.Save()
			} else {
				var test = err.Error()
				f.SetCellValue("Sheet2", fmt.Sprintf("A%v", count), s+"."+colCell[0])
				f.SetCellValue("Sheet2", fmt.Sprintf("B%v", count), colCell[1])
				f.SetCellValue("Sheet2", fmt.Sprintf("F%v", count), test)
				count++
				f.Save()
				//fmt.Println(test)
			}
		}

	}

}
func getSSLInfo(domain string, prefix string) (string, string, string, error) {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", prefix+"."+domain+":443", conf)
	if err != nil {

		return "", "", "", err
	}
	/*	conn, err := net.DialTimeout("tcp", "nnnn.com.tw:443", 5*time.Second)
		conn = tls.Client(conn, conf)*/
	if conn != nil {
		defer conn.Close()
		certs := conn.ConnectionState().PeerCertificates
		notAfter := certs[0].NotAfter.Format("2006-January-02")
		issuer := certs[0].Issuer.CommonName
		var certType string
		var normalBool = false
		//fmt.Printf("Issuer Name: %s\n", issuer)
		//fmt.Printf("Expiry: %s \n", notAfter)

		PID := certs[0].PolicyIdentifiers
		for _, newP := range PID {
			if strings.HasPrefix(newP.String(), "2.23.140") {
				//	fmt.Printf("PolicyIdentifiers: %s \n", newP)
				switch newP.String() {
				case "2.23.140.1.2.1":
					certType = "DV"
				case "2.23.140.1.2.2":
					certType = "OV"
				case "2.23.140.1.2.3":
					certType = "IV"
				case "2.23.140.1.1":
					certType = "EV"
				}
				normalBool = true
			}
		}
		if !normalBool {
			for _, newP := range PID {
				certType += newP.String()
			}
		}
		return issuer, notAfter, certType, nil
	}

	return "", "", "", nil
}
