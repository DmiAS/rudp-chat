package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"

	"chat/internal/model"
)

const (
	fileStoragePath = "/Users/d.antsibor/university/network/course/chat/files"
)

func (s *Server) readFiles(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		// receive file from gui and send it to opponent
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Debug().Err(err).Msg("read from socket error")
			break
		}
		log.Debug().Msgf("receive %d bytes from chat thread", len(msg))
		log.Debug().Msgf("message  type = %d", mt)

		// pack file in packet and send it
		packet := &model.Packet{IsFile: true, Data: msg}
		data, err := json.Marshal(packet)
		if err != nil {
			log.Error().Err(err).Msg("failure to send file")
		} else {
			s.manager.SendData(data)
		}
	}
}

func (s *Server) writeFiles(c *websocket.Conn) {
	fileCnt := 0
	for {
		// receive row binary file data
		msg := s.manager.ReceiveFile()

		// save file into local storage
		fileName := fmt.Sprintf("%s/file_%d", fileStoragePath, fileCnt)
		if err := os.WriteFile(fileName, msg, 0666); err != nil {
			log.Error().Err(err).Msg("failure to create file")
		} else {
			fileCnt++
			// send path to file to gui
			if err := c.WriteMessage(websocket.TextMessage, []byte(fileName)); err != nil {
				log.Error().Err(err).Msgf("failure to write message to websocket")
			}
		}
	}
}
