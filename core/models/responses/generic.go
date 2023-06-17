package responses

type Arrayed struct {
	Data any `json:"data"`
}

var Acknowledged = struct {
	Ack bool `json:"ack"`
}{Ack: true}
