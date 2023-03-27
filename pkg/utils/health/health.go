package health

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/m1kx/go-vtr-backend/pkg/utils/structs"
)

var alive bool = true
var last_words string = ""

func get_message() structs.HealthResponse {
	status := ""
	if alive {
		status = "alive"
	} else {
		status = "dead"
	}
	res := structs.HealthResponse{
		Status:     status,
		Last_Words: last_words,
	}
	return res
}

func Health(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, get_message())
}

func Dead(cause string) {
	alive = false
	last_words = fmt.Sprintf("%s%s |||", last_words, cause)
}
