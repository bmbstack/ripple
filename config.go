package ripple

import (
	"fmt"
	"github.com/bmbstack/ripple/cache"
	"github.com/spf13/viper"
)

type BaseConfig struct {
	AutoMigrate bool             `mapstructure:"autoMigrate,omitempty"` // autoMigrate
	Domain      string           `mapstructure:"domain,omitempty"`      // domain=127.0.0.1:8090
	Static      string           `mapstructure:"static,omitempty"`      // static=frontend/static
	Templates   string           `mapstructure:"templates,omitempty"`   // templates=frontend/templates
	Databases   []DatabaseConfig `mapstructure:"databases,omitempty"`   // databases
	Caches      []CacheConfig    `mapstructure:"caches,omitempty"`      // caches
	Nacos       Nacos            `mapstructure:"nacos,omitempty"`       // nacos
}

type DatabaseConfig struct {
	Alias        string `mapstructure:"alias"`        // alias=forum
	Dialect      string `mapstructure:"dialect"`      // dialect=mysql
	Host         string `mapstructure:"host"`         // host=127.0.0.1
	Port         int    `mapstructure:"port"`         // port=3306
	Name         string `mapstructure:"name"`         // name=forum
	Username     string `mapstructure:"username"`     // username=root
	Password     string `mapstructure:"password"`     // password=123456
	MaxIdleConns int    `mapstructure:"maxIdleConns"` // maxIdleConns
	MaxOpenConns int    `mapstructure:"maxOpenConns"` // maxOpenConns
}

type CacheConfig struct {
	Alias    string `mapstructure:"alias"`    // alias=forum
	Section  string `mapstructure:"section"`  // section=forum
	Adapter  string `mapstructure:"adapter"`  // adapter=redis
	Host     string `mapstructure:"host"`     // host=127.0.0.1
	Port     int    `mapstructure:"port"`     // port=6379
	Password string `mapstructure:"password"` // password=123456
	DB       int    `mapstructure:"db"`       // db, select db
}

type Nacos struct {
	Host        string `mapstructure:"host"`
	Port        uint64 `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespaceId"`
	Cluster     string `mapstructure:"cluster"`
	Group       string `mapstructure:"group"`
	CacheDir    string `mapstructure:"cacheDir"`
	LogDir      string `mapstructure:"logDir"`
	Server      string `mapstructure:"server"`
}

var (
	e  string
	v  *viper.Viper
	bc *BaseConfig
)

func InitConfig(env string) {
	InitConfigWithPath(env, "./config/")
}

func InitConfigWithPath(env string, configPath string) {
	fmt.Println(fmt.Sprintf("执行环境: %s", env))
	e = env
	v = viper.New()
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	if err := v.ReadInConfig(); err != nil {
		fmt.Println(fmt.Sprintf("Viper ReadInConfig err:%s\n", err))
	}

	baseConfig := BaseConfig{}
	err := v.Unmarshal(&baseConfig)
	if err != nil {
		fmt.Println("yaml parse err: ", err)
	}
	bc = &baseConfig
}

func GetBaseConfig() *BaseConfig {
	if bc == nil {
		panic("Please init Config")
	}
	return bc
}

func GetConfig() *viper.Viper {
	if v == nil {
		panic("Please init Config")
	}
	return v
}

func GetEnv() string {
	if e == "" {
		panic("Please init Config and Env")
	}
	return e
}

func (cacheConfig CacheConfig) GetCacheAdapterConfig() cache.AdapterConfig {
	return cache.AdapterConfig{
		Addr:     fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port),
		Password: cacheConfig.Password,
		DB:       cacheConfig.DB,
	}
}
