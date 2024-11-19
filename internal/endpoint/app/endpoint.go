package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rs/zerolog"

	config "Audio2TextService/internal/config/app"
)

type Handler interface {
	HandleRequest(c echo.Context) error
}

type Endpoint struct {
	handler Handler
	Logger  *zerolog.Logger
}

func New(handler Handler, logger *zerolog.Logger) *Endpoint {
	return &Endpoint{handler: handler, Logger: logger}
}

func (e *Endpoint) Start() error {

	// Cоздаем новый Echo-сервер и привязываем его к порту
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.Gzip())

	server.POST("/", e.handler.HandleRequest)
	return server.Start(config.ADDR)
}
