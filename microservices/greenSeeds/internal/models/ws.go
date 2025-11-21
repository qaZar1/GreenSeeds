package models

type WSMessage struct {
	Type     string     `json:"type"`
	Status   *string    `json:"status,omitempty"`
	Params   *Params    `json:"params,omitempty"`
	Payload  *PayloadWs `json:"payload,omitempty"`
	Error    *string    `json:"error,omitempty"`
	Solution *string    `json:"solution,omitempty"`
	Id       *int64     `json:"id,omitempty"`
}

type Params struct {
	Shift          int    `json:"shift"`
	Number         int    `json:"number"`
	Receipt        int    `json:"receipt"`
	Turn           int    `json:"turn"`
	RequiredAmount int    `json:"required_amount"`
	Bunker         int    `json:"bunker"`
	Gcode          string `json:"gcode"`
	ExtraMode      bool   `json:"extraMode"`
	Seed           string `json:"seed"`
}

type PayloadWs struct {
	Control *bool   `json:"control"`
	Reason  *string `json:"reason"`
	Photo   *[]byte `json:"photo"`
}
