package binding

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(obj interface{}, c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
