package jsonrpc2ws

import (
	"net/http"

	"github.com/ggicci/jsonrpc2"
	"github.com/gorilla/websocket"
)

// Server is a JSON-RPC 2.0 WebSocket server.
type Server struct {
	http.Server
	jsonrpc2.Handler
}

func (s *Server) Serve() error {
	var (
		mux      = http.NewServeMux()
		upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	)

}
