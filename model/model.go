package model

import "encoding/json"

type TxRequest struct {
	TxHex   string `json:"txHex"`
	ByTxHex string `json:"byTxHex"`
}

type TxResponse struct {
	TxId    string `json:"txId"`
	Index   int    `json:"index"`
	SigBE   string `json:"sigBE"`
	SigLE   string `json:"sigLE"`
	Padding string `json:"padding"`
	Payload string `json:"payload"`
	Script  string `json:"script"`

	ByTxId string `json:"byTxId"`
}

type TxTwoSigResponse struct {
	TxId    string `json:"txId"`
	Index   int    `json:"index"`
	SigBE   string `json:"sigBE"`
	SigLE   string `json:"sigLE"`
	Padding string `json:"padding"`
	Payload string `json:"payload"`
	Script  string `json:"script"`

	ByTxId      string `json:"byTxId"`
	ByTxIndex   int    `json:"byTxIndex"`
	ByTxSigBE   string `json:"byTxSigBE"`
	ByTxSigLE   string `json:"byTxSigLE"`
	ByTxPadding string `json:"byTxPadding"`
	ByTxPayload string `json:"byTxPayload"`
	ByTxScript  string `json:"byTxScript"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (t *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(*t)
}

////////////////
type Welcome struct {
	PubKey  string `json:"pubKey"`
	Contact string `json:"contact"`
	Job     string `json:"job"`
	Github  string `json:"github"`
}
