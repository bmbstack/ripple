package ripple

import (
	"fmt"
	"github.com/bmbstack/ripple/middleware/binding"
	"github.com/bmbstack/ripple/middleware/logger"
	"github.com/bmbstack/ripple/middleware/cache"
	. "github.com/bmbstack/ripple/helper"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	"os"
)

var Logger *logger.Logger
var baseRipple *Ripple

var firstRegController = true
var firstRegModel = true
var line1 = "=============================="
var line2 = "================================"

// Ripple ripple struct
type Ripple struct {
	Echo   *echo.Echo
	Config *Config
	Orms   map[string]*Orm
}

func init() {
	Logger = NewLogger()
	baseRipple = NewRipple()
}

func NewLogger() *logger.Logger {
	log, err := logger.NewLogger("ripple", 1, os.Stdout)
	if err != nil {
		log.Error(err.Error())
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
		for _, database := range config.Databases {
			orms[database.Alias] = NewOrm(database)
		}
	}
	r.Orms = orms

	if IsNotEmpty(config.Caches) {
		for _, itemCache := range config.Caches {
			r.Echo.Use(cache.EchoCacher(itemCache.Alias, cache.Options{
				Adapter:       itemCache.Adapter,
				AdapterConfig: itemCache.GetCacheConfig(),
				Section:       itemCache.Section, Interval: 5}))
		}
	}
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

// RegisterControllers register a controller for ripple App
func RegisterController(c Controller) {
	if firstRegController {
		fmt.Println(fmt.Sprintf("%s%s%s",
			color.White(line1),
			color.Bold(color.Green("Controller information")),
			color.White(line1)))
	}
	AddController(*baseRipple.Echo, c)
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
	orm.AddModels(modelItems...)
	firstRegModel = false
}

// Run run ripple application
func Run() {
	for alias := range baseRipple.Orms {
		baseRipple.Orms[alias].AutoMigrateAll()
	}
	Logger.Info(fmt.Sprintf("Ripple ListenAndServe: %s", color.Green(baseRipple.Config.Domain)))
	baseRipple.Echo.Debug = baseRipple.Config.DebugOn
	baseRipple.Echo.Start(baseRipple.Config.Domain)
}
