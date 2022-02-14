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
	ch := make(chan []byte, 1)
	for {
		select {
		case msg := <-ch:
			if err := s.handleChat(msg); err != nil {
				log.Debug().Err(err).Msg("something in socket")
				if err := c.WriteJSON(model.ErrResponse{Msg: err.Error()}); err != nil {
					log.Debug().Err(err).Msg("failure to send json via socket")
				}
			}
		case name := <-s.cli.GetSignalChan():
			log.Debug().Msg("received signal")
			if err := s.manager.StartServer(); err != nil {
				log.Error().Err(err).Msg("failure to start server")
				continue
			}
			if err := c.WriteMessage(websocket.TextMessage, name); err != nil {
				log.Error().Err(err).Msg("failure to send ping message to server")
				continue
			}
		}
	}
}

func (s *Server) readFromConn(c *websocket.Conn, ch chan []byte) {
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Debug().Err(err).Msg("read from socket error")
			continue
		}
		log.Debug().Msgf("receive %d bytes from chat thread", len(msg))
		log.Debug().Msgf("message  type = %d", mt)
		ch <- msg
	}
}

func (s *Server) handleChat(data []byte) error {
	event := &chatEvent{}
	if err := json.Unmarshal(data, event); err != nil {
		return fmt.Errorf("failure to unmarshal chat action")
	}

	log.Debug().Msgf("value = %+v", event)
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
