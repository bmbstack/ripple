package ripple

import (
	"fmt"
	"github.com/bmbstack/ripple/cache"
	. "github.com/bmbstack/ripple/helper"
	"github.com/bmbstack/ripple/middleware/binding"
	"github.com/bmbstack/ripple/middleware/logger"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"os"
)

var Logger *logger.Logger
var baseRipple *Ripple

var firstRegController = true
var firstRegModel = true
var line1 = "=============================="
var line2 = "================================"

const VersionName = "0.4.0"

// Init init ripple
func init() {
	Logger = NewLogger()
	baseRipple = NewRipple()
}

func Version() string {
	return VersionName
}

// Ripple ripple struct
type Ripple struct {
	Echo   *echo.Echo
	Config *Config
	Orms   map[string]*Orm
	Caches map[string]*cache.Cache
}

func NewLogger() *logger.Logger {
	log, err := logger.NewLogger("ripple", 1, os.Stdout)
	if err != nil {
		panic(err) // Check for error
	}
	return log
}

// NewRipple new a ripple instance
func NewRipple() *Ripple {
	config := NewConfig()
	r := &Ripple{}

	r.Config = config
	r.Echo = echo.New()

	r.Echo.Use(mw.Recover())
	r.Echo.Use(mw.Logger())

	r.Echo.Binder = binding.Binder{}
	r.Echo.Renderer = NewRenderer(config)
	r.Echo.Static("/static", config.Static)

	orms := make(map[string]*Orm)
	if IsNotEmpty(config.Databases) {
		for _, item := range config.Databases {
			orms[item.Alias] = NewOrm(item, r.Config.DebugOn)
		}
	}
	r.Orms = orms

	caches := make(map[string]*cache.Cache)
	if IsNotEmpty(config.Caches) {
		for _, item := range config.Caches {
			newCache, err := cache.NewCache(item.Alias, cache.Options{
				Adapter:       item.Adapter,
				AdapterConfig: item.GetCacheAdapterConfig(),
				Section:       item.Section,
			})
			if err != nil {
				fmt.Println(fmt.Sprintf("Connect.cache error: %s", err.Error()))
			} else {
				caches[item.Alias] = newCache
			}
		}
	}
	r.Caches = caches
	return r
}

// GetEcho  return echo
func GetEcho() *echo.Echo {
	return baseRipple.Echo
}

// GetConfig return config
func GetConfig() *Config {
	return baseRipple.Config
}

// GetOrm  return ripple model
func GetOrm(alias string) *Orm {
	if _, ok := baseRipple.Orms[alias]; !ok {
		panic(fmt.Errorf("GetOrm: cannot get orm alias '%s'", alias))
	}
	return baseRipple.Orms[alias]
}

// GetCache  return ripple cache
func GetCache(alias string) *cache.Cache {
	if _, ok := baseRipple.Caches[alias]; !ok {
		panic(fmt.Errorf("GetCache: cannot get cache alias '%s'", alias))
	}
	return baseRipple.Caches[alias]
}

// RegisterControllers register a controller for ripple App
func RegisterController(c Controller) {
	if firstRegController {
		fmt.Println(fmt.Sprintf("%s%s%s",
			color.White(line1),
			color.Bold(color.Green("Controller information")),
			color.White(line1)))
	}
	AddController(baseRipple.Echo, c)
	firstRegController = false
}

// RegisterModels registers models in the global ripple App.
func RegisterModels(orm *Orm, modelItems ...interface{}) {
	if firstRegModel {
		fmt.Println(fmt.Sprintf("%s%s%s",
			color.White(line2),
			color.Bold(color.Green("Orm information")),
			color.White(line2)))
	}
	_ = orm.AddModels(modelItems...)
	firstRegModel = false
}

// Run run ripple application
func Run() {
	for alias := range baseRipple.Orms {
		baseRipple.Orms[alias].AutoMigrateAll()
	}
	Logger.Info(fmt.Sprintf("Ripple ListenAndServe: %s", color.Green(baseRipple.Config.Domain)))
	baseRipple.Echo.Debug = baseRipple.Config.DebugOn
	err := baseRipple.Echo.Start(baseRipple.Config.Domain)
	if err != nil {
		Logger.Error(fmt.Sprintf("Ripple Start error: %s", color.Red(err)))
	}
}
