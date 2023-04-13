package models

import (
	"fmt"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type Wallet struct {
	gorm.Model
	ID     string  `json:"id" gorm:"primary_key"`
	CustId string  `json:"customerId"`
	Debit  float64 `json:"debit"`
	Credit float64 `json:"credit"`
}

func CreateTrascation(db *gorm.DB, Wallet *Wallet) (err error) {
	log.Info("Creating transaction")
	err = db.Create(Wallet).Error
	if err != nil {
		return err
	}
	valBalance := CheckBalance(db, Wallet.CustId)
	err = client.Set(Wallet.CustId, valBalance, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func GetWallet(db *gorm.DB, Wallet *Wallet, id string) {
	log.Info("Get transaction by id : ", id)
	db.Where("cust_id = ?", id).First(Wallet)
	return
}

func CheckBalance(db *gorm.DB, id string) (balance float64) {
	log.Info("Check Balance")
	db.Table("wallets").Select("sum(debit)-sum(credit)").Where("cust_id = ?", id).Row().Scan(&balance)

	return
}
