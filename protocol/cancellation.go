package protocol

// CancelledNotification represents a notification that a request has been cancelled
type CancelledNotification struct {
	RequestID RequestID `json:"requestId"`
	Reason    string    `json:"reason,omitempty"`
}

// NewCancelledNotification creates a new cancelled notification
func NewCancelledNotification(requestID RequestID, reason string) Params {
	params := map[string]interface{}{
		"requestId": requestID,
	}
	if reason != "" {
		params["reason"] = reason
	}
	return params
}
