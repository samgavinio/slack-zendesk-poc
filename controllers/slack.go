package controllers

import (
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"github.com/labstack/echo"
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
	}
	verificationResponse struct {
		Challenge string `json:"challenge"`
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
	}

	return c.JSON(http.StatusOK, nil)
}
