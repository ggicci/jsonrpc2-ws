package jsonrpc2ws

import (
	"errors"
	"fmt"

	"github.com/ggicci/jsonrpc2"
)

var (
	ErrMethodAlreadyRegistered = errors.New("duplicate method name")
)

// Handler is a JSON-RPC 2.0 WebSocket handler.
type Router struct {
	methods map[string]jsonrpc2.Handler
}

func NewRouter() *Router {
	return &Router{
		methods: make(map[string]jsonrpc2.Handler),
	}
}

func (rt *Router) RegisterMethod(method string, h jsonrpc2.Handler) error {
	if _, exists := rt.methods[method]; exists {
		return fmt.Errorf("%w: %q", ErrMethodAlreadyRegistered, method)
	}
	rt.methods[method] = h
	return nil
}

// Handle implements jsonrpc2.Handler interface.
func (rt *Router) Handle(conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	h := rt.methods[req.Method]
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
