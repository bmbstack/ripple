package one

import (
	"fmt"
	"encoding/json"
	"github.com/bmbstack/ripple"
	. "github.com/bmbstack/ripple/helper"
	"github.com/labstack/echo"
)

type BmbCity struct {
	ID   int64  `json:"ID" gorm:"column:id; type:int(11); primary_key; auto_increment"`
	Name string `json:"name" gorm:"column:name; type:varchar(50); not null; default:'北京'"`
	Pid  int64  `json:"pid" gorm:"column:pid; type:tinyint(4); not null; default:1"`
	BaseModel
}

func init() {
	ripple.RegisterModels(Orm, &BmbCity{})
}

func (BmbCity) TableName() string {
	return "bmb_city"
}

func (this *BmbCity) FindCityListByPid(ctx echo.Context, pid int64) (list []BmbCity) {
	cacheKey := fmt.Sprintf("FindCityListByPid:%d", pid)
	cacheValue := this.GetCache(ctx, cacheKey)
	if IsNotEmpty(cacheValue) {
		json.Unmarshal([]byte(cacheValue), &list)
	} else {
		Orm.DB.Where("pid=? AND is_deleted=0", pid).Find(&list)
		this.SetCache(ctx, cacheKey, list)
	}
	return list
}
