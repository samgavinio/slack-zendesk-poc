package controllers

import (
	"fmt"
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"github.com/nlopes/slack"
	"github.com/labstack/echo"
	"github.com/benmanns/goworker"
	"github.com/zendesk/slack-poc/config"
)

type (
	payload struct {
		Type string `json:"type" validate:"required"`
		Token string `json:"token" validate:"required"`
		Challenge string `json:"challenge"`
		Event event
	}
	event struct {
		Type string
		Text string
	}
	verificationResponse struct {
		Challenge string `json:"challenge"`
	}
	saveAuth struct {
		Code string `query:"code" validate:"required"`
		State string `query:"state" validate:"required"`
	}
	initiateAuth struct {
		Workspace string `form:"workspace" validate:"required"`
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
	// cfg.SlackVerificationToken != request.Token
	if false {
		return echo.NewHTTPError(http.StatusBadRequest, "Verification token does not match.");
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
				Args: []interface{}{request.Event},
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
	token, _, err := slack.GetOAuthToken(cfg.SlackAppClientId, cfg.SlackAppClientSecret, request.Code,  returnUrl, false)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "admin", nil)
}
