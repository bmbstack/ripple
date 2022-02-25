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
	"strings"
	"sync"
)

var firstRegController = true
var firstRegModel = true
var line1 = "=============================="
var line2 = "================================"

// VersionName 0.8.0以后使用yaml配置文件
const VersionName = "0.8.0"

func Version() string {
	return VersionName
}

var ins *Ripple
var once sync.Once

func Default() *Ripple {
	once.Do(func() {
		ins = NewRipple()
	})
	return ins
}

// Ripple ripple struct
type Ripple struct {
	Logger *logger.Logger
	Echo   *echo.Echo
	Orms   map[string]*Orm
	Caches map[string]*cache.Cache
}

// NewLogger new a logger instance
func NewLogger() *logger.Logger {
	log, err := logger.NewLogger("ripple", 1, os.Stdout)
	if err != nil {
		panic(err) // Check for error
	}
	return log
}

// NewRipple new a ripple instance
func NewRipple() *Ripple {
	config := GetBaseConfig()

	r := &Ripple{}
	r.Logger = NewLogger()
	r.Echo = echo.New()

	r.Echo.Use(mw.Recover())
	r.Echo.Use(mw.Logger())

	r.Echo.Binder = binding.Binder{}
	r.Echo.Renderer = NewRenderer(config)
	r.Echo.Static("/static", config.Static)

	orms := make(map[string]*Orm)
	if IsNotEmpty(config.Databases) {
		for _, item := range config.Databases {
			orms[item.Alias] = NewOrm(item, !strings.EqualFold("prod", GetEnv()))
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

//================================================================
//                      Ripple func
//================================================================
// GetEcho  return echo
func (this *Ripple) GetEcho() *echo.Echo {
	return this.Echo
}

// GetOrm  return ripple model
func (this *Ripple) GetOrm(alias string) *Orm {
	if _, ok := this.Orms[alias]; !ok {
		panic(fmt.Errorf("GetOrm: cannot get orm alias '%s'", alias))
	}
	return this.Orms[alias]
}

// GetCache  return ripple cache
func (this *Ripple) GetCache(alias string) *cache.Cache {
	if _, ok := this.Caches[alias]; !ok {
		panic(fmt.Errorf("GetCache: cannot get cache alias '%s'", alias))
	}
	return this.Caches[alias]
}

// RegisterControllers register a controller for ripple App
func (this *Ripple) RegisterController(c Controller) {
	if firstRegController {
		fmt.Println(fmt.Sprintf("%s%s%s",
			color.White(line1),
			color.Bold(color.Green("Controller information")),
			color.White(line1)))
	}
	AddController(this.Echo, c)
	firstRegController = false
}

// RegisterModels registers models in the global ripple App.
func (this *Ripple) RegisterModels(orm *Orm, modelItems ...interface{}) {
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
func (this *Ripple) Run() {
	// autoMigrate all orms
	if GetBaseConfig().AutoMigrate {
		for alias := range this.Orms {
			this.Orms[alias].AutoMigrateAll()
		}
	}

	this.Logger.Info(fmt.Sprintf("Ripple ListenAndServe: %s", color.Green(GetBaseConfig().Domain)))
	this.Echo.Debug = !strings.EqualFold("prod", GetEnv())
	err := this.Echo.Start(GetBaseConfig().Domain)
	if err != nil {
		this.Logger.Error(fmt.Sprintf("Ripple Start error: %s", color.Red(err)))
	}
}
