package controllers

import (
	"RestApiGin/database"
	"RestApiGin/models"
	"strconv"

	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status string             `json:"status"`
	Data   map[string]float64 `json:"data"`
}

type WalletRepo struct {
	Db *gorm.DB
}

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func New() *WalletRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Wallet{})
	return &WalletRepo{Db: db}
}
func errorFunc(status string, message string) (err ErrorResponse) {
	err = ErrorResponse{
		Status:  status,
		Message: message,
	}
	return err
}

func (repository *WalletRepo) DebitTransaction(c *gin.Context) {
	idParam := c.Param("id")
	var wallet models.Wallet
	c.BindJSON(&wallet)

	if reflect.TypeOf(wallet.Debit).Kind() != reflect.Float64 {
		log.Debugf("Error : Values not Float64")
		response := errorFunc("error", "Check your type of values")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if wallet.Debit <= 0 {
		log.Debugf("Error : Debit values should be more then 0", idParam)
		response := errorFunc("error", "Error : Debit values should be more then 0")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	dt := time.Now()
	strID := wallet.CustId + dt.String()
	var id = base64.StdEncoding.EncodeToString([]byte(strID))
	wallet.ID = id
	wallet.CustId = idParam

	err := models.CreateTrascation(repository.Db, &wallet)
	if err != nil {
		log.Debugf("Error : ", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := client.Get(idParam).Result()
	if err != nil {
		fmt.Println(err)
	}
	balance, _ := strconv.ParseFloat(res, 64)

	resp := SuccessResponse{
		Status: "Debit Success",
		Data: map[string]float64{
			"debit":   wallet.Debit,
			"balance": balance,
		},
	}
	c.JSON(http.StatusOK, resp)
}

func (repository *WalletRepo) CreditTransaction(c *gin.Context) {
	idParam := c.Param("id")
	var wallet models.Wallet
	c.BindJSON(&wallet)
	res, err := client.Get(idParam).Result()
	if err != nil {
		fmt.Println(err)
	}
	balance, _ := strconv.ParseFloat(res, 64)
	if reflect.TypeOf(wallet.Credit).Kind() != reflect.Float64 {
		log.Debugf("Error : Values not Float64")
		response := errorFunc("error", "Check your type of values")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if balance < wallet.Credit {
		log.Debugf("Error : Balance not enough")
		response := errorFunc("error", "Balance not enough")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if wallet.Credit <= 0 {
		log.Debugf("Error : Credit values should be more then 0", idParam)
		response := errorFunc("error", "Error : Credit values should be more then 0")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	dt := time.Now()
	strID := wallet.CustId + dt.String()
	var id = base64.StdEncoding.EncodeToString([]byte(strID))
	wallet.ID = id
	wallet.CustId = idParam
	c.BindJSON(&wallet)
	err = models.CreateTrascation(repository.Db, &wallet)
	if err != nil {
		log.Debugf("Error : ", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	resp := SuccessResponse{
		Status: "Credit Success",
		Data: map[string]float64{
			"credit":  wallet.Credit,
			"balance": balance,
		},
	}
	c.JSON(http.StatusOK, resp)

}

func (repository *WalletRepo) CheckTransaction(c *gin.Context) {
	idParam := c.Param("id")
	var wallet models.Wallet
	models.GetWallet(repository.Db, &wallet, idParam)
	res, err := client.Get(idParam).Result()
	if err != nil {
		fmt.Println(err)
	}
	balance, _ := strconv.ParseFloat(res, 64)
	if wallet.ID != "" {
		resp := SuccessResponse{
			Status: "Balance Success",
			Data: map[string]float64{
				"balance": balance,
			},
		}
		c.JSON(http.StatusOK, resp)
	} else {
		log.Debugf("Error : Invalid Customer ID : ", idParam)
		response := errorFunc("error", "Invalid Customer ID : "+idParam)
		c.JSON(http.StatusBadRequest, response)
		return
	}

}
