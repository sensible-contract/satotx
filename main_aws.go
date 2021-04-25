package main

import (
	"context"
	"log"
	"satotx/controller"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
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

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
