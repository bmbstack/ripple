package cache

import (
	"github.com/labstack/echo"
	"fmt"
)

const EchoCacheStoreKey = "EchoCacheStore"

func Store(alias string, value interface{}) Cache {
	var cacher Cache
	switch v := value.(type) {
	case echo.Context:
		contextKey := fmt.Sprintf("%s_%s", EchoCacheStoreKey, alias)
		cacher = v.Get(contextKey).(Cache)
		if cacher == nil {
			panic("EchoStore not found, forget to Use Middleware ?")
		}
	default:
		panic("unknown Context")
	}

	if cacher == nil {
		panic("cache context not found")
	}
	return cacher
}

func EchoCacher(alias string, opt Options) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cacher, err := New(alias, opt)
			if err != nil {
				return err
			}
			contextKey := fmt.Sprintf("%s_%s", EchoCacheStoreKey, alias)
			ctx.Set(contextKey, cacher)

			if err = next(ctx); err != nil {
				return err
			}
			return nil
		}
	}
}
