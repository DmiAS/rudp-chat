package server

import (
	"encoding/json"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"

	"chat/internal/model"
)

func (s *Server) readMessages(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		// receive message from gui and send it to another user
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Debug().Err(err).Msg("read from socket error")
			break
		}
		log.Debug().Msgf("receive %d bytes from read message thread", len(msg))
		log.Debug().Msgf("message  type = %d", mt)

		// pack message in packet and send it
		packet := &model.Packet{IsMessage: true, Data: msg}
		data, err := json.Marshal(packet)
		if err != nil {
			log.Error().Err(err).Msg("failure to send message")
		} else {
			log.Debug().Msgf("prepare to send message %+v", packet)
			s.manager.SendData(data)
		}
	}
}

func (s *Server) writeMessages(c *websocket.Conn) {
	for {
		msg := s.manager.ReceiveMessage()
		log.Debug().Msgf("received message:%s", string(msg))
		if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Debug().Err(err).Msg("failure to send message to gui")
		}
	}
}
