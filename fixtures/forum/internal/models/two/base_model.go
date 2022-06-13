package two

import (
	"encoding/json"
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/cache"
	. "github.com/bmbstack/ripple/fixtures/forum/internal/helper"
	. "github.com/bmbstack/ripple/helper"
	"gorm.io/gorm"
	"time"
)

var databaseAlias = "two"
var cacheAlias = "two"

type BaseModel struct {
	CreatedTimeStr string     `json:"createdTime" gorm:"-"`
	CreatedTime    *time.Time `json:"-" gorm:"column:created_time; type:datetime; not null; default:current_timestamp"`
	UpdatedTime    *time.Time `json:"-" gorm:"column:updated_time; type:datetime"`
	DeletedTime    *time.Time `json:"-" gorm:"column:deleted_time; type:datetime"`
	IsDeleted      int64      `json:"-" gorm:"column:is_deleted; type:tinyint(1); not null; default:0"`
}

func Orm() *ripple.Orm {
	return ripple.Default().GetOrm(databaseAlias)
}

func Cache() *cache.Cache {
	return ripple.Default().GetCache(cacheAlias)
}

func (this *BaseModel) AfterFind(*gorm.DB) error {
	if IsNotEmpty(this.CreatedTime) {
		this.CreatedTimeStr = this.CreatedTime.Format(DateFullLayout)
	}
	return nil
}

func (this *BaseModel) GetCache(cacheKey string) string {
	result, _ := Cache().Get(cacheKey)
	return result
}

func (this *BaseModel) SetCache(cacheKey string, data interface{}) {
	if IsNotEmpty(data) {
		bytes, err := json.Marshal(data)
		if err == nil {
			Cache().Set(cacheKey, string(bytes), CacheSeconds)
		}
	}
}

func (this *BaseModel) DeleteCache(key string) {
	Cache().Delete(key)
}

func (this *BaseModel) DeleteCacheByPrefix(prefix string) {
	Cache().DeleteByPrefix(prefix)
}
