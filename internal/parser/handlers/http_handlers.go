package handlers

import (
	"net/http"

	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
)

const AddressQueryKey = "address"

func (h Handler) getCurrentBlock(w http.ResponseWriter, r *http.Request) {
	blockNumber, err := h.parserSvc.GetCurrentBlock(r.Context())
	if err != nil {
		h.logger.Printf("error retrieving current block: %v", err)
		h.HandleError(w, err)
		return
	}

	h.OK(w, newCurrentBlockResponse(blockNumber))
}

func (h Handler) subscribeAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get(AddressQueryKey)

	err := h.parserSvc.Subscribe(r.Context(), evm.Address(address))
	if err != nil {
		h.logger.Printf("error subscribing address: %s: %v", address, err)

		h.HandleError(w, err)
		return
	}

	h.OK(w, struct{}{})
}

func (h Handler) getTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get(AddressQueryKey)

	txs, err := h.parserSvc.GetTransactions(r.Context(), evm.Address(address))
	if err != nil {
		h.logger.Printf("error retrieving transactions for address: %s: %v", address, err)

		h.HandleError(w, err)
		return
	}

	h.OK(w, newTransactionsResponse(txs))
}
