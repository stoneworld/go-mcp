package protocol

// ListPromptsRequest represents a request to list available prompts
type ListPromptsRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// ListPromptsResult represents the response to a list prompts request
type ListPromptsResult struct {
	Prompts    []Prompt `json:"prompts"`
	NextCursor string   `json:"nextCursor,omitempty"`
}

// Prompt related types
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// GetPromptRequest represents a request to get a specific prompt
type GetPromptRequest struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments,omitempty"`
}

// GetPromptResult represents the response to a get prompt request
type GetPromptResult struct {
	Messages    []PromptMessage `json:"messages"`
	Description string          `json:"description,omitempty"`
}

type PromptMessage struct {
	Role    Role    `json:"role"`
	Content Content `json:"content"`
}

// PromptListChangedNotification represents a notification that the prompt list has changed
type PromptListChangedNotification struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

// NewListPromptsRequest creates a new list prompts request
func NewListPromptsRequest(cursor string) *ListPromptsRequest {
	return &ListPromptsRequest{Cursor: cursor}
}

// NewListPromptsResponse creates a new list prompts response
func NewListPromptsResponse(prompts []Prompt, nextCursor string) *ListPromptsResult {
	return &ListPromptsResult{
		Prompts:    prompts,
		NextCursor: nextCursor,
	}
}

// NewGetPromptRequest creates a new get prompt request
func NewGetPromptRequest(name string, arguments map[string]string) *GetPromptRequest {
	return &GetPromptRequest{
		Name:      name,
		Arguments: arguments,
	}
}

// NewGetPromptResponse creates a new get prompt response
func NewGetPromptResponse(messages []PromptMessage, description string) *GetPromptResult {
	return &GetPromptResult{
		Messages:    messages,
		Description: description,
	}
}

// NewPromptListChangedNotification creates a new prompt list changed notification
func NewPromptListChangedNotification() *PromptListChangedNotification {
	return &PromptListChangedNotification{}
}
