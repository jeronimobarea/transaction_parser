package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetCurrentBlock_Success(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c, teardown := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req rpcRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Fatalf("unmarshal request: %v", err)
			}

			if req.Method != EthBlockNumberMethod {
				t.Errorf("want method %q, got %q", EthBlockNumberMethod, req.Method)
			}

			if len(req.Params) != 0 {
				t.Errorf("expected no params, got %v", req.Params)
			}

			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":"0x1a"}`)
		}))
		defer teardown()

		got, err := c.GetCurrentBlock(context.Background())
		if err != nil {
			t.Fatalf("GetCurrentBlock error: %v", err)
		}

		const want = "0x1a"
		if got != want {
			t.Errorf("GetCurrentBlock = %s; want %s", got, want)
		}
	})

	t.Run("rpc error", func(t *testing.T) {
		c, teardown := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"Method not found"}}`)
		}))
		defer teardown()

		_, err := c.GetCurrentBlock(context.Background())
		if err == nil || err.Error() != "Method not found" {
			t.Errorf("expected RPC error \"Method not found\", got %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		c, teardown := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `invalid-json`)
		}))
		defer teardown()

		_, err := c.GetCurrentBlock(context.Background())
		if err == nil {
			t.Fatal("expected JSON unmarshal error, got nil")
		}
	})
}

func TestGetBlock_Success(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var (
			blockID     = "0x2"
			expectedTxs = []TransactionResponse{
				{
					Hash:        "h1",
					From:        "0xfoo",
					To:          "0xbar",
					Value:       "123",
					BlockNumber: "0x2",
				},
			}
		)

		cli, teardown := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req rpcRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Fatalf("unmarshal request: %v", err)
			}
			if req.Method != EthGetBlockByNumber {
				t.Errorf("want method %q, got %q", EthGetBlockByNumber, req.Method)
			}
			// params: [ "0x2", true ]
			if len(req.Params) != 2 {
				t.Fatalf("expected 2 params, got %d", len(req.Params))
			}
			if req.Params[0] != "0x2" {
				t.Errorf("want first param \"0x2\", got %v", req.Params[0])
			}
			if req.Params[1] != ReturnFullTransactionObjects {
				t.Errorf("want second param %v, got %v", ReturnFullTransactionObjects, req.Params[1])
			}

			resp := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"result": map[string]interface{}{
					"number":       "0x2",
					"transactions": expectedTxs,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer teardown()

		txs, err := cli.GetBlock(context.Background(), blockID)
		if err != nil {
			t.Fatalf("GetBlock error: %v", err)
		}
		if !reflect.DeepEqual(txs, expectedTxs) {
			t.Errorf("GetBlock = %+v; want %+v", txs, expectedTxs)
		}
	})

	t.Run("rpc error", func(t *testing.T) {
		cli, teardown := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"block not found"}}`)
		}))
		defer teardown()

		_, err := cli.GetBlock(context.Background(), "0x5")
		if err == nil || err.Error() != "block not found" {
			t.Errorf("expected RPC error \"block not found\", got %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		cli, teardown := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":INVALID}`)
		}))
		defer teardown()

		_, err := cli.GetBlock(context.Background(), "0x1")
		if err == nil {
			t.Fatal("expected JSON unmarshal error, got nil")
		}
	})
}

func newTestClient(handler http.Handler) (Client, func()) {
	ts := httptest.NewServer(handler)
	cli := NewClient(ts.URL, log.Default())
	return cli, ts.Close
}
