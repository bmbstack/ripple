package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/bmbstack/ripple/cmd/ripple/templates/internal/ecode"
)

type CityController struct {
	Group *echo.Group
	BaseController
}

func (this CityController) Setup() {
	this.Group.GET("/list", this.ActionIndex)
}

type RequestCity struct {
	Pid int64 `form:"pid" json:"pid"`
}

func (this CityController) ActionIndex(ctx echo.Context) error {
	params := &RequestCity{}
	err := ctx.Bind(params)
	if err != nil {
		return ecode.Error(ctx, ecode.ParamError)
	}

	var result []map[string]interface{}
	result = append(result, map[string]interface{}{
		"id":   1,
		"name": "北京市",
	})
	result = append(result, map[string]interface{}{
		"id":   2,
		"name": "上海市",
	})
	return ecode.OK(ctx, result)
}
