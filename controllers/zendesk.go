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
		PushClientId string `json:"push_client_id"`
		Urls urls `json:"urls"`
	}

	urls struct {
		AdminUi string `json:"admin_ui"`
		ChannelbackUrl string `json:"channelback_url"`
	}

	setupRequest struct {
		Subdomain string `form:"subdomain" validate:"required"`
		InstancePushId string `form:"instance_push_id" validate:"required"`
		ZendeskAccessToken string `form:"zendesk_access_token" validate:"required"`
		ReturnUrl string `form:"return_url" validate:"required"`
	}
)

func (handler *Controller) Manifest(c echo.Context) (err error) {
	urls := urls{
		AdminUi: handler.Echo.Reverse("zendesk.setup"),
		ChannelbackUrl: "/zendesk/channel-back",
	}

	manifest := manifest{
		Name: "Zendesk-Slack POC",
		Id: "zendesk_slack_poc",
		Author: "samgavinio@gmail.com",
		Version: "0.0.0",
		PushClientId: "shopify-integration",
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
