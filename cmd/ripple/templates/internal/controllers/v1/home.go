package v1

import (
	"github.com/bmbstack/ripple"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HomeController struct {
	Index  echo.HandlerFunc `controller:"GET /"`
	Html   echo.HandlerFunc `controller:"GET html"`
	String echo.HandlerFunc `controller:"GET string"`
}

func (this HomeController) Register() {
	ripple.Default().RegisterController(&HomeController{})
}

func (this HomeController) Path() string {
	return "/v1/"
}

func (this HomeController) ActionIndex(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "home/index.html", map[string]interface{}{
		"title": "Hello world",
	})
}

func (this HomeController) ActionHtml(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "home/html.html", map[string]interface{}{
		"title": "Hello, this is a html template",
	})
}

func (this HomeController) ActionString(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, this is a string")
}
