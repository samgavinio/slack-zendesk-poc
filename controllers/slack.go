package controllers

import (
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
	"github.com/nlopes/slack"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/benmanns/goworker"

	"github.com/zendesk/slack-poc/config"
	"github.com/zendesk/slack-poc/operation"
	"github.com/zendesk/slack-poc/models"
)

type (
	payload struct {
		Type string `json:"type" validate:"required"`
		Token string `json:"token" validate:"required"`
		Challenge string `json:"challenge"`
		TeamId string `json:"team_id"`
		Event event `json:"event"`
	}
	event struct {
		Type string `json:"type"`
		Text string `json:"text"`
		EventTs string `json:"event_ts"`
		User string `json:"user"`
	}
	verificationResponse struct {
		Challenge string `json:"challenge"`
	}
	saveAuth struct {
		Code string `query:"code" validate:"required"`
		State string `query:"state" validate:"required"`
	}
	initiateAuth struct {
		Subdomain string `form:"subdomain" validate:"required"`
		ReturnUrl string `form:"return_url" validate:"required"`
		Workspace string `form:"workspace" validate:"required"`
		ChannelsToken string `form:"token" validate:"required"`
		ChannelsPushClientId string `form:"push_client_id" validate:"required"`
	}
)

func (handler *Controller) SlackEvent(c echo.Context) (err error) {
	request := new(payload)
	if err = c.Bind(request); err != nil {
		return err
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error());
	}

	cfg := config.GetConfig()
	if cfg.SlackVerificationToken != request.Token {
		return echo.NewHTTPError(http.StatusBadRequest, "Verification token is invalid.");
	}

	switch request.Type {
	case "url_verification":
		return handler.handleUrlVerification(request, c)
	case "event_callback":
		return handler.handleEventCallback(request, c)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Request type is not supported.");
	}
}

func (handler *Controller) handleUrlVerification(request *payload, c echo.Context) (err error) {
	return c.JSON(http.StatusOK, verificationResponse{Challenge: request.Challenge})
}

func (handler *Controller) handleEventCallback(request *payload, c echo.Context) (err error) {
	if request.Event.Type == "message" {
		goworker.Enqueue(&goworker.Job{
			Queue: "zendesk",
			Payload: goworker.Payload{
				Class: "ProcessSlackMessage",
				Args: []interface{}{request},
			},
		})
	}

	return c.JSON(http.StatusOK, nil)
}

func (handler *Controller) InitiateOAuth (c echo.Context) (err error) {
	request := new(initiateAuth)
	if err = c.Bind(request); err != nil {
		return err
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error());
	}
		
	cfg := config.GetConfig()
	returnUrl := c.Scheme() + "://" + c.Request().Host + handler.Echo.Reverse("slack.oauth.save")
	redirectTo := fmt.Sprintf("https://slack.com/oauth/authorize?client_id=%s&client_secret=%s&scope=channels:read&team=%s&state=state&redirect_url=%s",
		cfg.SlackAppClientId,
		cfg.SlackAppClientSecret,
		request.Workspace,
		returnUrl,
	)

	sess, _ := session.Get("session", c)
	sess.Values["zendesk_subdomain"] = request.Subdomain
	sess.Values["slack_workspace"] = request.Workspace
	sess.Values["channels_token"] = request.ChannelsToken
	sess.Values["channels_push_client_id"] = request.ChannelsPushClientId
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
}

func (handler *Controller) SaveOAuth (c echo.Context) (err error) {
	request := new(saveAuth)
	if err = c.Bind(request); err != nil {
		return err
	}

	// todo: validate state nonce
	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error());
	}
	
	cfg := config.GetConfig()
	returnUrl := c.Scheme() + "://" + c.Request().Host + handler.Echo.Reverse("slack.oauth.save")
	response, err := slack.GetOAuthResponse(cfg.SlackAppClientId, cfg.SlackAppClientSecret, request.Code,  returnUrl, false)
	if err != nil {
		return err
	}

	sess, _ := session.Get("session", c)
	integration := &models.Integration{
		SlackToken: response.AccessToken,
		SlackTeamId: response.TeamID,
		SlackWorkspace: sess.Values["slack_workspace"].(string),
		ZendeskSubdomain: sess.Values["zendesk_subdomain"].(string),
		ZendeskToken: sess.Values["channels_token"].(string),
		ZendeskInstancePushId: sess.Values["channels_push_client_id"].(string),
	}
	if err = operation.DB.Create(&integration).Error; err != nil {
		return err
	}

	return c.Render(http.StatusOK, "success", nil)
}

func (handler *Controller) IsConfigured (c echo.Context) (err error) {
	request := new(initiateAuth)
	if err = c.Bind(request); err != nil {
		return err
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error());
	}

	var integration models.Integration
	query := operation.DB.Where(&models.Integration{
		ZendeskSubdomain: request.Subdomain,
		SlackWorkspace: request.Workspace,
	})
	if query.First(&integration).RecordNotFound() {
		return c.String(http.StatusBadRequest, "Integration not found")
	}

	return c.JSON(http.StatusOK, "Integration found.")
}
