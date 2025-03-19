package protocol

// NewPingRequest creates a new ping request
func NewPingRequest() Params {
	return nil
}

// NewPingResponse creates a new ping response
func NewPingResponse() Result {
	return struct{}{}
}
