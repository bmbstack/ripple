package v1

import (
	"net/http"

	. "github.com/bmbstack/ripple/fixtures/forum/helper"
	"github.com/bmbstack/ripple/fixtures/forum/models/one"
	"github.com/bmbstack/ripple/fixtures/forum/models/two"
	"github.com/labstack/echo"
)

type CityController struct {
	Group *echo.Group
	BaseController
}

func (this CityController) Setup() {
	this.Group.GET("/citys", this.ActionIndex)
	this.Group.GET("/citys/two", this.ActionTwoIndex)
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

	bmbCity := &one.BmbCity{}
	result := bmbCity.FindCityListByPid(ctx, params.Pid)
	return ctx.JSON(http.StatusOK, SuccessJSON(result))
}

func (this CityController) ActionTwoIndex(ctx echo.Context) error {
	params := &RequestCity{}
	err := ctx.Bind(params)
	if err != nil {
		return ctx.JSON(http.StatusOK, ErrorJSON(ErrorMsgParamsValidateFailed, ErrorCodeParamsValidateFailed))
	}

	bmbCity := &two.BmbCity{}
	result := bmbCity.FindCityListByPid(ctx, params.Pid)
	return ctx.JSON(http.StatusOK, SuccessJSON(result))
}
