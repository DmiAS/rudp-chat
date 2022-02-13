package model

type Packet struct {
	IsFile    bool   `json:"is_file"`
	IsMessage bool   `json:"is_message"`
	Data      []byte `json:"data"`
}
