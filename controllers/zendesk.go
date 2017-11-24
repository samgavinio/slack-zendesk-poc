package controllers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"

	"gopkg.in/go-playground/validator.v9"
	"github.com/labstack/echo"
	"github.com/nlopes/slack"

	"github.com/zendesk/slack-poc/operation"
	"github.com/zendesk/slack-poc/models"
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

	channelbackRequest struct {
		Message string `form:"message"`
		ParentId string `form:"parent_id"`
		RecipientId string `form:"recipient_id"`
		RequestUniqueIdentifier string `form:"request_unique_identifier"`
		Metadata string `form:"metadata"`
	}

	metadata struct {
		Subdomain string `json:"subdomain"`
	}

	channelbackResponse struct {
		ExternalId string `json:"message"`
		AllowChannelback bool `json:"allow_channel_back"`
	}
)

func (handler *Controller) Manifest(c echo.Context) (err error) {
	urls := urls{
		AdminUi: handler.Echo.Reverse("zendesk.setup"),
		ChannelbackUrl: handler.Echo.Reverse("zendesk.channel-back"),
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

func (handler *Controller) ChannelBack (c echo.Context) (err error) {
	request := new(channelbackRequest)
	if err = c.Bind(request); err != nil {
		return err
	}

	var meta metadata
	if err = json.Unmarshal([]byte(request.Metadata), &meta); err != nil {
		return err
	}

	// todo: validation of the channelback request/metadata here
	var integration models.Integration
	query := operation.DB.Where(&models.Integration{
		ZendeskSubdomain: meta.Subdomain,
	})
	if query.First(&integration).RecordNotFound() {
		return c.JSON(http.StatusBadRequest, nil)
	} else {
		channel := strings.Split(request.ParentId, "-")[0]
		api := slack.New(integration.SlackToken)
		_, timestamp, err := api.PostMessage(channel, request.Message, slack.PostMessageParameters{})
		if err != nil {
			return err
		}

		response := channelbackResponse{
			ExternalId: fmt.Sprintf("%s-%s", channel, timestamp),
			AllowChannelback: true,
		}

		return c.JSON(http.StatusOK, response)

	}
}
