package controllers

import (
	"net/http"
	"github.com/labstack/echo"
)

type (
	manifest struct {
		Name string `json:"name"`
		Id string `json:"id"`
		Author string `json:"author"`
		Version string `json:"version"`
		Urls urls `json:"url"`
	}

	urls struct {
		AdminUi string `json:"admin_ui"`
	}
)

func (handler *Controller) Manifest(c echo.Context) (err error) {
	urls := urls{
		AdminUi: "/zendesk/admin_ui",
	}

	manifest := manifest{
		Name: "Zendesk-Slack POC",
		Id: "com.zendesk.slack.integration.poc",
		Author: "samgavinio@gmail.com",
		Version: "0.0.0",
		Urls: urls,
	}

	return c.JSON(http.StatusOK, manifest)
}
