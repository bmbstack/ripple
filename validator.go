package ripple

import (
	"github.com/go-playground/validator/v10"
)

// Validator 是 ripple 内置的 echo.Validator 实现，基于 go-playground/validator。
// 用于支持 controller 中的 ctx.Validate(req) 调用，校验 struct tag。
type Validator struct {
	validate *validator.Validate
}

// NewValidator 创建一个 Validator 实例。
// tag 为 struct 上使用的校验 tag 名（如 "validate"、"binding"），传空时默认使用 "validate"。
func NewValidator(tag string) *Validator {
	v := validator.New()
	if tag == "" {
		tag = "validate"
	}
	v.SetTagName(tag)
	return &Validator{validate: v}
}

// Validate 实现 echo.Validator 接口。
func (v *Validator) Validate(i interface{}) error {
	if v == nil || v.validate == nil {
		return nil
	}
	return v.validate.Struct(i)
}
