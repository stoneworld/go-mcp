package protocol

// LoggingLevel represents the severity of a log message
type LoggingLevel string

const (
	LogEmergency LoggingLevel = "emergency"
	LogAlert     LoggingLevel = "alert"
	LogCritical  LoggingLevel = "critical"
	LogError     LoggingLevel = "error"
	LogWarning   LoggingLevel = "warning"
	LogNotice    LoggingLevel = "notice"
	LogInfo      LoggingLevel = "info"
	LogDebug     LoggingLevel = "debug"
)

// SetLoggingLevelRequest represents a request to set the logging level
type SetLoggingLevelRequest struct {
	Level LoggingLevel `json:"level"`
}

// SetLoggingLevelResult represents the response to a set logging level request
type SetLoggingLevelResult struct {
	Success bool `json:"success"`
}

// LogMessageNotification represents a log message notification
type LogMessageNotification struct {
	Level   LoggingLevel           `json:"level"`
	Message string                 `json:"message"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// NewSetLoggingLevelRequest creates a new set logging level request
func NewSetLoggingLevelRequest(level LoggingLevel) Params {
	return SetLoggingLevelRequest{
		Level: level,
	}
}

// NewSetLoggingLevelResponse creates a new set logging level response
func NewSetLoggingLevelResponse(success bool) Result {
	return SetLoggingLevelResult{
		Success: success,
	}
}

// NewLogMessageNotification creates a new log message notification
func NewLogMessageNotification(level LoggingLevel, message string, meta map[string]interface{}) Params {
	params := map[string]interface{}{
		"level":   level,
		"message": message,
	}
	if meta != nil {
		params["meta"] = meta
	}
	return params
}
