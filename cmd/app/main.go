package main

import (
	"net/http"
	"os"

	"github.com/joshjms/firefly-executor/pkg/controller"
	"github.com/joshjms/firefly-executor/pkg/handlers/submit"
	"github.com/labstack/echo/v4"
)

func main() {
	controller.NewBoxController(1000)

	os.Mkdir("mounts", 0755)
	os.Mkdir("metadata", 0755)

	e := echo.New()
	e.POST("/submit", Submit)

	e.Logger.Fatal(e.Start(":8080"))
}

func Submit(c echo.Context) error {
	var j submit.Job
	if err := c.Bind(&j); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	resp := submit.Submit(&j)

	return c.JSON(http.StatusOK, resp)
}
