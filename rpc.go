// Based on JSON-RPC 2.0 (https://www.jsonrpc.org/specification)

package gv

type RpcMethod struct {
	Name    string
	Handler func(RpcRequest) RpcResponse
}

type RpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	Id      string `json:"id"`
}

type RpcErrorCode int

const (
	RpcErrorCodeParseError     RpcErrorCode = -32700
	RpcErrorCodeInvalidRequest RpcErrorCode = -32600
	RpcErrorCodeMethodNotFound RpcErrorCode = -32601
	RpcErrorCodeInvalidParams  RpcErrorCode = -32602
	RpcErrorCodeInternalError  RpcErrorCode = -32603
	RpcErrorCodeServerError    RpcErrorCode = -32000
)

type RpcError struct {
	Code    RpcErrorCode `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data"`
}

type RpcResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	Result  *any      `json:"result,omitempty"`
	Error   *RpcError `json:"error,omitempty"`
	Id      string    `json:"id"`
}

type RpcServer struct{}

func (r *RpcServer) Handle() {}

func NewRpcServer() *RpcServer {
	return &RpcServer{}
}
