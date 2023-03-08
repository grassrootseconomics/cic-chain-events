package main

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/labstack/echo/v4"
)

func initApiServer() *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	if ko.Bool("metrics.go_process") {
		server.GET("/metrics", func(c echo.Context) error {
			metrics.WritePrometheus(c.Response(), true)
			return nil
		})
	}

	return server
}
