package controllers

import (
	"net/http"
	"github.com/labstack/echo"
)

func (handler *Controller) SlackEventOnboard(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, handler.Success(nil, "Hello world."))
}
