package binding

import (
	"github.com/labstack/echo/v4"
	"strings"
)

type Binding interface {
	Name() string
	Bind(i interface{}, c echo.Context) error
}

type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface{}) error
}

var Validator StructValidator = &defaultValidator{}

var (
	XML           = xmlBinding{}
	JSON          = jsonBinding{}
	Form          = formBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
)

func Default(method, contentType string) Binding {
	if method == "GET" {
		return Form
	} else {
		// POST
		switch {
		case strings.HasPrefix(contentType, echo.MIMEApplicationJSON):
			return JSON
		case strings.HasPrefix(contentType, echo.MIMEApplicationXML), strings.HasPrefix(contentType, echo.MIMETextXML):
			return XML
		case strings.HasPrefix(contentType, echo.MIMEApplicationForm), strings.HasPrefix(contentType, echo.MIMEMultipartForm):
			return Form
		default:
			return Form
		}
	}
}

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}

type Binder struct {
}

func (Binder) Bind(i interface{}, c echo.Context) (err error) {
	b := Default(c.Request().Method, c.Request().Header.Get("Content-Type"))
	err = b.Bind(i, c)
	return
}
