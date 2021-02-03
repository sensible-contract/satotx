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

func txValidator(ctx *gin.Context) (code int, err error) {
	return
}

// SignUtxoSpendBy
func SignUtxoSpendBy(ctx *gin.Context) {
	log.Printf("SignUtxoSpendBy enter")

	txIdHex := ctx.Param("txid")
	txIndexString := ctx.Param("index")

	byTxIdHex := ctx.Param("byTxid")

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

	byTxIdReverse, err := hex.DecodeString(byTxIdHex)
	if err != nil {
		log.Printf("byTxid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "byTxid invalid"})
		return
	}
	byTxId := utils.ReverseBytes(byTxIdReverse)

	// check body
	req := model.TxRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("Bind json failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "json error"})
		return
	}

	// tx
	txRaw, err := hex.DecodeString(req.TxHex)
	if err != nil {
		log.Printf("txRaw invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "tx invalid"})
		return
	}

	txObjHash := blkparser.GetShaString(txRaw)
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

	// by tx
	byTxRaw, err := hex.DecodeString(req.ByTxHex)
	if err != nil {
		log.Printf("byTxRaw invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "byTx invalid"})
		return
	}

	byTxObjHash := blkparser.GetShaString(byTxRaw)
	if !bytes.Equal(byTxObjHash, byTxId) {
		log.Printf("byTxId(%s) not match hash(txBody)(%s)", hex.EncodeToString(byTxId), hex.EncodeToString(byTxObjHash))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "byTxId not match byTxBody"})
		return
	}

	byTxObj, byTxoffset := blkparser.NewTx(byTxRaw)
	if int(byTxoffset) < len(byTxRaw) {
		log.Printf("byTxRaw invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "byTx invalid"})
		return
	}

	foundInput := false
	for _, txIn := range byTxObj.TxIns {
		if !bytes.Equal(txIn.InputHash, txId) {
			continue
		}
		if txIndex != int(txIn.InputVout) {
			continue
		}
		foundInput = true
		break
	}
	if !foundInput {
		log.Printf("not spend by tx: %s", req.ByTxHex)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "not spend by tx"})
		return
	}

	txIndexRaw := make([]byte, 4)
	binary.LittleEndian.PutUint32(txIndexRaw, uint32(txIndex))
	payloadMsg := bytes.Join([][]byte{
		txId,
		txIndexRaw,
		txObj.TxOuts[txIndex].ValueRaw,
		txObj.TxOuts[txIndex].Pkscript,
		byTxId,
	}, []byte{})

	sig, padding := rb.Sign(payloadMsg)
	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: &model.TxResponse{
			TxId:    txIdHex,
			Index:   txIndex,
			ByTxId:  byTxIdHex,
			Sig:     hex.EncodeToString(sig),
			Padding: hex.EncodeToString(padding),
			Payload: hex.EncodeToString(payloadMsg),
		}})

}
