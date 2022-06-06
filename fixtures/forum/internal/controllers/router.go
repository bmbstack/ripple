package controllers

import (
	"context"
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/internal/controllers/v1"
	. "github.com/bmbstack/ripple/fixtures/forum/internal/helper"
	"github.com/bmbstack/ripple/fixtures/forum/proto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func RouteAPI() {
	echoMux := ripple.Default().GetEcho()
	echoMux.GET("/", func(ctx echo.Context) error {
		userClient := proto.NewUserClient()
		req := &proto.GetInfoReq{
			Id: 1,
		}
		reply, _ := userClient.GetInfo(context.Background(), req)
		result := map[string]interface{}{
			"username": reply.Name,
		}
		return ctx.JSON(http.StatusOK, SuccessJSON(result))
	})

	//===========================================================
	//                      common
	//===========================================================
	echoMux.Use(middleware.Gzip())
	echoMux.Use(middleware.Secure())

	allowOriginsValue := []string{"*"}
	//allowOriginsValue := []string{"http://www.ripple.com", "http://api.ripple.com"}
	echoMux.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOriginsValue,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
	}))
	echoMux.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLength:  32,
		TokenLookup:  "header:X-CSRF-TOKEN",
		ContextKey:   "csrf",
		CookieName:   "_csrf",
		CookieMaxAge: 86400,
	}))

	//===========================================================
	//                      v1
	//===========================================================
	v1Group := echoMux.Group("/v1")

	homes := v1.HomeController{}
	homes.Register()

	citys := v1.CityController{Group: v1Group}
	citys.Setup()

	users := v1.UserController{Group: v1Group}
	users.Setup()
}
