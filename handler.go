package jsonrpc2ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ggicci/jsonrpc2"
	wsstream "github.com/ggicci/jsonrpc2/websocket"
	"github.com/gorilla/websocket"
)

var (
	ErrMethodAlreadyRegistered = errors.New("duplicate method name")
)

type JSONRPC2Handler struct {
	methods map[string]jsonrpc2.Handler
}

func NewJSONRPC2Handler() *JSONRPC2Handler {
	return &JSONRPC2Handler{
		methods: make(map[string]jsonrpc2.Handler),
	}
}

func (jh *JSONRPC2Handler) RegisterMethod(method string, h jsonrpc2.Handler) error {
	if _, exists := jh.methods[method]; exists {
		return fmt.Errorf("%w: %q", ErrMethodAlreadyRegistered, method)
	}
	jh.methods[method] = h
	return nil
}

func (jh *JSONRPC2Handler) MustRegisterMethod(method string, h jsonrpc2.Handler) {
	panicOnError(jh.RegisterMethod(method, h))
}

// Handle implements jsonrpc2.Handler interface.
func (jh *JSONRPC2Handler) Handle(conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	h := jh.methods[req.Method]
	if h == nil {
		if !req.Notif {
			conn.ReplyWithError(req.Context(), req.ID, &jsonrpc2.Error{
				Code:    jsonrpc2.CodeMethodNotFound,
				Message: fmt.Sprintf("method %q not found", req.Method),
			})
		}
		return
	}

	h.Handle(conn, req)
}

// ServeHTTP implements http.Handler interface.
func (jh *JSONRPC2Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		http.Error(rw, fmt.Errorf("could not upgrade to WebSocket: %w", err).Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	<-jsonrpc2.NewConn(
		context.Background(),
		wsstream.NewObjectStream(conn),
		jh,
	).DisconnectNotify()
}
