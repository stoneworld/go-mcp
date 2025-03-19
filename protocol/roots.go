package protocol

// ListRootsRequest represents a request to list root directories
type ListRootsRequest struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

// ListRootsResult represents the response to a list roots request
type ListRootsResult struct {
	Roots []Root `json:"roots"`
}

// Root represents a root directory or file that the server can operate on
type Root struct {
	Name string `json:"name,omitempty"`
	URI  string `json:"uri"`
}

// RootsListChangedNotification represents a notification that the roots list has changed
type RootsListChangedNotification struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

// NewListRootsRequest creates a new list roots request
func NewListRootsRequest() Params {
	return nil
}

// NewListRootsResponse creates a new list roots response
func NewListRootsResponse(id RequestID, roots []Root) Result {
	return &ListRootsResult{
		Roots: roots,
	}
}

// NewRootsListChangedNotification creates a new roots list changed notification
func NewRootsListChangedNotification() Params {
	return nil
}
