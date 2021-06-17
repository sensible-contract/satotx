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

// SignUtxoSpendByUtxo
func SignUtxoSpendByUtxo(ctx *gin.Context) {
	log.Printf("SignUtxoSpendByUtxo enter")

	txIdHex := ctx.Param("txid")
	txIndexString := ctx.Param("index")

	byTxIdHex := ctx.Param("byTxid")
	byTxIndexString := ctx.Param("byTxindex")

	// check txid
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		log.Printf("txid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	// check txindex
	txIndex, err := strconv.Atoi(txIndexString)
	if err != nil || txIndex < 0 {
		log.Printf("txindex invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txindex invalid"})
		return
	}

	// check bytxid
	byTxIdReverse, err := hex.DecodeString(byTxIdHex)
	if err != nil {
		log.Printf("byTxid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "byTxid invalid"})
		return
	}
	byTxId := utils.ReverseBytes(byTxIdReverse)

	// check txindex
	byTxIndex, err := strconv.Atoi(byTxIndexString)
	if err != nil || byTxIndex < 0 {
		log.Printf("byTxIndex invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "byTxIndex invalid"})
		return
	}

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

	// by tx
	byTxRaw, err := hex.DecodeString(req.ByTxHex)
	if err != nil {
		log.Printf("byTxRaw invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "byTx invalid"})
		return
	}

	byTxObjHash := blkparser.GetHash256(byTxRaw)
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
	if len(byTxObj.TxOuts) <= byTxIndex {
		log.Printf("byTxIndex exceed limit")
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "byTxindex exceed limit"})
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

	// sig spend by
	txIndexRaw := make([]byte, 4)
	binary.LittleEndian.PutUint32(txIndexRaw, uint32(txIndex))
	payloadMsg := bytes.Join([][]byte{
		txId,
		txIndexRaw,
		txObj.TxOuts[txIndex].ValueRaw,
		blkparser.GetHash160(txObj.TxOuts[txIndex].Pkscript),
		byTxId,
	}, []byte{})
	sig, padding := rb.Sign(payloadMsg)
	if len(sig) == 0 {
		log.Printf("padding over 256 times.")
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "padding over 256 times"})
		return
	}

	// sig by utxo
	byTxIndexRaw := make([]byte, 4)
	binary.LittleEndian.PutUint32(byTxIndexRaw, uint32(byTxIndex))
	byPayloadMsg := bytes.Join([][]byte{
		byTxId,
		byTxIndexRaw,
		byTxObj.TxOuts[byTxIndex].ValueRaw,
		blkparser.GetHash160(byTxObj.TxOuts[byTxIndex].Pkscript),
	}, []byte{})

	bySig, byPadding := rb.Sign(byPayloadMsg)
	if len(bySig) == 0 {
		log.Printf("by padding over 256 times.")
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "by padding over 256 times"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: &model.TxTwoSigResponse{
			TxId:    txIdHex,
			Index:   txIndex,
			SigBE:   hex.EncodeToString(sig),
			SigLE:   hex.EncodeToString(utils.ReverseBytes(sig)),
			Padding: hex.EncodeToString(padding),
			Payload: hex.EncodeToString(payloadMsg),

			ByTxId:      byTxIdHex,
			ByTxIndex:   byTxIndex,
			ByTxSigBE:   hex.EncodeToString(bySig),
			ByTxSigLE:   hex.EncodeToString(utils.ReverseBytes(bySig)),
			ByTxPadding: hex.EncodeToString(byPadding),
			ByTxPayload: hex.EncodeToString(byPayloadMsg),
		}})
}
