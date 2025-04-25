package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/test/evmtest"
	"github.com/jeronimobarea/transaction_parser/internal/test/parsertest"
)

func TestHandler_GetCurrentBlockHandler(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		fake := &parsertest.FakeParserSvc{GetCurrentBlockResp: 123}
		h := Handler{
			parserSvc: fake,
			logger:    log.Default(),
		}

		req := httptest.NewRequest("GET", "/current-block", nil)
		rec := httptest.NewRecorder()

		h.getCurrentBlock(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}

		var resp currentBlockResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if resp.BlockNumber != 123 {
			t.Errorf("expected CurrentBlock=123, got %d", resp.BlockNumber)
		}
	})

	t.Run("parser error", func(t *testing.T) {
		fake := &parsertest.FakeParserSvc{GetCurrentBlockErr: errors.New("rpc fail")}
		h := Handler{
			parserSvc: fake,
			logger:    log.Default(),
		}

		req := httptest.NewRequest("GET", "/current-block", nil)
		rec := httptest.NewRecorder()

		h.getCurrentBlock(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestHandler_SubscribeAddress(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		fake := &parsertest.FakeParserSvc{SubscribeErr: nil}
		h := Handler{
			parserSvc: fake,
			logger:    log.Default(),
		}

		url := "/subscribe?" + AddressQueryKey + "=0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		req := httptest.NewRequest("POST", url, nil)
		rec := httptest.NewRecorder()

		h.subscribeAddress(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}
		if body := rec.Body.String(); body != "{}\n" {
			t.Errorf("expected empty body, got %q", body)
		}
	})

	t.Run("parser error", func(t *testing.T) {
		fake := &parsertest.FakeParserSvc{SubscribeErr: errors.New("subscribe fail")}
		h := Handler{
			parserSvc: fake,
			logger:    log.Default(),
		}

		url := "/subscribe?" + AddressQueryKey + "=0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		req := httptest.NewRequest("POST", url, nil)
		rec := httptest.NewRecorder()

		h.subscribeAddress(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestHandler_GetTransactions(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		wantTxs := []parser.Transaction{
			{Hash: "h1", From: evmtest.EVMZeroValueAddress, To: evmtest.EVMZeroValueAddress, Value: "10", BlockNumber: "0x1"},
		}
		fake := &parsertest.FakeParserSvc{GetTransactionsResp: wantTxs}
		h := Handler{
			parserSvc: fake,
			logger:    log.Default(),
		}

		url := "/transactions?" + AddressQueryKey + "=" + evmtest.EVMZeroValueAddress.String()
		req := httptest.NewRequest("GET", url, nil)
		rec := httptest.NewRecorder()

		h.getTransactions(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}

		var resp []*transactionResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if len(resp) != len(wantTxs) {
			t.Errorf("exepceted equal transactions length %d but got %d", len(resp), len(wantTxs))
		}
		if wantTxs[0].BlockNumber != resp[0].BlockNumber {
			t.Errorf("expected transactions %+v, got %+v", wantTxs, resp)
		}
	})

	t.Run("parser error", func(t *testing.T) {
		fake := &parsertest.FakeParserSvc{GetTransactionsErr: errors.New("fetch fail")}
		h := Handler{
			parserSvc: fake,
			logger:    log.Default(),
		}

		url := "/transactions?" + AddressQueryKey + "=" + evmtest.EVMZeroValueAddress.String()
		req := httptest.NewRequest("GET", url, nil)
		rec := httptest.NewRecorder()

		h.getTransactions(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d on service error, got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}
