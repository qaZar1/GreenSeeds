package models

type WSMessage struct {
	Type    string      `json:"type"`
	TaskID  string      `json:"task_id,omitempty"`
	Status  string      `json:"status,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}
