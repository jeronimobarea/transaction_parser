package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jeronimobarea/transaction_parser/internal/pkg/svcerrors"
)

type Handler struct{}

func (h Handler) OK(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.HandleError(w, err)
		return
	}
}

func (h Handler) HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, svcerrors.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}

	if errors.Is(err, svcerrors.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
