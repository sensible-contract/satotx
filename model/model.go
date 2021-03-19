package model

import "encoding/json"

type TxRequest struct {
	TxHex   string `json:"txHex"`
	ByTxHex string `json:"byTxHex"`
}

type TxResponse struct {
	TxId    string `json:"txId"`
	Index   int    `json:"index"`
	ByTxId  string `json:"byTxId"`
	SigBE   string `json:"sigBE"`
	SigLE   string `json:"sigLE"`
	Padding string `json:"padding"`
	Payload string `json:"payload"`
	Script  string `json:"script"`
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
