# jsonrpc2-ws

JSON-RPC 2.0 over WebSocket for Go

## Usage

```go
import (
    "github.com/ggicci/jsonrpc2-ws"
    "github.com/ggicci/jsonrpc2"
)

websocketHandler := jsonrpc2ws.NewJSONRPC2Handler()
func CloneRepo(conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
    // ...
}

func DeployApp(conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
    // ...
}

websocketHandler.RegisterMethod("clone", jsonrpc2.HandlerFunc(CloneRepo))
websocketHandler.RegisterMethod("deploy", jsonrpc2.HandlerFunc(DeployApp))

http.Handle("/", websocketHanler)
http.ListenAndServe(":8080", nil)
```
