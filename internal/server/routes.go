package server

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

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
	// serve html for gui
	s.app.Static("/", "/Users/d.antsibor/university/network/course/chat/static/build")
	api := s.app.Group(apiPath).Use(cors.New())
	{
		api.Post(registerPath+"/:"+clientName, s.register)
		api.Post(connectPath+"/:"+clientName, s.connect)
	}
	chat := s.app.Group(websocketPath + chatPath).Use(upgradeWesocket)
	{
		chat.Get("/thread", websocket.New(s.chatThread))
		chat.Get(messagePath, websocket.New(s.workWithMessages))
		chat.Get(filesPath, websocket.New(s.workWithFiles))
	}
}
