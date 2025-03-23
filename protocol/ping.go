package protocol

type PingRequest struct{}

type PingResult struct{}

// NewPingRequest creates a new ping request
func NewPingRequest() *PingRequest {
	return &PingRequest{}
}

// NewPingResponse creates a new ping response
func NewPingResponse() *PingResult {
	return &PingResult{}
}
