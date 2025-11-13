package models

type WSMessage struct {
	Type    string     `json:"type"`
	Status  *string    `json:"status,omitempty"`
	Params  *Params    `json:"params,omitempty"`
	Payload *PayloadWs `json:"payload,omitempty"`
	Error   *string    `json:"error,omitempty"`
}

type Params struct {
	Shift           int    `json:"shift"`
	Number          int    `json:"number"`
	Receipt         int    `json:"receipt"`
	Amount          int    `json:"amount"`
	CompletedAmount int    `json:"completed_amount"`
	RequiredAmount  int    `json:"required_amount"`
	Bunker          int    `json:"bunker"`
	Gcode           string `json:"gcode"`
	ExtraMode       bool   `json:"extraMode"`
}

type PayloadWs struct {
	Control bool `json:"control"`
}
