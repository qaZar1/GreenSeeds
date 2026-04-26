package models

type WSRequest struct {
	Type     WSMessageType `json:"type"`
	Id       *int64        `json:"id"`
	Params   *Params       `json:"params"`
	Token    *string       `json:"token"`
	Solution *string       `json:"solution"`
}

type WSResponse struct {
	Type WSMessageType `json:"type"`
	Status  string     `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Iteration int
	Data    interface{} `json:"data,omitempty"`
	Actions *[]string   `json:"actions,omitempty"`
}

type Params struct {
	Shift          int    `json:"shift"`
	Number         int    `json:"number"`
	Receipt        int `json:"receipt"`
	Turn           int    `json:"turn"`
	RequiredAmount int    `json:"required_amount"`
	Bunker         int    `json:"bunker"`
	Gcode          string `json:"gcode"`
	ExtraMode      bool   `json:"extraMode"`
	Seed           string `json:"seed"`
}