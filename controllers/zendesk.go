package controllers

import (
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"github.com/labstack/echo"
)

type (
	manifest struct {
		Name string `json:"name"`
		Id string `json:"id"`
		Author string `json:"author"`
		Version string `json:"version"`
		Urls urls `json:"urls"`
	}

	urls struct {
		AdminUi string `json:"admin_ui"`
		PullUrl string `json:"pull_url"`
		ChannelbackUrl string `json:"channelback_url"`
	}

	setupRequest struct {
		Subdomain string `form:"subdomain" validate:"required"`
		ReturnUrl string `form:"return_url" validate:"required"`
	}
)

func (handler *Controller) Manifest(c echo.Context) (err error) {
	urls := urls{
		AdminUi: handler.Echo.Reverse("zendesk.setup"),
		PullUrl: "/zendesk/pull",
		ChannelbackUrl: "/zendesk/channel-back",
	}

	manifest := manifest{
		Name: "Zendesk-Slack POC",
		Id: "zendesk_slack_poc",
		Author: "samgavinio@gmail.com",
		Version: "0.0.0",
		Urls: urls,
	}

	return c.JSON(http.StatusOK, manifest)
}

func (handler *Controller) SetupForm (c echo.Context) (err error) {
	request := new(setupRequest)
	if err = c.Bind(request); err != nil {
		return err
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error());
	}

	return c.Render(http.StatusOK, "admin", request)
}
