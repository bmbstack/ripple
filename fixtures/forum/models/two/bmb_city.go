package two

import (
	"encoding/json"
	"fmt"
	. "github.com/bmbstack/ripple/helper"
	"time"
)

type BmbCity struct {
	ID   int64  `json:"ID" gorm:"column:id; type:int(11); primary_key; auto_increment"`
	Name string `json:"name" gorm:"column:name; type:varchar(50); not null; default:'北京'"`
	Pid  int64  `json:"pid" gorm:"column:pid; type:tinyint(4); not null; default:1"`
	BaseModel
}

func (BmbCity) TableName() string {
	return "bmb_city"
}

func (this *BmbCity) Save() bool {
	this.DeleteCacheByPrefix("FindCity")
	nowTime := time.Now()
	if IsNotEmpty(this.ID) {
		location, _ := time.LoadLocation(TimeLocationName)
		createdDate, _ := time.ParseInLocation(DateFullLayout, this.CreatedTimeStr, location)
		this.CreatedTime = &createdDate

		this.UpdatedTime = &nowTime
	}
	if this.IsDeleted == 1 {
		this.DeletedTime = &nowTime
	}
	return Orm().DB.Save(this).RowsAffected == 1
}

func (this *BmbCity) FindCityListByPid(pid int64) (list []BmbCity) {
	cacheKey := fmt.Sprintf("FindCityListByPid:%d", pid)
	cacheValue := this.GetCache(cacheKey)
	if IsNotEmpty(cacheValue) {
		_ = json.Unmarshal([]byte(cacheValue), &list)
	} else {
		Orm().DB.Where("pid=? AND is_deleted=0", pid).Find(&list)
		this.SetCache(cacheKey, list)
	}
	return list
}

func (this *BmbCity) FindCityID(id int64) (one BmbCity) {
	cacheKey := fmt.Sprintf("%s:%d", CurrentMethodName(), id)
	cacheValue := this.GetCache(cacheKey)
	if IsNotEmpty(cacheValue) {
		_ = json.Unmarshal([]byte(cacheValue), &one)
	} else {
		Orm().DB.Where("id=? AND is_deleted=0", id).Find(&one)
		this.SetCache(cacheKey, one)
	}
	return one
}
