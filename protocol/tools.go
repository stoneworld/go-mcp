package protocol

// ListToolsRequest represents a request to list available tools
type ListToolsRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// ListToolsResult represents the response to a list tools request
type ListToolsResult struct {
	Tools      []Tool `json:"tools"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// Tool represents a tool definition that the client can call
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// CallToolRequest represents a request to call a specific tool
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResult represents the response to a tool call
type CallToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// ToolListChangedNotification represents a notification that the tool list has changed
type ToolListChangedNotification struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

// NewListToolsRequest creates a new list tools request
func NewListToolsRequest(cursor string) *ListToolsRequest {
	return &ListToolsRequest{
		Cursor: cursor,
	}
}

// NewListToolsResponse creates a new list tools response
func NewListToolsResponse(id RequestID, tools []Tool, nextCursor string) *ListToolsResult {
	return &ListToolsResult{
		Tools:      tools,
		NextCursor: nextCursor,
	}
}

// NewCallToolRequest creates a new call tool request
func NewCallToolRequest(name string, arguments map[string]interface{}) *CallToolRequest {
	return &CallToolRequest{
		Name:      name,
		Arguments: arguments,
	}
}

// NewCallToolResponse creates a new call tool response
func NewCallToolResponse(id RequestID, content []Content, isError bool) *CallToolResult {
	return &CallToolResult{
		Content: content,
		IsError: isError,
	}
}

// NewToolListChangedNotification creates a new tool list changed notification
func NewToolListChangedNotification() *ToolListChangedNotification {
	return &ToolListChangedNotification{}
}
