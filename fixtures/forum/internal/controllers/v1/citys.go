package v1

import (
	"fmt"
	"github.com/bmbstack/ripple/fixtures/forum/internal/ecode"
	"github.com/bmbstack/ripple/fixtures/forum/internal/models/one"
	"github.com/bmbstack/ripple/fixtures/forum/internal/models/two"
	"github.com/labstack/echo/v4"
)

type CityController struct {
	Group *echo.Group
	BaseController
}

func (this CityController) Setup() {
	this.Group.GET("/citys/one", this.ActionIndex)
	this.Group.GET("/citys/oneid", this.ActionIndexID)
	this.Group.GET("/citys/one/save", this.ActionIndexSave)
	this.Group.GET("/citys/two", this.ActionTwoIndex)
	this.Group.GET("/citys/twoid", this.ActionTwoIndexID)
	this.Group.GET("/citys/two/save", this.ActionTwoIndexSave)
}

//================================================================================================
type RequestCity struct {
	Pid int64 `form:"pid" json:"pid"`
}

type RequestCityID struct {
	ID int64 `form:"id" json:"id"`
}

func (this CityController) ActionIndex(ctx echo.Context) error {
	params := &RequestCity{}
	err := ctx.Bind(params)
	if err != nil {
		return ecode.Error(ctx, ecode.ParamError)
	}

	bmbCity := &one.BmbCity{}
	result := bmbCity.FindCityListByPid(params.Pid)
	return ecode.OK(ctx, result)
}

func (this CityController) ActionIndexID(ctx echo.Context) error {
	params := &RequestCityID{}
	err := ctx.Bind(params)
	if err != nil {
		return ecode.Error(ctx, ecode.ParamError)
	}

	bmbCity := &one.BmbCity{}
	result := bmbCity.FindCityID(params.ID)
	return ecode.OK(ctx, result)
}

func (this CityController) ActionIndexSave(ctx echo.Context) error {
	params := &RequestCityID{}
	err := ctx.Bind(params)
	if err != nil {
		return ecode.Error(ctx, ecode.ParamError)
	}

	bmbCity := &one.BmbCity{}
	result := bmbCity.FindCityID(params.ID)
	result.Name = fmt.Sprintf("%s_up", result.Name)
	result.Save()
	return ecode.OK(ctx, result)
}

func (this CityController) ActionTwoIndex(ctx echo.Context) error {
	params := &RequestCity{}
	err := ctx.Bind(params)
	if err != nil {
		return ecode.Error(ctx, ecode.ParamError)
	}

	bmbCity := &two.BmbCity{}
	result := bmbCity.FindCityListByPid(params.Pid)
	return ecode.OK(ctx, result)
}

func (this CityController) ActionTwoIndexID(ctx echo.Context) error {
	params := &RequestCityID{}
	err := ctx.Bind(params)
	if err != nil {
		return ecode.Error(ctx, ecode.ParamError)
	}

	bmbCity := &two.BmbCity{}
	result := bmbCity.FindCityID(params.ID)
	return ecode.OK(ctx, result)
}

func (this CityController) ActionTwoIndexSave(ctx echo.Context) error {
	params := &RequestCityID{}
	err := ctx.Bind(params)
	if err != nil {
		return ecode.Error(ctx, ecode.ParamError)
	}

	bmbCity := &two.BmbCity{}
	result := bmbCity.FindCityID(params.ID)
	result.Name = fmt.Sprintf("%s_up", result.Name)
	result.Save()
	return ecode.OK(ctx, result)
}
