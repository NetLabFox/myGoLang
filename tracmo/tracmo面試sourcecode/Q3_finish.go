package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

type TracMoData struct {
	ID        string `json:"id"`
	Timestamp int    `json:"timestamp"`
	KeyA      string `json:"Key_A"`
	KeyB      string `json:"Key_B"`
	KeyC      string `json:"Key_C"`
	KeyD      string `json:"Key_D"`
	Userid    string `json:"user_id"`
}
type TracMoDataArray struct {
	TracMo []TracMoData `json:"tracmo"`
}

func s3Handler(c *gin.Context) {
	var filePath = "./tracmoAWS.csv"
	RDbytes, err := c.GetRawData() //抓post過來的資料
	if err != nil {

		return
	}
	var TracMoDataA TracMoDataArray
	json.Unmarshal(RDbytes, &TracMoDataA)
	var header = make([]string, 1)
	header[0] = "#id,#timestamp,#user_id,#Key_A,#Key_B,#Key_C,#Key_D" //產生csv'sheader

	//將結構內資料轉成csv格式
	for i := 0; i < len(TracMoDataA.TracMo); i++ {
		if TracMoDataA.TracMo[i].ID == "" {
			continue
		}
		if TracMoDataA.TracMo[i].Timestamp == 0 {
			continue
		}
		if TracMoDataA.TracMo[i].KeyA == "" {
			TracMoDataA.TracMo[i].KeyA = " "
		}
		if TracMoDataA.TracMo[i].KeyB == "" {
			TracMoDataA.TracMo[i].KeyC = " "
		}
		if TracMoDataA.TracMo[i].KeyC == "" {
			TracMoDataA.TracMo[i].KeyB = " "
		}
		if TracMoDataA.TracMo[i].KeyD == "" {
			TracMoDataA.TracMo[i].KeyD = " "
		}
		if TracMoDataA.TracMo[i].Userid == "" {
			TracMoDataA.TracMo[i].Userid = " "
		}
		var record []string
		record = append(record, TracMoDataA.TracMo[i].ID)
		record = append(record, strconv.Itoa(TracMoDataA.TracMo[i].Timestamp))
		record = append(record, TracMoDataA.TracMo[i].Userid)
		record = append(record, TracMoDataA.TracMo[i].KeyA)
		record = append(record, TracMoDataA.TracMo[i].KeyB)
		record = append(record, TracMoDataA.TracMo[i].KeyC)
		record = append(record, TracMoDataA.TracMo[i].KeyD)
		header = append(header, strings.Join(record, ", "))
	}
	var result string = strings.Join(header, "\r\n")

	cred := credentials.NewStaticCredentials("AKIAINQVSGZ5H2FRWZEA", "KBAFktoJJcUAH6qWv+AvqJrCnRuTXu1dbV9/7zKH", ``) //權限
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: cred,
	})
	svc := s3.New(sess)        //s3連線開啟
	buffer := []byte(result)   //字串變成bytes
	size := int64(len(buffer)) //s3參數只吃int64

	params := &s3.PutObjectInput{
		Bucket:        aws.String("tracmobucket"),
		Key:           aws.String(filePath),
		ACL:           aws.String("public-read"), //要公開
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("application/octet-stream"), //csv格式
	}

	resp, err := svc.PutObject(params)
	_ = resp
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	c.String(http.StatusOK, "https://s3-ap-southeast-1.amazonaws.com/tracmobucket/tracmoAWS.csv")
}
func routerEngine() *gin.Engine {
	// set server mode
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/pAPI2", s3Handler)

	return r
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	addr := ":" + os.Getenv("PORT")
	//addr = "localhost:4000"
	log.Fatal(gateway.ListenAndServe(addr, routerEngine()))

}
