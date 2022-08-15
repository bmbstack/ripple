package controllers

import (
	"context"
	"fmt"
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/internal/controllers/v1"
	. "github.com/bmbstack/ripple/fixtures/forum/internal/helper"
	"github.com/bmbstack/ripple/fixtures/forum/internal/rpcclient"
	"github.com/bmbstack/ripple/fixtures/forum/proto"
	"github.com/bmbstack/ripple/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func RouteAPI() {
	echoMux := ripple.Default().GetEcho()
	echoMux.GET("/", func(ctx echo.Context) error {
		reply1, _ := rpcclient.GetPlayerClient().GetPlayerInfo(context.Background(), &proto.GetPlayerInfoReq{PlayerId: 1377693})
		reply2, _ := rpcclient.GetStudentClient().Learn(context.Background(), &proto.LearnReq{Id: 1})
		reply3, _ := rpcclient.GetTeacherClient().Teach(context.Background(), &proto.TeachReq{Id: 1})
		result := map[string]interface{}{
			"remoteName":        reply1.Nickname,
			"localName.student": reply2.Name,
			"localName.teacher": reply3.Name,
		}
		logger.With(nil).Info(fmt.Sprintf("time: %d, value1: %s", time.Now().Unix(), reply1.Nickname))
		logger.With(nil).Info(fmt.Sprintf("time: %d, value2: %s", time.Now().Unix(), reply2.Name))
		logger.With(nil).Info(fmt.Sprintf("time: %d, value3: %s", time.Now().Unix(), reply3.Name))
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
	//echoMux.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	//	TokenLength:  32,
	//	TokenLookup:  "header:X-CSRF-TOKEN",
	//	ContextKey:   "csrf",
	//	CookieName:   "_csrf",
	//	CookieMaxAge: 86400,
	//}))

	//===========================================================
	//                      v1
	//===========================================================
	v1Group := echoMux.Group("/v1")

	homes := v1.HomeController{}
	homes.Register()

	citys := v1.CityController{Group: v1Group}
	citys.Setup()

	students := v1.StudentController{Group: v1Group}
	students.Setup()

	teachers := v1.TeacherController{Group: v1Group}
	teachers.Setup()
}
