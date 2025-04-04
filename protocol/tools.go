package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/pkg"
)

// ListToolsRequest represents a request to list available tools
type ListToolsRequest struct{}

// ListToolsResult represents the response to a list tools request
type ListToolsResult struct {
	Tools      []*Tool `json:"tools"`
	NextCursor string  `json:"nextCursor,omitempty"`
}

// Tool represents a tool definition that the client can call
type Tool struct {
	// Name is the unique identifier of the tool
	Name string `json:"name"`

	// Description is a human-readable description of the tool
	Description string `json:"description,omitempty"`

	// InputSchema defines the expected parameters for the tool using JSON Schema
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchemaType string

const Object InputSchemaType = "object"

// InputSchema represents a JSON Schema object defining the expected parameters for a tool
type InputSchema struct {
	Type       InputSchemaType        `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
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

// UnmarshalJSON implements the json.Unmarshaler interface for CallToolResult
func (r *CallToolResult) UnmarshalJSON(data []byte) error {
	type Alias CallToolResult
	aux := &struct {
		Content []json.RawMessage `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := pkg.JsonUnmarshal(data, &aux); err != nil {
		return err
	}

	r.Content = make([]Content, len(aux.Content))
	for i, content := range aux.Content {
		// Try to unmarshal content as TextContent first
		var textContent TextContent
		if err := pkg.JsonUnmarshal(content, &textContent); err == nil {
			r.Content[i] = textContent
			continue
		}

		// Try to unmarshal content as ImageContent
		var imageContent ImageContent
		if err := pkg.JsonUnmarshal(content, &imageContent); err == nil {
			r.Content[i] = imageContent
			continue
		}

		// Try to unmarshal content as embeddedResource
		var embeddedResource EmbeddedResource
		if err := pkg.JsonUnmarshal(content, &embeddedResource); err == nil {
			r.Content[i] = embeddedResource
			return nil
		}

		return fmt.Errorf("unknown content type at index %d", i)
	}

	return nil
}

// ToolListChangedNotification represents a notification that the tool list has changed
type ToolListChangedNotification struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

// NewListToolsRequest creates a new list tools request
func NewListToolsRequest() *ListToolsRequest {
	return &ListToolsRequest{}
}

// NewListToolsResult creates a new list tools response
func NewListToolsResult(tools []*Tool, nextCursor string) *ListToolsResult {
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

// NewCallToolResult creates a new call tool response
func NewCallToolResult(content []Content, isError bool) *CallToolResult {
	return &CallToolResult{
		Content: content,
		IsError: isError,
	}
}

// NewToolListChangedNotification creates a new tool list changed notification
func NewToolListChangedNotification() *ToolListChangedNotification {
	return &ToolListChangedNotification{}
}
