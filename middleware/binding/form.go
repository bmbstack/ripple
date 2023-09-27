package binding

import (
	"github.com/labstack/echo/v4"
)

type formBinding struct{}
type formPostBinding struct{}
type formMultipartBinding struct{}

func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(obj interface{}, c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	// 32 MB
	_ = c.Request().ParseMultipartForm(32 << 10)
	if err := mapForm(obj, c.Request().Form); err != nil {
		return err
	}
	return validate(obj)
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (formPostBinding) Bind(obj interface{}, c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, c.Request().PostForm); err != nil {
		return err
	}
	return validate(obj)
}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (formMultipartBinding) Bind(obj interface{}, c echo.Context) error {
	if err := c.Request().ParseMultipartForm(32 << 10); err != nil {
		return err
	}
	if err := mapForm(obj, c.Request().MultipartForm.Value); err != nil {
		return err
	}
	return validate(obj)
}
