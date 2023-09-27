package ripple

import (
	"bufio"
	"fmt"
	"github.com/bmbstack/g"
	"github.com/bmbstack/ripple/cache"
	. "github.com/bmbstack/ripple/helper"
	log "github.com/bmbstack/ripple/logger"
	"github.com/bmbstack/ripple/logger/cls"
	"github.com/bmbstack/ripple/logger/sls"
	"github.com/bmbstack/ripple/middleware/bind"
	"github.com/bmbstack/ripple/middleware/binding"
	"github.com/bmbstack/ripple/middleware/logger"
	"github.com/bmbstack/ripple/middleware/recover"
	"github.com/bmbstack/ripple/nacos/rpcxnacos/serverplugin"
	"github.com/bmbstack/ripple/util"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var firstRegController = true
var firstRegModel = true
var line1 = "=============================="
var line2 = "================================"

type LogType string

const (
	LogTypeSLS = "sls"
	LogTypeCLS = "cls"
)

// VersionName 0.8.2以后使用yaml配置文件, 1.0.1升级了脚手架(protoc, ast gen)
const VersionName = "1.2.4"

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
	Logger              *logger.Logger
	Logrus              *logrus.Logger
	Echo                *echo.Echo
	Orms                map[string]*Orm
	Caches              map[string]*cache.Cache
	RpcServer           *server.Server
	NacosRegisterPlugin *serverplugin.NacosRegisterPlugin
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

	g.SetDefaultRecoverFunc(func(v interface{}, done chan<- error) {
		// 上报panic日志到阿里云......
		stackArr := make([]byte, 2048)
		length := runtime.Stack(stackArr, true)
		stack := string(stackArr[:length])
		msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", v, stack)
		log.With(nil).Error(msg)
		done <- fmt.Errorf("panic in goroutine successfully recovered")
	})

	r := &Ripple{}
	r.Logger = NewLogger()
	r.Logrus = log.StandardLogger()
	r.Echo = echo.New()

	r.Echo.Use(recover.Recover())
	r.Echo.Use(mw.Logger())

	r.Echo.Binder = binding.Binder{}
	if config.BindAllTag {
		r.Echo.Binder = &bind.DefaultBinder{}
	}
	r.Echo.Renderer = NewRenderer(config)
	r.Echo.Static("/static", config.Static)

	// orm
	orms := make(map[string]*Orm)
	if IsNotEmpty(config.Databases) {
		for _, item := range config.Databases {
			orms[item.Alias] = NewOrm(item, GetBaseConfig().Debug)
		}
	}
	r.Orms = orms

	// cache
	caches := make(map[string]*cache.Cache)
	if IsNotEmpty(config.Caches) {
		for _, item := range config.Caches {
			newCache, err := cache.NewCache(cache.Options{
				Alias:         item.Alias,
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

// AddLogType  add log type (ripple.LogTypeSLS, ripple.LogTypeCLS)
func (this *Ripple) AddLogType(value LogType) {
	fmt.Println(color.Green(fmt.Sprintf("Logger: Add LogType %s", value)))

	conf := GetBaseConfig()
	formatter := &logrus.JSONFormatter{
		DisableHTMLEscape: true,
	}
	if LogTypeSLS == value && IsNotEmpty(conf.SLS) {
		h := sls.NewSLSHook(
			conf.SLS.AccessKeyId,
			conf.SLS.AccessKeySecret,
			conf.SLS.Endpoint,
			conf.SLS.AllowLogLevel,
			sls.SetProject(conf.SLS.Project),
			sls.SetLogstore(conf.SLS.Logstore),
			sls.SetTopic(conf.SLS.Topic),
			sls.SetSource(conf.SLS.Source),
		)
		if conf.SLS.CloseStdout {
			f, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				fmt.Println("SLS.CloseStdout Open file err: ", err)
			}
			this.Logrus.SetOutput(bufio.NewWriter(f))
		}
		this.Logrus.SetFormatter(formatter)
		this.Logrus.AddHook(h)
	} else if LogTypeCLS == value && IsNotEmpty(conf.CLS) {
		h := cls.NewCLSHook(
			conf.CLS.AccessKeyId,
			conf.CLS.AccessKeySecret,
			conf.CLS.Endpoint,
			conf.CLS.AllowLogLevel,
			cls.SetTopic(conf.CLS.Topic),
		)
		if conf.CLS.CloseStdout {
			f, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				fmt.Println("CLS.CloseStdout Open file err: ", err)
			}
			this.Logrus.SetOutput(bufio.NewWriter(f))
		}
		this.Logrus.SetFormatter(formatter)
		this.Logrus.AddHook(h)
	} else {
		// do nothing
	}
}

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

// RegisterController register a controller for ripple App
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

// checkOrNewRpcServer check or new rpc server
func (this *Ripple) checkOrNewRpcServer() {
	config := GetBaseConfig()
	if this.RpcServer == nil && IsNotEmpty(config.Nacos) {
		this.RpcServer, this.NacosRegisterPlugin = NewRpcServerNacos(config.Nacos)
	}
}

// RegisterRpc register rpc service
func (this *Ripple) RegisterRpc(name string, rpc interface{}, metadata string) {
	this.checkOrNewRpcServer()
	if this.RpcServer != nil {
		err := this.RpcServer.RegisterName(name, rpc, metadata)
		if err != nil {
			this.Logger.Error(fmt.Sprintf("Rpc register service error: %s", err.Error()))
		} else {
			this.Logger.Notice(fmt.Sprintf("Rpc register service success: %s, %v", name, rpc))
		}
	}
}

// RunRpc run rpc server
func (this *Ripple) RunRpc() {
	this.checkOrNewRpcServer()
	if this.RpcServer != nil {
		conf := GetBaseConfig()
		if IsNotEmpty(conf.Nacos) {
			if !strings.Contains(conf.Nacos.Server, ":") {
				this.Logger.Error("Rpc run error: nacos server address format is wrong, must contains `:`")
				return
			}
			arr := strings.Split(conf.Nacos.Server, ":")
			address := fmt.Sprintf("%s:%s", util.InternalIP(), arr[len(arr)-1:][0])
			this.Logger.Notice(fmt.Sprintf("Rpc run, address: %s", address))
			go func() {
				err := this.RpcServer.Serve("tcp", address)
				if err != nil {
					this.Logger.Error(fmt.Sprintf("Rpc run error: %s", err.Error()))
				} else {
					this.Logger.Notice(fmt.Sprintf("Rpc run success, address: %s", address))
				}
			}()
		}
	}
}

// UnregisterRpc unregisters all rpc services.
func (this *Ripple) UnregisterRpc() {
	if this.RpcServer != nil {
		err := this.RpcServer.UnregisterAll()
		if err != nil {
			this.Logger.Error(fmt.Sprintf("Rpc unregisters all services error: %s", err.Error()))
		} else {
			this.Logger.Notice("Rpc unregisters all service success")
		}
	}
}

// StopRpc stop rpc server, close all rpc connections
func (this *Ripple) StopRpc() {
	if this.RpcServer != nil {
		err := this.NacosRegisterPlugin.Stop()
		if err != nil {
			this.Logger.Error(fmt.Sprintf("Rpc nacosRegisterPlugin stop error: %s", err.Error()))
		} else {
			this.Logger.Notice("Rpc nacosRegisterPlugin stop success")
		}
		err = this.RpcServer.Close()
		if err != nil {
			this.Logger.Error(fmt.Sprintf("Rpc server close error: %s", err.Error()))
		} else {
			this.Logger.Notice("Rpc server close success")
		}
	}
}

func (this *Ripple) CloseOrm() {
	if IsNotEmpty(this.Orms) {
		for key, item := range this.Orms {
			err := item.Close()
			if err != nil {
				this.Logger.Error(fmt.Sprintf("Close orm (alias: %s) error: %s", key, err.Error()))
			} else {
				this.Logger.Notice(fmt.Sprintf("Close orm (alias: %s) success", key))
			}
		}
	}
}

func (this *Ripple) CloseCache() {
	if IsNotEmpty(this.Caches) {
		for key, item := range this.Caches {
			err := item.Close()
			if err != nil {
				this.Logger.Error(fmt.Sprintf("Close cache (alias: %s) error: %s", key, err.Error()))
			} else {
				this.Logger.Notice(fmt.Sprintf("Close cache (alias: %s) success", key))
			}
		}
	}
}

// Run run ripple application
func (this *Ripple) Run() {
	this.RunWith(GetBaseConfig().Domain)
}

// RunWith run ripple application
func (this *Ripple) RunWith(domain string) {
	// autoMigrate all orms
	if GetBaseConfig().AutoMigrate {
		for alias := range this.Orms {
			this.Orms[alias].AutoMigrateAll()
		}
	}

	this.Logger.Info(fmt.Sprintf("Ripple ListenAndServe: %s", color.Green(domain)))
	this.Echo.Debug = GetBaseConfig().Debug
	err := this.Echo.Start(domain)
	if err != nil {
		this.Logger.Error(fmt.Sprintf("Ripple Start error: %s", color.Red(err)))
	}
}

// RunScript run script
func RunScript(commands []string) {
	entireScript := strings.NewReader(strings.Join(commands, "\n"))
	bash := exec.Command("/bin/bash")
	stdin, _ := bash.StdinPipe()
	stdout, _ := bash.StdoutPipe()
	stderr, _ := bash.StderrPipe()

	wait := sync.WaitGroup{}
	wait.Add(3)
	go func() {
		_, _ = io.Copy(stdin, entireScript)
		_ = stdin.Close()
		wait.Done()
	}()
	go func() {
		_, _ = io.Copy(os.Stdout, stdout)
		wait.Done()
	}()
	go func() {
		_, _ = io.Copy(os.Stderr, stderr)
		wait.Done()
	}()

	_ = bash.Start()
	wait.Wait()
	_ = bash.Wait()
}
