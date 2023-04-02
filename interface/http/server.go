package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func StartHttpServer() (err error) {
	server := echo.New()
	server.HideBanner = true

	server.GET("/", func(ctx echo.Context) error {
		msg := fmt.Sprintf("Preparation Server Up ,,,,, (%s)", time.Now().Format(time.RFC3339))
		return ctx.String(http.StatusOK, msg)
	})

	return server.Start(fmt.Sprintf("%s:%d", "0.0.0.0", 4000))
}
