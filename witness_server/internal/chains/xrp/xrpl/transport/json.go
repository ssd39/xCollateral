package transport

import (
	"encoding/json"
	"fmt"
)

// Request is a jsonrpc request
type Request struct {
	Id     uint64 `json:"id"`
	Name   string `json:"command"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
	*Response
}

// ErrorObject is a jsonrpc error
type ErrorObject struct {
	Name    string `json:"error"`
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

// Response is a jsonrpc response
type Response struct {
	*ErrorObject
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result"`
	Status string          `json:"status"`
	Type   string          `json:"type"`
}

// Subscription is a jsonrpc subscription
type Subscription struct {
	ID     string          `json:"subscription"`
	Result json.RawMessage `json:"result"`
}

func (r *Response) HasError() bool {
	return r.ErrorObject != nil
}

// Error implements error interface
func (r *Response) Error() string {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("jsonrpc.internal marshal error: %v", err)
	}
	return string(data)
}
