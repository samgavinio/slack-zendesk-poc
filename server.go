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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
