package models

// RemoteLog ...
type RemoteLog struct {
	LogLevel string                 `json:"log_level" track:"log_level"`
	Message  string                 `json:"message" track:"message"`
	Data     map[string]interface{} `json:"data" track:"data"`
}

// Event ...
func (a RemoteLog) Event() string {
	return "remote_log"
}

// Model ...
func (a RemoteLog) Model() interface{} {
	return a
}

// UserID ...
func (a RemoteLog) UserID() string {
	if stepID, ok := a.Data["step_id"].(string); ok {
		return stepID
	}
	return ""
}
