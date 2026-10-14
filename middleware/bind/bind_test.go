package bind

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// stubValidator 是一个简单的 echo.Validator 桩实现。
type stubValidator struct {
	err error
}

func (s *stubValidator) Validate(i interface{}) error { return s.err }

type bindReq struct {
	Name string `json:"name" validate:"required"`
}

// TestDefaultBinderBindWithValidator 验证：绑定完成后会自动调用 Echo.Validator 校验，
// 校验失败时 Bind 返回校验错误。
func TestDefaultBinderBindWithValidator(t *testing.T) {
	e := echo.New()
	e.Binder = &DefaultBinder{}
	e.Validator = &stubValidator{err: errors.New("validation failed")}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"x"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var dst bindReq
	if err := e.Binder.Bind(&dst, c); err == nil {
		t.Fatal("expected validation error from Binder, got nil")
	}
}

// TestDefaultBinderBindWithoutValidator 验证：未设置 Validator 时，Bind 只做绑定不做校验，
// 不返回错误。
func TestDefaultBinderBindWithoutValidator(t *testing.T) {
	e := echo.New()
	e.Binder = &DefaultBinder{}
	// 不设置 e.Validator

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"x"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var dst bindReq
	if err := e.Binder.Bind(&dst, c); err != nil {
		t.Fatalf("expected nil error without validator, got %v", err)
	}
	if dst.Name != "x" {
		t.Errorf("expected bound name %q, got %q", "x", dst.Name)
	}
}
