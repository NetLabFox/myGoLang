package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

type TracMoData struct {
	ID        string `json:"id"`
	Timestamp int    `json:"timestamp"` //unixtimestamp
	KeyA      string `json:"Key_A"`
	KeyB      string `json:"Key_B"`
	KeyC      string `json:"Key_C"`
	KeyD      string `json:"Key_D"`
	Userid    string `json:"user_id"`
}

func accessHandler(c *gin.Context) {

	var resultString []TracMoData
	timestamp := c.Query("timestamp")
	if timestamp != "" {
		iTimestamp, err := strconv.ParseUint(timestamp, 10, 64) //Get過來是字串需轉換
		if err != nil {
			return
		}
		cred := credentials.NewStaticCredentials("AKIAINQVSGZ5H2FRWZEA", "KBAFktoJJcUAH6qWv+AvqJrCnRuTXu1dbV9/7zKH", ``)
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String("ap-southeast-1"),
			Credentials: cred,
		})
		svc := dynamodb.New(sess)                                                           //DB連線
		filt := expression.Name("timestamp").GreaterThanEqual(expression.Value(iTimestamp)) //區間查詢
		expr, err := expression.NewBuilder().WithFilter(filt).Build()
		if err != nil {
			fmt.Println(" Error expression:")
			fmt.Println(err.Error())
		}
		//執行查詢
		params := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			TableName:                 aws.String("T_Tracmo"),
		}
		result, err := svc.Scan(params)

		if err != nil {
			fmt.Println("Query Failed:")
			fmt.Println((err.Error()))
		}
		for _, i := range result.Items {
			tracmodata := TracMoData{} //建立TracMoData的buffer

			err = dynamodbattribute.UnmarshalMap(i, &tracmodata) //塞進去之前的buffer

			if err != nil {
				fmt.Println("UnmarshalMap Error:")
				fmt.Println(err.Error())

			}
			resultString = append(resultString, tracmodata)

		}
		b, err := json.Marshal(resultString) //json格式化
		if err != nil {
			fmt.Println("Marshal Error:", err)
		}
		c.String(http.StatusOK, string(b)) //回傳
	}
}
func routerEngine() *gin.Engine {
	// set server mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/gAPI", accessHandler)

	return r
}
func main() {

	//gin.SetMode(gin.ReleaseMode)
	addr := ":" + os.Getenv("PORT")
	//addr = "localhost:4000"
	log.Fatal(gateway.ListenAndServe(addr, routerEngine()))

}
