package ripple

import (
	"fmt"

	"github.com/bmbstack/ripple/cache"
	"github.com/spf13/viper"
)

type BaseConfig struct {
	AutoMigrate bool             `mapstructure:"autoMigrate,omitempty"` // autoMigrate
	BindAllTag  bool             `mapstructure:"bindAllTag,omitempty"`  // false=binding
	Debug       bool             `mapstructure:"debug,omitempty"`       // debug, default false
	Domain      string           `mapstructure:"domain,omitempty"`      // domain=127.0.0.1:8090
	Static      string           `mapstructure:"static,omitempty"`      // static=frontend/static
	Templates   string           `mapstructure:"templates,omitempty"`   // templates=frontend/templates
	Databases   []DatabaseConfig `mapstructure:"databases,omitempty"`   // databases
	Caches      []CacheConfig    `mapstructure:"caches,omitempty"`      // caches
	Logs        []LogConfig      `mapstructure:"logs,omitempty"`        // logs
	Nacos       NacosConfig      `mapstructure:"nacos,omitempty"`       // nacos
	SLS         SlsConfig        `mapstructure:"sls,omitempty"`         // sls
	CLS         ClsConfig        `mapstructure:"cls,omitempty"`         // cls
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
	Alias        string `mapstructure:"alias"`        // alias=forum
	Section      string `mapstructure:"section"`      // section=forum
	Adapter      string `mapstructure:"adapter"`      // adapter=redis
	Host         string `mapstructure:"host"`         // host=127.0.0.1
	Port         int    `mapstructure:"port"`         // port=6379
	Password     string `mapstructure:"password"`     // password=123456
	DB           int    `mapstructure:"db"`           // db, select db
	PoolFIFO     bool   `mapstructure:"poolFifo"`     // poolFifo, true for FIFO pool, false for LIFO pool
	PoolSize     int    `mapstructure:"poolSize"`     // poolSize
	MinIdleConns int    `mapstructure:"minIdleConns"` // minIdleConns

	// redis://user:password@localhost:6789/3?dial_timeout=3&db=1&read_timeout=6s&max_retries=2
	Url string `mapstructure:"url"` // url, 如果使用url，其他配置无效,
}

type NacosConfig struct {
	// common
	Host        string `mapstructure:"host"`
	Port        uint64 `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespaceId"`
	Cluster     string `mapstructure:"cluster"`
	Group       string `mapstructure:"group"`
	CacheDir    string `mapstructure:"cacheDir"`
	LogDir      string `mapstructure:"logDir"`

	// server
	Server string `mapstructure:"server"`

	// client
	FailMode       string `mapstructure:"failMode"`
	SelectMode     string `mapstructure:"selectMode"`
	ClientPoolSize int    `mapstructure:"clientPoolSize"`
}

type LumberjackConfig struct {
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"maxSize"` // MB
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"` // days
	Compress   bool   `mapstructure:"compress"`
}

type LogConfig struct {
	Alias string `mapstructure:"alias"`
	Type  string `mapstructure:"type"` // sls, cls

	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	AllowLogLevel   string `mapstructure:"allowLogLevel"`
	CloseStdout     bool   `mapstructure:"closeStdout"`

	LumberjackConfig LumberjackConfig `mapstructure:"lumberjackConfig,omitempty"`

	Project  string `mapstructure:"project"`
	Logstore string `mapstructure:"logstore"`
	Topic    string `mapstructure:"topic"`
	Source   string `mapstructure:"source"`
}

type SlsConfig struct {
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	AllowLogLevel   string `mapstructure:"allowLogLevel"`
	CloseStdout     bool   `mapstructure:"closeStdout"`

	LumberjackConfig LumberjackConfig `mapstructure:"lumberjackConfig,omitempty"`

	Project  string `mapstructure:"project"`
	Logstore string `mapstructure:"logstore"`
	Topic    string `mapstructure:"topic"`
	Source   string `mapstructure:"source"`
}

type ClsConfig struct {
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	AllowLogLevel   string `mapstructure:"allowLogLevel"`
	CloseStdout     bool   `mapstructure:"closeStdout"`

	LumberjackConfig LumberjackConfig `mapstructure:"lumberjackConfig,omitempty"`

	Topic string `mapstructure:"topic"`
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
	fmt.Println(fmt.Sprintf("配置文件路径: %s", configPath))
	fmt.Println(fmt.Sprintf("执行环境: %s", env))
	e = env
	v = viper.New()
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	if err := v.ReadInConfig(); err != nil {
		fmt.Println(fmt.Sprintf("Viper ReadInConfig err:%s\n", err))
		panic(err)
	}

	baseConfig := BaseConfig{}
	err := v.Unmarshal(&baseConfig)
	if err != nil {
		fmt.Println("yaml parse err: ", err)
		panic(err)
	}
	bc = &baseConfig
}

func GetBaseConfig() *BaseConfig {
	if bc == nil {
		panic("Please init Config")
	}
	return bc
}

func UpdateBaseConfig() {
	if v == nil {
		panic("Please init Config")
	}
	baseConfig := BaseConfig{}
	err := v.Unmarshal(&baseConfig)
	if err != nil {
		fmt.Println("yaml parse err: ", err)
		panic(err)
	}
	bc = &baseConfig
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
		Addr:         fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port),
		Password:     cacheConfig.Password,
		DB:           cacheConfig.DB,
		PoolFIFO:     cacheConfig.PoolFIFO,
		PoolSize:     cacheConfig.PoolSize,
		MinIdleConns: cacheConfig.MinIdleConns,
	}
}
