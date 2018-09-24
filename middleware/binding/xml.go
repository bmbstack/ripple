// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/xml"
	"github.com/labstack/echo"
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
