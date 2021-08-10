package ripple

import (
	"encoding/json"
	"fmt"
	"github.com/bmbstack/ripple/cache"
	"github.com/labstack/gommon/color"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	DebugOn     bool             `json:"debugOn"`             // debugOn
	AutoMigrate bool             `json:"autoMigrate"`         // autoMigrate
	Domain      string           `json:"domain"`              // domain=127.0.0.1:8090
	Static      string           `json:"static"`              // static=frontend/static
	Templates   string           `json:"templates"`           // templates=frontend/templates
	Databases   []DatabaseConfig `json:"databases,omitempty"` // databases
	Caches      []CacheConfig    `json:"caches,omitempty"`    // caches
}

type DatabaseConfig struct {
	Alias    string `json:"alias"`    // alias=forum
	Dialect  string `json:"dialect"`  // dialect=mysql
	Host     string `json:"host"`     // host=127.0.0.1
	Port     int    `json:"port"`     // port=3306
	Name     string `json:"name"`     // name=forum
	Username string `json:"username"` // username=root
	Password string `json:"password"` // password=123456
}

type CacheConfig struct {
	Alias    string `json:"alias"`    // alias=forum
	Section  string `json:"section"`  // section=forum
	Adapter  string `json:"adapter"`  // adapter=redis
	Host     string `json:"host"`     // host=127.0.0.1
	Port     int    `json:"port"`     // port=6379
	Password string `json:"password"` // password=123456
}

// GetCacheConfig return cache config
func (cacheConfig CacheConfig) GetCacheAdapterConfig() cache.AdapterConfig {
	return cache.AdapterConfig{
		Addr:     fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port),
		Password: cacheConfig.Password,
	}
}

func NewConfig() *Config {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config := &Config{}
	data, err := ioutil.ReadFile(filepath.Join(workingDir, "config.json"))
	if err != nil {
		data, err = ioutil.ReadFile(filepath.Join(workingDir, "config.json.example"))
		if err != nil {
			Logger.Info(fmt.Sprintf(color.Red("config.json or config.json.example is not exist!")))
			panic(err)
		}
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		Logger.Info(fmt.Sprintf(color.Red("config file parse error!")))
		panic(err)
	}

	// static
	if !filepath.IsAbs(config.Static) {
		config.Static = Getwd(config.Static)
	}
	// templates
	if !filepath.IsAbs(config.Templates) {
		config.Templates = Getwd(config.Templates)
	}

	return config
}

// Getwd return the path's abs path string
func Getwd(path string) string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(workingDir, path)
}
