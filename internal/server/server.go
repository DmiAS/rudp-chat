package server

import (
	"net/http"

	"github.com/DmiAS/rendezvous/pkg/proto/client"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"chat/internal/chat"
)

type Server struct {
	hostAddress string
	app         *fiber.App
	cli         *client.Client
	manager     *chat.Manager
}

const (
	serverAddress = "localhost"
)

func NewServer(cli *client.Client) *Server {
	srv := &Server{
		cli:         cli,
		hostAddress: serverAddress,
		manager:     chat.NewManager(cli.GetConnection()),
	}
	srv.app = fiber.New(
		fiber.Config{
			DisableStartupMessage: false,
			CaseSensitive:         true,
			StrictRouting:         true,
		},
	)
	srv.app.Use(recover.New(), logger.New())
	srv.initRoutes()
	return srv
}

func (s *Server) Run() {
	log.Info().Msgf("server started listen on address: %s", s.hostAddress)
	if err := s.app.Listen(s.hostAddress); err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Err(err).
			Msgf("the HTTP rest stopped with unknown error")
	}
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
