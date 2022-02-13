package server

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"

	"chat/internal/model"
)

type chatEvent struct {
	Address string `json:"address"`
	Action  string `json:"action"`
}

const (
	startChat = "start"
	stopChat  = "stop"
)

// handle messages for connecting go clients
func (s *Server) chatThread(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Debug().Err(err).Msg("read from socket error")
			break
		}
		log.Debug().Msgf("receive %d bytes from chat thread", len(msg))
		log.Debug().Msgf("message  type = %d", mt)

		if err := s.handleChat(msg); err != nil {
			if err := c.WriteJSON(model.ErrResponse{Msg: err.Error()}); err != nil {
				log.Debug().Err(err).Msg("failure to send json via socket")
			}
		}
	}
}

func (s *Server) handleChat(data []byte) error {
	event := &chatEvent{}
	if err := json.Unmarshal(data, event); err != nil {
		return fmt.Errorf("failure to unmarshal chat action")
	}

	switch action := event.Action; action {
	case startChat:
		if err := s.manager.Connect(event.Address); err != nil {
			return fmt.Errorf("failure to start chat: %s", err)
		}
	case stopChat:
		s.manager.Disconnect()
	default:
		return fmt.Errorf("unknown event: %s", action)
	}

	return nil
}

func (s *Server) workWithMessages(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Debug().Err(err).Msg("read from socket error")
			break
		}
		log.Debug().Msgf("receive %d bytes from chat thread", len(msg))
		log.Debug().Msgf("message  type = %d", mt)

	}
}

func (s *Server) workWithFiles(c *websocket.Conn) {

}
