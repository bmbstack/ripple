package controllers

import (
	"github.com/bmbstack/ripple"
	v12 "github.com/bmbstack/ripple/cmd/ripple/templates/internal/controllers/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RouteAPI() {
	echoMux := ripple.Default().GetEcho()

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

	homes := v12.HomeController{}
	homes.Register()

	citys := v12.CityController{Group: v1Group}
	citys.Setup()
}