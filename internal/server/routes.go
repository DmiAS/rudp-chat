package server

import "github.com/gofiber/websocket/v2"

const (
	// http
	apiPath      = "/api/v1"
	registerPath = "/register"
	connectPath  = "/connect"
	clientName   = "name"

	// websockets
	chatPath      = "/chat"
	websocketPath = "/ws"
	messagePath   = "/message"
	filesPath     = "/files"
)

func (s *Server) initRoutes() {
	s.app.Use(upgradeWesocket)

	s.app.Post(apiPath+registerPath+"/:"+clientName, s.register)
	s.app.Post(apiPath+connectPath+"/:"+clientName, s.connect)
	chat := s.app.Group(websocketPath + chatPath)
	{
		chat.Get("/", websocket.New(s.chatThread))
		chat.Get(messagePath, websocket.New(s.workWithMessages))
		chat.Get(filesPath, websocket.New(s.workWithFiles))
	}
}
