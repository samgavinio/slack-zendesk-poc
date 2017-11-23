package main

import (
	"html/template"
	"io"
	"net/http"
	"github.com/labstack/echo"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/zendesk/slack-poc/controllers"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}


func main() {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	handler := &controllers.Controller{Echo:e}

	e.POST("/slack/events", handler.SlackEvent).Name = "slack.events"
	e.POST("/slack/oauth/initiate", handler.InitiateOAuth).Name = "slack.oauth.initiate"
	e.GET("/slack/oauth/save", handler.SaveOAuth).Name = "slack.oauth.save"
	e.POST("/slack/oauth/is-configured", handler.IsConfigured).Name = "slack.oauth.is-configured"

	e.GET("/zendesk/manifest.json", handler.Manifest).Name = "zendesk.manifest"
	e.POST("/zendesk/admin-ui", handler.SetupForm).Name = "zendesk.setup"

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Zendesk-Slack POC")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
