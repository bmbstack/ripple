package ripplepprof

import (
	"net/http/pprof"
	"github.com/labstack/echo"
)

func Wrap(e *echo.Echo) {
	e.Get("/debug/pprof/", IndexHandler)
	e.Get("/debug/pprof/heap", HeapHandler)
	e.Get("/debug/pprof/goroutine", GoroutineHandler)
	e.Get("/debug/pprof/block", BlockHandler)
	e.Get("/debug/pprof/threadcreate", ThreadCreateHandler)
	e.Get("/debug/pprof/cmdline", CmdlineHandler)
	e.Get("/debug/pprof/profile", ProfileHandler)
	e.Get("/debug/pprof/symbol", SymbolHandler)
}

var Wrapper = Wrap

func IndexHandler(ctx *echo.Context) error {
	pprof.Index(ctx.Response(), ctx.Request())
	return nil
}

func HeapHandler(ctx *echo.Context) error {
	pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func GoroutineHandler(ctx *echo.Context) error {
	pprof.Handler("goroutine").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func BlockHandler(ctx *echo.Context) error {
	pprof.Handler("block").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func ThreadCreateHandler(ctx *echo.Context) error {
	pprof.Handler("threadcreate").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func CmdlineHandler(ctx *echo.Context) error {
	pprof.Cmdline(ctx.Response(), ctx.Request())
	return nil
}

func ProfileHandler(ctx *echo.Context) error {
	pprof.Profile(ctx.Response(), ctx.Request())
	return nil
}

func SymbolHandler(ctx *echo.Context) error {
	pprof.Symbol(ctx.Response(), ctx.Request())
	return nil
}
