package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//TableName gorm透過結構會自動加s
func (Ssldomainlist Ssldomainlist) TableName() string {
	return "Ssldomainlist"
}

// Ssldomainlist [...]
type Ssldomainlist struct {
	Domain string `gorm:"primary_key;column:domain;type:varchar(80);not null" json:"-"`
	TaxID  string `gorm:"column:taxID;type:varchar(8);not null" json:"tax_id"`
	Status int    `gorm:"column:status;type:int;not null" json:"status"`
}

func main() {

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
	// batch size 100
	//var results []Ssldomainlist
	results := []Ssldomainlist{}
	//var results []map[string]interface{}
	result := db.Where("status = ?", 1).FindInBatches(&results, 1000, func(tx *gorm.DB, batch int) error {
		/*一次修正完畢
		if err := tx.Update("status", 1).Error; err != nil {
			// return any error will rollback
			return err
		}*/

		for i, result := range results {
			//result.Status = 0
			//result.Status = 0
			results[i].Status = 0
			fmt.Println(result)
		}

		/*	for _, result := range results {
			fmt.Println(result)
		}*/
		//fmt.Println("1", tx.RowsAffected)
		tx.Save(&results)
		fmt.Println("2", tx.RowsAffected) // number of records in this batch
		fmt.Println(batch)                // Batch 1, 2, 3
		/*	if tx.Error != nil {

			fmt.Println(tx.Error)
		}*/
		// returns error will stop future batches
		return nil
	})
	fmt.Println(result.Error)        // returned error
	fmt.Println(result.RowsAffected) // processed records count in all batches

}
