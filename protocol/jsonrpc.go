package protocol

import "encoding/json"

const jsonrpc_version = "2.0"

// Standard JSON-RPC error codes
const (
	PARSE_ERROR      = -32700 // Invalid JSON
	INVALID_REQUEST  = -32600 // The JSON sent is not a valid Request object
	METHOD_NOT_FOUND = -32601 // The method does not exist / is not available
	INVALID_PARAMS   = -32602 // Invalid method parameter(s)
	INTERNAL_ERROR   = -32603 // Internal JSON-RPC error
)

type RequestID interface{} // 字符串/数值

type Params interface{}

type Result interface{}

type JSONRPCRequest struct {
	JSONRPC   string          `json:"jsonrpc"`
	ID        RequestID       `json:"id"`
	Method    Method          `json:"method"`
	Params    Params          `json:"params,omitempty"`
	RawParams json.RawMessage `json:"-"`
}

func (r *JSONRPCRequest) UnmarshalJSON(data []byte) error {
	type alias JSONRPCRequest
	temp := &struct {
		Params json.RawMessage `json:"params,omitempty"`
		*alias
	}{
		alias: (*alias)(r),
	}

	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	r.RawParams = temp.Params

	if len(r.RawParams) != 0 {
		if err := json.Unmarshal(r.RawParams, &r.Params); err != nil {
			return err
		}
	}

	return nil
}

// IsValid checks if the request is valid according to JSON-RPC 2.0 spec
func (r *JSONRPCRequest) IsValid() bool {
	return r.JSONRPC == jsonrpc_version && r.Method != ""
}

// JSONRPCResponse represents a response to a request.
type JSONRPCResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      RequestID `json:"id"`
	Result  Result    `json:"result"`
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

type JSONRPCResponseErr struct {
}

type JSONRPCNotification struct {
	JSONRPC string `json:"jsonrpc"`
	Method  Method `json:"method"`
	Params  Params `json:"params,omitempty"`
}

// NewJSONRPCRequest creates a new JSON-RPC request
func NewJSONRPCRequest(id RequestID, method Method, params Params) *JSONRPCRequest {
	return &JSONRPCRequest{
		JSONRPC: jsonrpc_version,
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

// NewJSONRPCSuccessResponse creates a new JSON-RPC response
func NewJSONRPCSuccessResponse(id RequestID, result Result) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: jsonrpc_version,
		ID:      id,
		Result:  result,
	}
}

// NewJSONRPCErrorResponse NewError creates a new JSON-RPC error response
func NewJSONRPCErrorResponse(id RequestID, code int, message string) *JSONRPCResponse {
	err := &JSONRPCResponse{
		JSONRPC: jsonrpc_version,
		ID:      id,
	}
	err.Error.Code = code
	err.Error.Message = message
	return err
}

// NewJSONRPCNotification creates a new JSON-RPC notification
func NewJSONRPCNotification(method Method, params Params) *JSONRPCNotification {
	return &JSONRPCNotification{
		JSONRPC: jsonrpc_version,
		Method:  method,
		Params:  params,
	}
}
