package controller

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net/http"
	"satotx/lib/blkparser"
	"satotx/lib/utils"
	"satotx/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SignUtxo
func SignUtxo(ctx *gin.Context) {
	log.Printf("SignUtxo enter")

	txIdHex := ctx.Param("txid")
	txIndexString := ctx.Param("index")

	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		log.Printf("txid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	// check index
	txIndex, err := strconv.Atoi(txIndexString)
	if err != nil || txIndex < 0 {
		log.Printf("txindex invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txindex invalid"})
		return
	}

	// check body
	req := model.TxRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("Bind json failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "json error"})
		return
	}

	txRaw, err := hex.DecodeString(req.TxHex)
	if err != nil {
		log.Printf("txRaw invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "tx invalid"})
		return
	}

	txObjHash := blkparser.GetHash256(txRaw)
	if !bytes.Equal(txObjHash, txId) {
		log.Printf("txId(%s) not match hash(txBody)(%s)", hex.EncodeToString(txId), hex.EncodeToString(txObjHash))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txId not match txBody"})
		return
	}

	txObj, txoffset := blkparser.NewTx(txRaw)
	if int(txoffset) < len(txRaw) {
		log.Printf("txRaw invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "tx invalid"})
		return
	}
	if len(txObj.TxOuts) <= txIndex {
		log.Printf("txIndex exceed limit")
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txindex exceed limit"})
		return
	}

	txIndexRaw := make([]byte, 4)
	binary.LittleEndian.PutUint32(txIndexRaw, uint32(txIndex))
	payloadMsg := bytes.Join([][]byte{
		txId,
		txIndexRaw,
		txObj.TxOuts[txIndex].ValueRaw,
		blkparser.GetHash160(txObj.TxOuts[txIndex].Pkscript),
	}, []byte{})

	sig, padding := rb.Sign(payloadMsg)

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: &model.TxResponse{
			TxId:    txIdHex,
			Index:   txIndex,
			ByTxId:  "",
			Sig:     hex.EncodeToString(sig),
			Padding: hex.EncodeToString(padding),
			Payload: hex.EncodeToString(payloadMsg),
			Script:  hex.EncodeToString(txObj.TxOuts[txIndex].Pkscript),
		}})
}
