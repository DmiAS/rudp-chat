package server

const (
	apiPath       = "/api/v1"
	websocketPath = "/ws"
	registerPath  = "/register"
	connectPath   = "/connect"
	clientName    = "name"
)

func (s *Server) initRoutes() {
	s.app.Use(upgradeWesocket)

	s.app.Post(apiPath+registerPath+"/:"+clientName, s.register)
	s.app.Post(apiPath+connectPath+"/:"+clientName, s.connect)
}
