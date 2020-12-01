package binding

import (
	"encoding/xml"
	"github.com/labstack/echo/v4"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(obj interface{}, c echo.Context) error {
	decoder := xml.NewDecoder(c.Request().Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
