package ripplepprof

import (
	"github.com/labstack/echo/v4"
	"net/http/pprof"
)

func Wrap(e *echo.Echo) {
	e.GET("/debug/pprof/", IndexHandler)
	e.GET("/debug/pprof/heap", HeapHandler)
	e.GET("/debug/pprof/goroutine", GoroutineHandler)
	e.GET("/debug/pprof/block", BlockHandler)
	e.GET("/debug/pprof/threadcreate", ThreadCreateHandler)
	e.GET("/debug/pprof/cmdline", CmdlineHandler)
	e.GET("/debug/pprof/profile", ProfileHandler)
	e.GET("/debug/pprof/symbol", SymbolHandler)
}

var Wrapper = Wrap

func IndexHandler(ctx echo.Context) error {
	pprof.Index(ctx.Response(), ctx.Request())
	return nil
}

func HeapHandler(ctx echo.Context) error {
	pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func GoroutineHandler(ctx echo.Context) error {
	pprof.Handler("goroutine").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func BlockHandler(ctx echo.Context) error {
	pprof.Handler("block").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func ThreadCreateHandler(ctx echo.Context) error {
	pprof.Handler("threadcreate").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func CmdlineHandler(ctx echo.Context) error {
	pprof.Cmdline(ctx.Response(), ctx.Request())
	return nil
}

func ProfileHandler(ctx echo.Context) error {
	pprof.Profile(ctx.Response(), ctx.Request())
	return nil
}

func SymbolHandler(ctx echo.Context) error {
	pprof.Symbol(ctx.Response(), ctx.Request())
	return nil
}
