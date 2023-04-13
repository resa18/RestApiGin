package main

import (
	"RestApiGin/controllers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "debug"
	}
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)
}
func main() {

	r := setupRouter()
	_ = r.Run(":8000")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	userRepo := controllers.New()
	r.POST("/wallet/:id/debit", userRepo.DebitTransaction)
	r.POST("/wallet/:id/credit", userRepo.CreditTransaction)
	r.GET("/wallet/:id/balance", userRepo.CheckTransaction)

	return r
}
