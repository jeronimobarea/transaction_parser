package handlers

import "github.com/jeronimobarea/transaction_parser/internal/parser"

type currentBlockResponse struct {
	BlockNumber int64 `json:"block_number"`
}

func newCurrentBlockResponse(blockNumber int64) *currentBlockResponse {
	return &currentBlockResponse{
		BlockNumber: blockNumber,
	}
}

type transactionResponse struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	BlockNumber string `json:"blockNumber"`
}

func newTransactionResponse(tx parser.Transaction) *transactionResponse {
	return &transactionResponse{
		Hash:        tx.Hash,
		From:        string(tx.From),
		To:          string(tx.To),
		Value:       tx.Value,
		BlockNumber: tx.BlockNumber,
	}
}

func newTransactionsResponse(txs []parser.Transaction) []*transactionResponse {
	txsView := make([]*transactionResponse, len(txs))
	for i := range txs {
		txsView[i] = newTransactionResponse(txs[i])
	}
	return txsView
}
