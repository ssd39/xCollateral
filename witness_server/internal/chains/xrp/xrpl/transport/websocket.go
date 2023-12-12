package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"peersyst/bridge-witness-go/internal/common/utils"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func newWebsocket(url string) (Transport, error) {
	wsConn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		return nil, err
	}
	codec := &websocketCodec{
		conn: wsConn,
	}
	return newStream(codec)
}

// ErrTimeout happens when the websocket requests times out
var ErrTimeout = fmt.Errorf("timeout")

type ackMessage struct {
	buf []byte
	err error
}

type callback func(b []byte, err error)

type stream struct {
	seq   uint64
	codec Codec

	// call handlers
	handlerLock sync.Mutex
	handler     map[uint64]callback

	// subscriptions
	subsLock sync.Mutex
	subs     map[string]func(b []byte)

	closeCh chan struct{}
	timer   *time.Timer
}

func newStream(codec Codec) (*stream, error) {
	w := &stream{
		codec:   codec,
		closeCh: make(chan struct{}),
		handler: map[uint64]callback{},
		subs:    map[string]func(b []byte){},
	}

	go w.listen()
	return w, nil
}

// Close implements the the transport interface
func (s *stream) Close() error {
	close(s.closeCh)
	return s.codec.Close()
}

func (s *stream) incSeq() uint64 {
	return atomic.AddUint64(&s.seq, 1)
}

func (s *stream) isClosed() bool {
	select {
	case <-s.closeCh:
		return true
	default:
		return false
	}
}

func (s *stream) listen() {
	buf := []byte{}

	for {
		var err error
		buf, err = s.codec.Read(buf[:0])
		if err != nil {
			if !s.isClosed() {
				log.Error().Msgf("Error reading buffer: '%+v'", err)
			}
			continue
		}

		var resp Response
		if err = json.Unmarshal(buf, &resp); err != nil {
			log.Error().Msgf("Error unmarshal response buffer to Response struct: %+v", err)
			continue
		}

		if resp.ID != 0 {
			go s.handleMsg(resp)
		} else {
			// handle subscription
			var respSub Request
			if err = json.Unmarshal(buf, &respSub); err != nil {
				continue
			}

			if respSub.Name == "subscribe" {
				go s.handleSubscription(respSub)
			}
		}
	}
}

func (s *stream) handleSubscription(response Request) {
	var sub Subscription
	if err := json.Unmarshal(response.Result, &sub); err != nil {
		panic(err)
	}

	s.subsLock.Lock()
	callback, ok := s.subs[sub.ID]
	s.subsLock.Unlock()

	if !ok {
		return
	}

	// call the callback function
	callback(sub.Result)
}

func (s *stream) handleMsg(response Response) {
	s.handlerLock.Lock()
	callback, ok := s.handler[response.ID]
	if !ok {
		s.handlerLock.Unlock()
		return
	}

	// delete handler
	delete(s.handler, response.ID)
	s.handlerLock.Unlock()

	if response.HasError() {
		callback(nil, &response)
	} else {
		callback(response.Result, nil)
	}
}

var timerLocker sync.Mutex

func (s *stream) setHandler(id uint64, ack chan *ackMessage) {
	callback := func(b []byte, err error) {
		select {
		case ack <- &ackMessage{b, err}:
		default:
		}
	}

	s.handlerLock.Lock()
	s.handler[id] = callback
	s.handlerLock.Unlock()

	timerLocker.Lock()
	defer timerLocker.Unlock()
	s.timer = time.AfterFunc(5*time.Second, func() {
		s.handlerLock.Lock()
		delete(s.handler, id)
		s.handlerLock.Unlock()

		select {
		case ack <- &ackMessage{nil, ErrTimeout}:
		default:
		}
	})
}

// Call implements the transport interface
func (s *stream) Call(method string, out interface{}, params interface{}) error {
	seq := s.incSeq()
	request := &Request{
		Id:   seq,
		Name: method,
	}
	var fullRequestMap map[string]interface{}
	ja, err := json.Marshal(request)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(ja, &fullRequestMap); err != nil {
		return err
	}

	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &fullRequestMap); err != nil {
			return err
		}
	}

	ack := make(chan *ackMessage)
	s.setHandler(seq, ack)

	raw, err := json.Marshal(fullRequestMap)
	if err != nil {
		return err
	}
	if err := s.codec.Write(raw); err != nil {
		return err
	}

	resp := <-ack
	if resp.err != nil {
		return resp.err
	}
	log.Debug().Msgf("Response buffer in string: '%+v'", utils.Truncate(string(resp.buf[:]), 2048))
	if err := json.Unmarshal(resp.buf, out); err != nil {
		return err
	}
	return nil
}

func (s *stream) unsubscribe(id string) error {
	s.subsLock.Lock()
	defer s.subsLock.Unlock()

	if _, ok := s.subs[id]; !ok {
		return fmt.Errorf("subscription %s not found", id)
	}
	delete(s.subs, id)

	var result bool
	if err := s.Call("eth_unsubscribe", &result, id); err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("failed to unsubscribe")
	}
	return nil
}

func (s *stream) setSubscription(id string, callback func(b []byte)) {
	s.subsLock.Lock()
	defer s.subsLock.Unlock()

	s.subs[id] = callback
}

// Subscribe implements the PubSubTransport interface
func (s *stream) Subscribe(method string, callback func(b []byte)) (func() error, error) {
	var out string
	if err := s.Call("subscribe", &out, method); err != nil {
		return nil, err
	}

	s.setSubscription(out, callback)
	cancel := func() error {
		return s.unsubscribe(out)
	}
	return cancel, nil
}

type websocketCodec struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (w *websocketCodec) Close() error {
	return w.conn.Close()
}

func (w *websocketCodec) Write(b []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteMessage(websocket.TextMessage, b)
}

func (w *websocketCodec) Read(b []byte) ([]byte, error) {
	_, buf, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	b = append(b, buf...)
	return b, nil
}

// Codec is the codec to write and read messages
type Codec interface {
	Read([]byte) ([]byte, error)
	Write([]byte) error
	Close() error
}
