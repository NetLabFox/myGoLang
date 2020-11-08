package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type error interface {
	Error() string
}

//TableName gorm透過結構會自動加s
func (Ssldomainlist Ssldomainlist) TableName() string {
	return "Ssldomainlist"
}

//TableName gorm透過結構會自動加s
func (Sslresult Sslresult) TableName() string {
	return "Sslresult"
}

// Ssldomainlist [...]
type Ssldomainlist struct {
	Domain string `gorm:"primary_key;column:domain;type:varchar(80);not null" json:"Domain"`
	TaxID  string `gorm:"column:taxID;type:varchar(8);not null" json:"tax_id"`
	Status int    `gorm:"column:status;type:int;not null" json:"status"`
}

// Sslresult [...]
type Sslresult struct {
	PrefixWithDomain string    `gorm:"primary_key;column:prefixWithDomain;type:varchar(70);not null" json:"PrefixWithDomain"`
	Issuer           string    `gorm:"column:issuer;type:varchar(100)" json:"issuer"`
	Notafter         time.Time `gorm:"column:notafter;type:date" json:"notafter"`
	CertType         string    `gorm:"column:certType;type:varchar(50)" json:"cert_type"`
	ErrorMsg         string    `gorm:"column:errorMsg;type:varchar(300)" json:"error_msg"`
	UpdateTime       time.Time `gorm:"column:updateTime;type:datetime" json:"update_time"`
	TaxID            string    `gorm:"column:taxID;type:varchar(8);not null" json:"tax_id"`
}

func main() {

	var wg sync.WaitGroup
	var rowscount int64
	//var results []Ssldomainlist
	dsn := "test:sslverify@tcp(127.0.0.1:3306)/sslverify?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("連線失敗")
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("取得DB失敗")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(150)
	sqlDB.SetConnMaxLifetime(time.Hour)
	count := 0
Loop:
	rows, err := db.Limit(10000).Model(&Ssldomainlist{}).Where("status = ?", 0).Rows()
	defer rows.Close()

	for rows.Next() {
		var Ssldomainlist Ssldomainlist
		db.ScanRows(rows, &Ssldomainlist)
		go getSSLInfo(Ssldomainlist.Domain, Ssldomainlist.TaxID, db, &wg)
		wg.Add(1)
		time.Sleep(300 * time.Millisecond)
		count++

	}
	db.Limit(10000).Model(&Ssldomainlist{}).Where("status = ?", 0).Update("status", 1)
	wg.Wait()
	db.Model(&Ssldomainlist{}).Where("status = ?", 0).Count(&rowscount)
	if rowscount > 0 {
		goto Loop
	}

	fmt.Println("all the tasks done...", count)
}
func getSSLInfo(domain string, TaxID string, db *gorm.DB, wg *sync.WaitGroup) {

	dic := [15]string{"www", "autodiscover", "app", "api", "shop", "eip", "uat", "portal", "vpn", "service", "mail", "webmail", "pop", "smtp", "imap"}

	var Sslresults []Sslresult
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	port := ":443"
	for _, prefix := range dic {
		var Sslresult Sslresult
		switch prefix {
		case "mail":
			port = ":587"
		case "webmail":
			port = ":587"
		case "pop":
			port = ":995"
		case "smtp":
			port = ":587"
		case "imap":
			port = ":993"
		default:
			port = ":443"
		}

		var wholeDomain = prefix + "." + domain + port

		conn, err := tls.DialWithDialer(&net.Dialer{
			Timeout: 5 * time.Second,
		}, "tcp", wholeDomain, conf)

		if err != nil {
			if prefix == "pop" {
				Sslresult = smtpSSL(domain, prefix, ":110", TaxID, db)
				Sslresults = append(Sslresults, Sslresult)
				continue
				//return issuer, notAfter, certType, err
			} else if prefix == "imap" {
				Sslresult = smtpSSL(domain, prefix, ":143", TaxID, db)
				Sslresults = append(Sslresults, Sslresult)
				continue
				//return issuer, notAfter, certType, err
			}
			Sslresult.CertType = ""
			Sslresult.Issuer = ""
			Sslresult.Notafter = time.Time{}

			Sslresult.PrefixWithDomain = prefix + "." + domain
			Sslresult.ErrorMsg = err.Error()
			Sslresult.UpdateTime = time.Now()
			Sslresult.TaxID = TaxID
			/*	result := db.Create(&Sslresult)
				if result.Error != nil {
					fmt.Println(result.Error)
				}*/
			Sslresults = append(Sslresults, Sslresult)
			continue
			//return "", "", "", err
		}
		/*  conn, err := net.DialTimeout("tcp", "nnnn.com.tw:443", 5*time.Second)
		    conn = tls.Client(conn, conf)*/
		if conn != nil {
			defer conn.Close()
			certs := conn.ConnectionState().PeerCertificates
			//notAfter := certs[0].NotAfter.Format("yyyy/MM/dd")
			issuer := certs[0].Issuer.CommonName
			if issuer == "" {
				if certs[0].Issuer.OrganizationalUnit != nil {
					issuer = certs[0].Issuer.OrganizationalUnit[0]
				}

			}
			var certType string
			var normalBool = false
			//fmt.Printf("Issuer Name: %s\n", issuer)
			//fmt.Printf("Expiry: %s \n", notAfter)

			PID := certs[0].PolicyIdentifiers
			for _, newP := range PID {
				if strings.HasPrefix(newP.String(), "2.23.140") {
					//  fmt.Printf("PolicyIdentifiers: %s \n", newP)
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
			Sslresult.Issuer = issuer
			Sslresult.Notafter = certs[0].NotAfter
			Sslresult.CertType = certType
			Sslresult.PrefixWithDomain = prefix + "." + domain
			Sslresult.UpdateTime = time.Now()
			Sslresult.TaxID = TaxID
			Sslresult.ErrorMsg = ""
			Sslresults = append(Sslresults, Sslresult)
			/*result := db.Create(&Sslresult)
			if result.Error != nil {
				fmt.Println(result.Error)
			}*/
			continue
		}
		//return issuer, notAfter, certType, nil
	}
	result := db.Create(&Sslresults)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	wg.Done()
}

func smtpSSL(domain string, prefix string, port string, TaxID string, db *gorm.DB) (Sslresult Sslresult) {
	//var Sslresult Sslresult
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	c, err := smtp.Dial(prefix + "." + domain + port)

	if err != nil {
		Sslresult.CertType = ""
		Sslresult.Issuer = ""
		Sslresult.Notafter = time.Time{}
		Sslresult.PrefixWithDomain = prefix + "." + domain
		Sslresult.ErrorMsg = err.Error()
		Sslresult.UpdateTime = time.Now()
		Sslresult.TaxID = TaxID
		/*result := db.Create(&Sslresult)
		if result.Error != nil {
			fmt.Println(result.Error)
		}*/

		return
		//	return "", "", "", err
	}
	c.StartTLS(conf)
	state, ok := c.TLSConnectionState()
	if ok {
		certs := state.PeerCertificates
		//notAfter := certs[0].NotAfter.Format("yyyy/MM/dd")
		issuer := certs[0].Issuer.CommonName
		if issuer == "" {
			if certs[0].Issuer.OrganizationalUnit != nil {
				issuer = certs[0].Issuer.OrganizationalUnit[0]
			}

		}
		var certType string
		var normalBool = false
		//fmt.Printf("Issuer Name: %s\n", issuer)
		//fmt.Printf("Expiry: %s \n", notAfter)

		PID := certs[0].PolicyIdentifiers
		for _, newP := range PID {
			if strings.HasPrefix(newP.String(), "2.23.140") {
				//  fmt.Printf("PolicyIdentifiers: %s \n", newP)
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
		Sslresult.Issuer = issuer
		//Sslresult.Notafter, err = time.Parse("yyyy/MM/dd", notAfter)
		Sslresult.Notafter = certs[0].NotAfter

		Sslresult.CertType = certType
		Sslresult.PrefixWithDomain = prefix + "." + domain
		Sslresult.UpdateTime = time.Now()
		Sslresult.TaxID = TaxID
		Sslresult.ErrorMsg = ""
		/*result := db.Create(&Sslresult)
		if result.Error != nil {
			fmt.Println(result.Error)
		}*/
		return
		//return issuer, notAfter, certType, nil
	}
	return
}
