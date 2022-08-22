package ripple

import (
	"errors"
	"fmt"
	. "github.com/bmbstack/ripple/helper"
	"github.com/labstack/gommon/color"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	mlogger "gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Orm facilitate database interactions, support mysql
type Orm struct {
	models map[string]reflect.Value
	*gorm.DB
}

// NewOrm creates a new model, and opens database connection based on cfg settings
func NewOrm(database DatabaseConfig, debug bool) *Orm {
	orm := &Orm{
		models: make(map[string]reflect.Value),
	}

	dialect := database.Dialect
	host := database.Host
	port := database.Port
	name := database.Name
	username := database.Username
	password := database.Password
	maxIdleConns := database.MaxIdleConns
	maxOpenConns := database.MaxOpenConns

	// logger
	logLevel := mlogger.Silent
	logColorful := false
	if debug {
		logLevel = mlogger.Info
		logColorful = true
	}
	newLogger := mlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		mlogger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLevel,    // Log level
			Colorful:      logColorful, // Disable color
		},
	)

	dsn := ""
	switch dialect {
	case "mysql":
		dsn = username + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN: dsn,
		}), &gorm.Config{Logger: newLogger})
		if err != nil {
			fmt.Println(fmt.Sprintf("%s: %s", color.Red(fmt.Sprintf("Connect.%s, error", dialect)), dsn))
			panic(err)
		}
		orm.DB = db
	case "postgresql":
		dsn = "host=" + host + " user=" + username + " password=" + password + " dbname=" + name + " port=" + strconv.Itoa(port) + " sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			fmt.Println(fmt.Sprintf("%s: %s", color.Red(fmt.Sprintf("Connect.%s, error", dialect)), dsn))
			panic(err)
		}
		orm.DB = db
	default:
		dsn = username + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{Logger: newLogger})
		if err != nil {
			fmt.Println(fmt.Sprintf("%s: %s", color.Red(fmt.Sprintf("Connect.%s, error", dialect)), dsn))
			panic(err)
		}
		orm.DB = db
	}
	fmt.Println(fmt.Sprintf("%s: %s", color.Green(fmt.Sprintf("Connect.%s", dialect)), dsn))

	sqlDB, err := orm.DB.DB()
	if err != nil {
		fmt.Println(color.Red(fmt.Sprintf("sqlDB config error.%s", err.Error())))
		return orm
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	return orm
}

// AutoMigrateAll runs migrations for all the registered models
func (orm *Orm) AutoMigrateAll() {
	for _, v := range orm.models {
		_ = orm.AutoMigrate(v.Interface())
	}
}

// AddModels add the values to the models registry
func (orm *Orm) AddModels(values ...interface{}) error {
	// do not work on them.models first, this is like an insurance policy
	// whenever we encounter any error in the values nothing goes into the registry
	models := make(map[string]reflect.Value)
	if len(values) > 0 {
		for _, val := range values {
			rVal := reflect.ValueOf(val)
			if rVal.Kind() == reflect.Ptr {
				rVal = rVal.Elem()
			}
			switch rVal.Kind() {
			case reflect.Struct:
				models[GetTypeName(rVal.Type())] = reflect.New(rVal.Type())
				fmt.Println(fmt.Sprintf("%s: %v", color.Bold("[RegisterModel]"), color.Bold(color.Blue(rVal.Type()))))
			default:
				return errors.New("ripple: model must be struct type")
			}
		}
	}
	for k, v := range models {
		orm.models[k] = v
	}
	return nil
}

func (orm *Orm) Close() error {
	db, err := orm.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
