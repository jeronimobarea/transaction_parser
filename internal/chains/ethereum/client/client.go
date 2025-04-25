package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const (
	JSONRPCVersion = "2.0"

	EthGetBlockByNumber = "eth_getBlockByNumber"

	ReturnFullTransactionObjects = true
	LatestBlock                  = "latest"
)

type Client interface {
	GetBlock(ctx context.Context, blockID string) ([]TransactionResponse, error)
}

type client struct {
	url        string
	httpClient *http.Client
	logger     *log.Logger
}

func NewClient(url string, logger *log.Logger) Client {
	return &client{
		url:        url,
		httpClient: http.DefaultClient,
		logger:     logger,
	}
}

func (c *client) GetBlock(ctx context.Context, blockID string) ([]TransactionResponse, error) {
	params := []interface{}{
		blockID,
		ReturnFullTransactionObjects,
	}

	resp, err := c.doRPCRequest(ctx, EthGetBlockByNumber, params...)
	if err != nil {
		c.logger.Printf("error making get block request: %v\n", err)
		return nil, err
	}

	var block blockResponse
	err = json.Unmarshal([]byte(resp), &block)
	if err != nil {
		c.logger.Printf("error unmarshalling response block response: %v\n", err)
		return nil, err
	}
	return block.Transactions, nil
}

func (c *client) doRPCRequest(ctx context.Context, method string, params ...interface{}) (json.RawMessage, error) {
	payload := rpcRequest{
		JSONRPC: JSONRPCVersion,
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp rpcResponse
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	return resp.Result, nil
}
