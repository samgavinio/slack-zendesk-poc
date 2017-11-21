package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/zendesk/slack-poc/controllers"
)

func main() {
	e := echo.New()
	handler := &controllers.Controller{}

	e.POST("/slack/events", handler.SlackEvent)
	e.GET("/zendesk/manifest", handler.Manifest)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Zendesk-Slack POC")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
