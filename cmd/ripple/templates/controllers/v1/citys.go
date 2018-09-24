package v1

import (
	"net/http"

	. "github.com/bmbstack/ripple/cmd/ripple/templates/helper"
		"github.com/labstack/echo"
)

type CityController struct {
	Group *echo.Group
	BaseController
}

func (this CityController) Setup() {
	this.Group.GET("/citys", this.ActionIndex)
}

//================================================================================================
type RequestCity struct {
	Pid int64 `form:"pid" json:"pid"`
}

func (this CityController) ActionIndex(ctx echo.Context) error {
	params := &RequestCity{}
	err := ctx.Bind(params)
	if err != nil {
		return ctx.JSON(http.StatusOK, ErrorJSON(ErrorMsgParamsValidateFailed, ErrorCodeParamsValidateFailed))
	}

	var result []map[string]interface{}
	result = append(result, map[string]interface{}{
		"id": 1,
		"name": "北京市",
	})
	result = append(result, map[string]interface{}{
		"id": 2,
		"name": "上海市",
	})
	return ctx.JSON(http.StatusOK, SuccessJSON(result))
}
