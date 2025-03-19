package protocol

const jsonrpc_version = "2.0"

// Standard JSON-RPC error codes
const (
	PARSE_ERROR      = -32700 // Invalid JSON
	INVALID_REQUEST  = -32600 // The JSON sent is not a valid Request object
	METHOD_NOT_FOUND = -32601 // The method does not exist / is not available
	INVALID_PARAMS   = -32602 // Invalid method parameter(s)
	INTERNAL_ERROR   = -32603 // Internal JSON-RPC error
)

type RequestId interface{}

type Method string

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      RequestId   `json:"id"`
	Method  Method      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// IsValid checks if the request is valid according to JSON-RPC 2.0 spec
func (r *JSONRPCRequest) IsValid() bool {
	return r.JSONRPC == jsonrpc_version && r.Method != ""
}

// JSONRPCResponse represents a successful (non-error) response to a request.
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      RequestId   `json:"id"`
	Result  interface{} `json:"result"`
}

// JSONRPCError represents a non-successful (error) response to a request.
type JSONRPCError struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      RequestId `json:"id"`
	Error   struct {
		// The error type that occurred.
		Code int `json:"code"`
		// A short description of the error. The message SHOULD be limited
		// to a concise single sentence.
		Message string `json:"message"`
		// Additional information about the error. The value of this member
		// is defined by the sender (e.g. detailed error information, nested errors etc.).
		Data interface{} `json:"data,omitempty"`
	} `json:"error"`
}

type JSONRPCNotification struct {
	JSONRPC string `json:"jsonrpc"`
	Notification
}

// Notification represents a JSON-RPC notification
type Notification struct {
	Method Method      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

// NewJSONRPCRequest NewRequest creates a new JSON-RPC request
func NewJSONRPCRequest(id RequestId, method Method, params interface{}) *JSONRPCRequest {
	return &JSONRPCRequest{
		JSONRPC: jsonrpc_version,
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

// NewJSONRPCResponse NewResponse creates a new JSON-RPC response
func NewJSONRPCResponse(id RequestId, result interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: jsonrpc_version,
		ID:      id,
		Result:  result,
	}
}

// NewJSONRPCError NewError creates a new JSON-RPC error response
func NewJSONRPCError(id RequestId, code int, message string) *JSONRPCError {
	err := &JSONRPCError{
		JSONRPC: jsonrpc_version,
		ID:      id,
	}
	err.Error.Code = code
	err.Error.Message = message
	return err
}

// NewJSONRPCNotification NewNotification creates a new JSON-RPC notification
func NewJSONRPCNotification(method Method, params interface{}) *JSONRPCNotification {
	return &JSONRPCNotification{
		JSONRPC: jsonrpc_version,
		Notification: Notification{
			Method: method,
			Params: params,
		},
	}
}
