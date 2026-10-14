package ripple

import "testing"

type validateStruct struct {
	Name  string `validate:"required"`
	Count int    `validate:"omitempty,min=1,max=100"`
}

func TestNewValidatorDefaultTag(t *testing.T) {
	v := NewValidator("")

	// validate tag 生效：Name 为空应报错
	if err := v.Validate(validateStruct{Name: ""}); err == nil {
		t.Error("expected error for empty Name with validate:required")
	}
	// Count 越界应报错
	if err := v.Validate(validateStruct{Name: "x", Count: 200}); err == nil {
		t.Error("expected error for Count=200 with validate:max=100")
	}
	// 合法值应通过
	if err := v.Validate(validateStruct{Name: "x", Count: 50}); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestNewValidatorCustomTag(t *testing.T) {
	v := NewValidator("binding")

	type s struct {
		Name string `binding:"required"`
	}
	if err := v.Validate(s{Name: ""}); err == nil {
		t.Error("expected error for empty Name with binding:required")
	}
	if err := v.Validate(s{Name: "x"}); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// 使用 binding tag 后，validate tag 应被忽略（不再校验）
	type s2 struct {
		Name string `validate:"required"`
	}
	if err := v.Validate(s2{Name: ""}); err != nil {
		t.Errorf("validate tag should be ignored when tag name is binding, got %v", err)
	}
}

func TestValidatorNilSafety(t *testing.T) {
	var v *Validator
	if err := v.Validate(validateStruct{Name: ""}); err != nil {
		t.Errorf("nil Validator should skip validation, got %v", err)
	}
}
