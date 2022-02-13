package server

import (
	"github.com/gofiber/fiber/v2"

	"chat/internal/model"
)

func (s *Server) register(ctx *fiber.Ctx) error {
	name := ctx.Params(clientName)
	if name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrResponse{Msg: "empty name"})
	}

	if err := s.cli.Register(name); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrResponse{Msg: err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s *Server) connect(ctx *fiber.Ctx) error {
	name := ctx.Params(clientName)
	if name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrResponse{Msg: "empty target name"})
	}

	addr, err := s.cli.ConnectTo(name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrResponse{Msg: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.AddrResponse{Address: addr.String()})
}
