package ripple

import (
	_ "github.com/joho/godotenv/autoload"
	"path/filepath"
	"os"
	"fmt"
	"github.com/labstack/gommon/color"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	DebugOn   bool       `json:"debugOn"`             // debug_on=true
	Domain    string     `json:"domain"`              // domain=127.0.0.1:8090
	Static    string     `json:"static"`              // static=frontend/static
	Templates string     `json:"templates"`           // templates=frontend/templates
	Databases []Database `json:"databases,omitempty"` // databases
	Caches    []Cache    `json:"caches,omitempty"`    // caches
}

type Database struct {
	Alias    string `json:"alias"`    // alias=forum
	Dialect  string `json:"dialect"`  // dialect=mysql
	Host     string `json:"host"`     // host=127.0.0.1
	Port     int    `json:"port"`     // port=3306
	Name     string `json:"name"`     // name=forum
	Username string `json:"username"` // username=root
	Password string `json:"password"` // password=123456
}

type Cache struct {
	Alias    string `json:"alias"`    // alias=forum
	Section  string `json:"section"`  // alias=ripple
	Adapter  string `json:"adapter"`  // adapter=redis
	Host     string `json:"host"`     // host=127.0.0.1
	Port     int    `json:"port"`     // port=6379
	Password string `json:"password"` // password=Bmbstack2016
}

// CacheAdapterConfig
type CacheAdapterConfig struct {
	Addr   string `json:"Addr"`
	Passwd string `json:"Passwd"`
}

// GetCacheConfig return cache config
func (cache Cache) GetCacheConfig() string {
	adapterConfig := &CacheAdapterConfig{
		Addr:   fmt.Sprintf("%s:%d", cache.Host, cache.Port),
		Passwd: cache.Password,
	}
	configByte, err := json.Marshal(adapterConfig)
	if err != nil {
		panic(err)
	}
	return string(configByte)
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
