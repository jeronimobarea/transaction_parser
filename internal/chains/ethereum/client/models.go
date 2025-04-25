package client

import "encoding/json"

type (
	rpcRequest struct {
		JSONRPC string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		ID      int           `json:"id"`
	}

	rpcResponse struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      int             `json:"id"`
		Result  json.RawMessage `json:"result"`
		Error   *rpcError       `json:"error,omitempty"`
	}

	rpcError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	blockResponse struct {
		Number       string                `json:"number"`
		Transactions []TransactionResponse `json:"transactions"`
	}

	TransactionResponse struct {
		Hash        string `json:"hash"`
		From        string `json:"from"`
		To          string `json:"to"`
		Value       string `json:"value"`
		BlockNumber string `json:"blockNumber"`
	}
)
