package models

type WSRequest struct {
	Type     WSMessageType `json:"type"`
	Id       *int64        `json:"id"`
	Params   *Params       `json:"params"`
	Token    *string       `json:"token"`
	Solution *string       `json:"solution"`
}

type WSResponse struct {
	Type    WSMessageType `json:"type,omitempty"`
	Status  string        `json:"status,omitempty"`
	Message string        `json:"message,omitempty"`

	Error *WSError `json:"error,omitempty"`

	// planting events
	Event     string    `json:"event,omitempty"`
	Iteration int       `json:"iteration,omitempty"`
	Step      string    `json:"step,omitempty"`
	Progress  *Progress `json:"progress,omitempty"`
	Data      *Params   `json:"data,omitempty"`
}

type Progress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
	Percent int `json:"percent"`
}

type WSError struct {
	Code    string `json:"code"`
	Stage   string `json:"stage"`
	Message string `json:"message"`
}

type Params struct {
	Shift          int    `json:"shift"`
	Number         int    `json:"number"`
	Recipe         int    `json:"recipe"`
	Turn           int    `json:"turn"`
	RequiredAmount int    `json:"required_amount"`
	Bunker         int    `json:"bunker"`
	Gcode          string `json:"gcode"`
	ExtraMode      bool   `json:"extra_mode"`
	Seed           string `json:"seed"`
}
