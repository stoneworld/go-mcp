package protocol

// Method represents the JSON-RPC method name
type Method string

const (
	// Core methods
	Ping                    Method = "ping"
	Initialize              Method = "initialize"
	NotificationInitialized Method = "notifications/initialized"

	// Root related methods
	RootsList                    Method = "roots/list"
	NotificationRootsListChanged Method = "notifications/roots/list_changed"

	// Resource related methods
	ResourcesList                    Method = "resources/list"
	ResourceListTemplates            Method = "resources/templates/list"
	ResourcesRead                    Method = "resources/read"
	ResourcesSubscribe               Method = "resources/subscribe"
	ResourcesUnsubscribe             Method = "resources/unsubscribe"
	NotificationResourcesListChanged Method = "notifications/resources/list_changed"
	NotificationResourcesUpdated     Method = "notifications/resources/updated"

	// Tool related methods
	ToolsList                    Method = "tools/list"
	ToolsCall                    Method = "tools/call"
	NotificationToolsListChanged Method = "notifications/tools/list_changed"

	// Prompt related methods
	PromptsList                    Method = "prompts/list"
	PromptsGet                     Method = "prompts/get"
	NotificationPromptsListChanged Method = "notifications/prompts/list_changed"

	// Sampling related methods
	SamplingCreateMessage Method = "sampling/createMessage"

	// Logging related methods
	LoggingSetLevel Method = "logging/setLevel"
	LogMessage      Method = "notifications/message"

	// Completion related methods
	CompletionComplete Method = "completion/complete"

	// progress related methods
	NotificationProgress  Method = "notifications/progress"
	NotificationCancelled Method = "notifications/cancelled"
)

// Role represents the sender or recipient of messages and data in a conversation
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)
