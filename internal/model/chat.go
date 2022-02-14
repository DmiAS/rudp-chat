package model

type Packet struct {
	FileName  string `json:"file_name"`
	IsFile    bool   `json:"is_file"`
	IsMessage bool   `json:"is_message"`
	Data      []byte `json:"data"`
	Close     bool   `json:"close"`
}
