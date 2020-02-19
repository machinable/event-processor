package webhook

// HookResult contains relevant information regarding the http response
type HookResult struct {
	WebHookID    string  `json:"webhook_id"`
	ProjectID    string  `json:"project_id"`
	StatusCode   int     `json:"status_code"`
	ResponseTime float64 `json:"response_time"`
	ErrorMessage string  `json:"error_message"`
}

// WebHook wraps the webhook object of the event
type WebHook struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"project_id"`
	Label     string   `json:"label"`
	IsEnabled bool     `json:"is_enabled"`
	Entity    string   `json:"entity"`
	EntityID  string   `json:"entity_id"`
	HookEvent string   `json:"event"`
	Headers   []Header `json:"headers"`
	HookURL   string   `json:"hook_url"`
}

// Header is a http header for a web hook
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HookEvent describes a single web hook event
type HookEvent struct {
	Hook      *WebHook               `json:"hook"`
	EntityKey string                 `json:"entity_key"`
	Payload   map[string]interface{} `json:"payload"`
}

// Format formats the event for post
func (h *HookEvent) Format() map[string]interface{} {
	payload := map[string]interface{}{}
	payload["payload"] = h.Payload
	payload["event"] = h.Hook.HookEvent
	payload["project_Id"] = h.Hook.ProjectID
	payload["entity_type"] = h.Hook.Entity
	payload["entity_key"] = h.EntityKey

	return payload
}
