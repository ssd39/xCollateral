package transport

// HTTP is an http transport
type HTTP struct {
	rpcUrl string
}

func newHTTP(rpcUrl string) *HTTP {
	return &HTTP{
		rpcUrl: rpcUrl,
	}
}

// Close implements the transport interface
func (h *HTTP) Close() error {
	return nil
}

// Call implements the transport interface
func (h *HTTP) Call(method string, out interface{}, params interface{}) error {
	return nil
}
