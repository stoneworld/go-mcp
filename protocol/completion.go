package protocol

// CompleteRequest represents a request for completion options
type CompleteRequest struct {
	Argument struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"argument"`
	Ref interface{} `json:"ref"` // Can be PromptReference or ResourceReference
}

// Reference types
type PromptReference struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type ResourceReference struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

// CompleteResult represents the response to a completion request
type CompleteResult struct {
	Completion struct {
		Values  []string `json:"values"`
		HasMore bool     `json:"hasMore,omitempty"`
		Total   int      `json:"total,omitempty"`
	} `json:"completion"`
}

// NewCompleteRequest creates a new completion request
func NewCompleteRequest(argName string, argValue string, ref interface{}) Params {
	return CompleteRequest{
		Argument: struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			Name:  argName,
			Value: argValue,
		},
		Ref: ref,
	}
}

// NewCompleteResponse creates a new completion response
func NewCompleteResponse(values []string, hasMore bool, total int) Result {
	return CompleteResult{
		Completion: struct {
			Values  []string `json:"values"`
			HasMore bool     `json:"hasMore,omitempty"`
			Total   int      `json:"total,omitempty"`
		}{
			Values:  values,
			HasMore: hasMore,
			Total:   total,
		},
	}
}
