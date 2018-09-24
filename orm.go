package ripple

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/color"
	"reflect"
	. "github.com/bmbstack/ripple/helper"
	"strconv"
)

// Orm facilitate database interactions, support mysql
type Orm struct {
	models map[string]reflect.Value
	*gorm.DB
}

// NewOrm creates a new model, and opens database connection based on cfg settings
func NewOrm(database Database) *Orm {
	orm := &Orm{
		models: make(map[string]reflect.Value),
	}

	dialect := database.Dialect
	host := database.Host
	port := database.Port
	name := database.Name
	username := database.Username
	password := database.Password

	connURI := ""
	switch dialect {
	case "mysql":
		connURI = username + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + name + "?charset=utf8&parseTime=True&loc=Local"
	default:
		dialect = "mysql"
		connURI = username + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + name + "?charset=utf8&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(dialect, connURI)
	if err != nil {
		Logger.Info(fmt.Sprintf("%s: %s, %s", color.Red("gorm.Open error"), dialect, connURI))
		panic(err)
	}
	orm.DB = db
	Logger.Info(fmt.Sprintf("%s: %s, %s", color.Green("gorm.Open"), dialect, connURI))
	return orm
}

// AutoMigrateAll runs migrations for all the registered models
func (orm *Orm) AutoMigrateAll() {
	for _, v := range orm.models {
		orm.AutoMigrate(v.Interface())
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


