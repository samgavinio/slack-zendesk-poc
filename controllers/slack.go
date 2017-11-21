package controllers

import (
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"github.com/labstack/echo"
	"github.com/zendesk/slack-poc/config"
)

type (
	verificationRequest struct {
		Type string `json:"type" validate:"required,eq=url_verification"`
		Token string `json:"token" validate:"required"`
		Challenge string `json:"challenge" validate:"required"`
	}
	verificationResponse struct {
		Challenge string `json:"challenge"`
	}
)

func (handler *Controller) SlackEventOnboard(c echo.Context) (err error) {
	request := new(verificationRequest)
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
	response :=  verificationResponse{Challenge: request.Challenge}

	return c.JSON(http.StatusOK, response)
}
