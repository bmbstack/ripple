package v1

import (
	"github.com/bmbstack/ripple"
	"github.com/labstack/echo"
	"net/http"
)

type HomeController struct {
	Index  echo.HandlerFunc `controller:"GET /"`
	Html   echo.HandlerFunc `controller:"GET html"`
	String echo.HandlerFunc `controller:"GET string"`
}

func init() {
	ripple.RegisterController(&HomeController{})
}

func (this HomeController) Path() string {
	return "/"
}

func (this HomeController) ActionIndex(ctx echo.Context) error {
	ctx.Render(http.StatusOK, "home/index.html", map[string]interface{}{
		"title": "Hello world",
	})

	return nil
}

func (this HomeController) ActionHtml(ctx echo.Context) error {
	ctx.Render(http.StatusOK, "home/html.html", map[string]interface{}{
		"title": "Hello, this is a html template",
	})

	return nil
}

func (this HomeController) ActionString(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, this is a string")

}
