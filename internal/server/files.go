package server

import (
	"bytes"
	"encoding/json"
	"io"
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

		log.Debug().Msgf("receive %d bytes from file thread", len(msg))
		log.Debug().Msgf("message  type = %d", mt)

		// start sending files
		s.sendFile(msg)
	}
}

func (s *Server) sendFile(data []byte) {
	fm := &fileMessage{}
	if err := json.Unmarshal(data, fm); err != nil {
		log.Error().Err(err).Msg("failure to read file meta")
		return
	}
	r := bytes.NewReader(fm.Data)
	flag := false
	for !flag {
		buff := make([]byte, 512)
		n, err := r.Read(buff)
		if err != nil {
			if err == io.EOF {
				log.Debug().Msg("stop sending messages")
				flag = true
			} else {
				log.Error().Err(err).Msg("failure to send files")
			}
		}
		// log.Debug().Msgf()
		chunk := &model.Packet{
			FileName:  fm.Name,
			IsFile:    true,
			IsMessage: false,
			Data:      buff[:n],
			Close:     flag,
		}
		msg, err := json.Marshal(chunk)
		if err != nil {
			log.Error().Err(err).Msgf("failure to create chunk")
		} else {
			log.Debug().Msgf("send %d bytes of chunk", len(msg))
			s.manager.SendData(msg)
		}
	}
}

type fileMessage struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

func (s *Server) writeFiles(c *websocket.Conn) {
	files := make(map[string]*bytes.Buffer)
	cnt := 0
	for {
		// receive row binary file data
		msg := s.manager.ReceiveFile()
		log.Debug().Msg("received file from manager")

		// get all info about packet
		chunk := &model.Packet{}
		if err := json.Unmarshal(msg, chunk); err != nil {
			log.Error().Err(err).Msg("fail to unmarshal packet")
			continue
		}

		// check file in map
		log.Debug().Msgf("get chunk of file %s", chunk.FileName)
		cnt += len(chunk.Data)
		if buff, ok := files[chunk.FileName]; !ok {
			files[chunk.FileName] = &bytes.Buffer{}
			files[chunk.FileName].Write(chunk.Data)
		} else {
			log.Debug().Msgf("write to buffer file %s %b", chunk.FileName, chunk.Close)
			// all daa for this file has been sent
			if chunk.Close {
				log.Debug().Msgf("LLLLLEEEEEn %d", cnt)
				saveFile(c, buff, chunk.FileName)
				delete(files, chunk.FileName)
			} else {
				buff.Write(chunk.Data)
			}
		}
	}
}

func saveFile(c *websocket.Conn, buff *bytes.Buffer, fileName string) {
	log.Debug().Msgf("collect all file %s data %d bytes", fileName, buff.Len())
	// save file into local storage
	if err := os.WriteFile(fileStoragePath+"/"+fileName, buff.Bytes(), 0666); err != nil {
		log.Error().Err(err).Msgf("failure to create file %s", fileName)
	} else {
		if err := c.WriteMessage(websocket.TextMessage, []byte(fileName)); err != nil {
			log.Error().Err(err).Msgf("failure to write message to websocket")
		}
	}
}
