// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"github.com/labstack/echo"
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
	c.Request().ParseMultipartForm(32 << 10) // 32 MB
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
