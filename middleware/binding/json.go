// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/json"
	"github.com/labstack/echo"
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
