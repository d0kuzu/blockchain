package network

import (
	"encoding/json"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (m *Message) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func NewMessage(msgType string, data interface{}) Message {
	jsonData, _ := json.Marshal(data)
	return Message{
		Type: msgType,
		Data: jsonData,
	}
}
