package protocol

import (
	"encoding/json"
	"fmt"

	"go-mcp/pkg"
)

// ListResourcesRequest represents a request to list available resources
type ListResourcesRequest struct{}

// ListResourcesResult represents the response to a list resources request
type ListResourcesResult struct {
	Resources  []Resource `json:"resources"`
	NextCursor string     `json:"nextCursor,omitempty"`
}

// ListResourceTemplatesRequest represents a request to list resource templates
type ListResourceTemplatesRequest struct{}

// ListResourceTemplatesResult represents the response to a list resource templates request
type ListResourceTemplatesResult struct {
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates"`
	NextCursor        string             `json:"nextCursor,omitempty"`
}

// ReadResourceRequest represents a request to read a specific resource
type ReadResourceRequest struct {
	URI string `json:"uri"`
}

// ReadResourceResult represents the response to a read resource request
type ReadResourceResult struct {
	Contents []ResourceContents `json:"contents"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for ReadResourceResult
func (r *ReadResourceResult) UnmarshalJSON(data []byte) error {
	type Alias ReadResourceResult
	aux := &struct {
		Contents []json.RawMessage `json:"contents"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := pkg.JsonUnmarshal(data, &aux); err != nil {
		return err
	}

	r.Contents = make([]ResourceContents, len(aux.Contents))
	for i, content := range aux.Contents {
		// Try to unmarshal content as TextResourceContents first
		var textContent TextResourceContents
		if err := pkg.JsonUnmarshal(content, &textContent); err == nil {
			r.Contents[i] = textContent
			continue
		}

		// Try to unmarshal content as BlobResourceContents
		var blobContent BlobResourceContents
		if err := pkg.JsonUnmarshal(content, &blobContent); err == nil {
			r.Contents[i] = blobContent
			continue
		}

		return fmt.Errorf("unknown content type at index %d", i)
	}

	return nil
}

// Resource related types
type Resource struct {
	Annotated
	Name        string `json:"name"`
	URI         string `json:"uri"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
	Size        int64  `json:"size,omitempty"`
}

type ResourceTemplate struct {
	Annotated
	Name        string `json:"name"`
	URITemplate string `json:"uriTemplate"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// Annotated represents base objects that include optional annotations
type Annotated struct {
	Annotations *Annotations `json:"annotations,omitempty"`
}

// Annotations represents optional annotations for objects
type Annotations struct {
	Audience []Role  `json:"audience,omitempty"`
	Priority float64 `json:"priority,omitempty"`
}

// ModelHint represents hints to use for model selection
type ModelHint struct {
	Name string `json:"name,omitempty"`
}

// ModelPreferences represents the server's preferences for model selection
type ModelPreferences struct {
	CostPriority         float64     `json:"costPriority,omitempty"`
	IntelligencePriority float64     `json:"intelligencePriority,omitempty"`
	SpeedPriority        float64     `json:"speedPriority,omitempty"`
	Hints                []ModelHint `json:"hints,omitempty"`
}

// Content interfaces and types
type Content interface {
	GetType() string
}

type TextContent struct {
	Annotated
	Type string `json:"type"`
	Text string `json:"text"`
}

func (t TextContent) GetType() string {
	return "text"
}

type ImageContent struct {
	Annotated
	Type     string `json:"type"`
	Data     []byte `json:"data"`
	MimeType string `json:"mimeType"`
}

func (i ImageContent) GetType() string {
	return "image"
}

// EmbeddedResource represents the contents of a resource, embedded into a prompt or tool call result.
// It is up to the client how best to render embedded resources for the benefit of the LLM and/or the user.
type EmbeddedResource struct {
	Type        string           `json:"type"` // Must be "resource"
	Resource    ResourceContents `json:"resource"`
	Annotations *Annotations     `json:"annotations,omitempty"`
}

// NewEmbeddedResource creates a new EmbeddedResource
func NewEmbeddedResource(resource ResourceContents, annotations *Annotations) *EmbeddedResource {
	return &EmbeddedResource{
		Type:        "resource",
		Resource:    resource,
		Annotations: annotations,
	}
}

func (i EmbeddedResource) GetType() string {
	return "resource"
}

type ResourceContents interface {
	GetURI() string
	GetMimeType() string
}

type TextResourceContents struct {
	URI      string `json:"uri"`
	Text     string `json:"text"`
	MimeType string `json:"mimeType,omitempty"`
}

func (t TextResourceContents) GetURI() string {
	return t.URI
}

func (t TextResourceContents) GetMimeType() string {
	return t.MimeType
}

type BlobResourceContents struct {
	URI      string `json:"uri"`
	Blob     []byte `json:"blob"`
	MimeType string `json:"mimeType,omitempty"`
}

func (b BlobResourceContents) GetURI() string {
	return b.URI
}

func (b BlobResourceContents) GetMimeType() string {
	return b.MimeType
}

// SubscribeRequest represents a request to subscribe to resource updates
type SubscribeRequest struct {
	URI string `json:"uri"`
}

// UnsubscribeRequest represents a request to unsubscribe from resource updates
type UnsubscribeRequest struct {
	URI string `json:"uri"`
}

type SubscribeResult struct{}

type UnsubscribeResult struct{}

// ResourceListChangedNotification represents a notification that the resource list has changed
type ResourceListChangedNotification struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

// ResourceUpdatedNotification represents a notification that a resource has been updated
type ResourceUpdatedNotification struct {
	URI string `json:"uri"`
}

// NewListResourcesRequest creates a new list resources request
func NewListResourcesRequest() *ListResourcesRequest {
	return &ListResourcesRequest{}
}

// NewListResourcesResponse creates a new list resources response
func NewListResourcesResponse(resources []Resource, nextCursor string) *ListResourcesResult {
	return &ListResourcesResult{
		Resources:  resources,
		NextCursor: nextCursor,
	}
}

// NewListResourceTemplatesRequest creates a new list resource templates request
func NewListResourceTemplatesRequest() *ListResourceTemplatesRequest {
	return &ListResourceTemplatesRequest{}
}

// NewListResourceTemplatesResponse creates a new list resource templates response
func NewListResourceTemplatesResponse(templates []ResourceTemplate, nextCursor string) *ListResourceTemplatesResult {
	return &ListResourceTemplatesResult{
		ResourceTemplates: templates,
		NextCursor:        nextCursor,
	}
}

// NewReadResourceRequest creates a new read resource request
func NewReadResourceRequest(uri string) *ReadResourceRequest {
	return &ReadResourceRequest{URI: uri}
}

// NewReadResourceResponse creates a new read resource response
func NewReadResourceResponse(contents []ResourceContents) *ReadResourceResult {
	return &ReadResourceResult{
		Contents: contents,
	}
}

// NewSubscribeRequest creates a new subscribe request
func NewSubscribeRequest(uri string) *SubscribeRequest {
	return &SubscribeRequest{URI: uri}
}

// NewUnsubscribeRequest creates a new unsubscribe request
func NewUnsubscribeRequest(uri string) *UnsubscribeRequest {
	return &UnsubscribeRequest{URI: uri}
}

func NewSubscribeResponse() *SubscribeResult {
	return &SubscribeResult{}
}

func NewUnsubscribeResponse() *UnsubscribeResult {
	return &UnsubscribeResult{}
}

// NewResourceListChangedNotification creates a new resource list changed notification
func NewResourceListChangedNotification() *ResourceListChangedNotification {
	return &ResourceListChangedNotification{}
}

// NewResourceUpdatedNotification creates a new resource updated notification
func NewResourceUpdatedNotification(uri string) *ResourceUpdatedNotification {
	return &ResourceUpdatedNotification{URI: uri}
}
