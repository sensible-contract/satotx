package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"satotx/controller"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// 0.0.0.0:8000
	listen_address = os.Getenv("LISTEN")
)

func main() {
	router := gin.Default()

	router.GET("/", controller.Satotx)
	router.POST("/utxo/:txid/:index", controller.SignUtxo)
	router.POST("/utxo-spend-by/:txid/:index/:byTxid", controller.SignUtxoSpendBy)
	router.POST("/utxo-spend-by-utxo/:txid/:index/:byTxid/:byTxindex", controller.SignUtxoSpendByUtxo)

	log.Printf("LISTEN: %s", listen_address)
	svr := &http.Server{
		Addr:    listen_address,
		Handler: router,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
