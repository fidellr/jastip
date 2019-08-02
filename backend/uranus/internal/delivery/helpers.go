package delivery

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type ErrorHTTPResponse struct {
	Message string `json:"error"`
}

func HandleUncaughtHTTPError(err error, c echo.Context) {
	logrus.Error(err)
	c.JSON(http.StatusInternalServerError, ErrorHTTPResponse{Message: err.Error()})
}
