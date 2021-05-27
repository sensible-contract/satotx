package controller

import (
	"log"
	"net/http"
	"satotx/model"

	"github.com/gin-gonic/gin"
)

// Satotx
func Satotx(ctx *gin.Context) {
	log.Printf("Satotx enter")

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "Welcome to use sensible contract on Bitcoin SV!",
		Data: &model.Welcome{
			PubKey:  rb.PubKeyHex,
			Contact: "",
			Job:     "",
			Github:  "https://github.com/sensible-contract",
		},
	})
}
