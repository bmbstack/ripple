package controllers

import (
	"context"
	"fmt"
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/internal/controllers/v1"
	. "github.com/bmbstack/ripple/fixtures/forum/internal/helper"
	"github.com/bmbstack/ripple/fixtures/forum/internal/services"
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
		for i := 0; i < 50; i++ {
			pc := services.GetPlayerClient()
			reply, _ := pc.GetPlayerInfo(context.Background(), &proto.GetPlayerInfoReq{PlayerId: 1377693})
			fmt.Println(reply.Nickname)
		}

		pc := services.GetPlayerClient()
		reply, _ := pc.GetPlayerInfo(context.Background(), &proto.GetPlayerInfoReq{PlayerId: 1377693})
		fmt.Println("services.GetPlayerClient().GetPlayerInfo.Nickname=====>", reply.Nickname)

		sc := services.GetStudentClient()
		reply2, _ := sc.Learn(context.Background(), &proto.LearnReq{Id: 1})
		result := map[string]interface{}{
			"remoteName": reply.Nickname,
			"localName":  reply2.Name,
		}
		logger.With(nil).Info(fmt.Sprintf("time: %d, name333: %s", time.Now().Unix(), reply.Nickname))
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
