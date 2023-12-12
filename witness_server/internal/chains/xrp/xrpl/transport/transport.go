package transport

import "strings"

// Transport is an inteface for transport methods to send jsonrpc requests
type Transport interface {
	// Call makes a jsonrpc request
	Call(method string, out interface{}, params interface{}) error

	// Close closes the transport connection if necessary
	Close() error
}

// PubSubTransport is a transport that allows subscriptions
type PubSubTransport interface {
	// Subscribe starts a subscription to a new event
	Subscribe(method string, callback func(b []byte)) (func() error, error)
}

// NewTransport creates a new transport object
func NewTransport(url string) (Transport, error) {
	if strings.HasPrefix(url, "ws://") || strings.HasPrefix(url, "wss://") {
		t, err := newWebsocket(url)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	return newHTTP(url), nil
}
