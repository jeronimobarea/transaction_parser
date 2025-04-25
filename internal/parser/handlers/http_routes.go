package handlers

import (
	"log"

	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/httphandler"
	"github.com/jeronimobarea/transaction_parser/pkg/httpx"
)

type Handler struct {
	httphandler.Handler
	parserSvc parser.Service
	logger    *log.Logger
}

func RegisterRoutes(
	router *httpx.Router,
	parserSvc parser.Service,
	logger *log.Logger,
) {
	handlers := Handler{
		parserSvc: parserSvc,
		logger:    logger,
	}

	router.Handle("GET", "/blocks/current", handlers.getCurrentBlock)
	router.Handle("POST", "/subscribe", handlers.subscribeAddress)
	router.Handle("GET", "/transactions", handlers.getTransactions)
}
