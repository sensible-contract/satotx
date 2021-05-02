package main

import (
	"context"
	"log"
	"satotx/controller"

	"github.com/gin-gonic/gin"
	ginadapter "github.com/linthan/scf-go-api-proxy/gin"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Println("Gin cold start")
	router := gin.Default()

	router.GET("/", controller.Satotx)
	router.POST("/utxo/:txid/:index", controller.SignUtxo)
	router.POST("/utxo-spend-by/:txid/:index/:byTxid", controller.SignUtxoSpendBy)
	router.POST("/utxo-spend-by-utxo/:txid/:index/:byTxid/:byTxindex", controller.SignUtxoSpendByUtxo)

	ginLambda = ginadapter.New(router)
}

func Handler(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	cloudfunction.Start(Handler)
}
