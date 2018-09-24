package two

import (
	"time"
	. "github.com/bmbstack/ripple/fixtures/forum/helper"
	"github.com/bmbstack/ripple"
	. "github.com/bmbstack/ripple/helper"
	"github.com/bmbstack/ripple/middleware/cache"
	"encoding/json"
	"github.com/labstack/echo"
)

var Orm *ripple.Orm
var databaseAlias = "two"
var cacheAlias = "two"

type BaseModel struct {
	CreatedTimeStr string     `json:"createdTime" gorm:"-"`
	CreatedTime    *time.Time `json:"-" gorm:"column:created_time; type:datetime; not null; default:current_timestamp"`
	UpdatedTime    *time.Time `json:"-" gorm:"column:updated_time; type:datetime"`
	DeletedTime    *time.Time `json:"-" gorm:"column:deleted_time; type:datetime"`
	IsDeleted      int64      `json:"-" gorm:"column:is_deleted; type:tinyint(1); not null; default:0"`
}

func init() {
	Orm = ripple.GetOrm(databaseAlias)
	Orm.DB.LogMode(ripple.GetConfig().DebugOn)
}

func (this *BaseModel) AfterFind() {
	this.CreatedTimeStr = this.CreatedTime.Format(DateShortLayout)
}

func (this *BaseModel) GetCache(ctx echo.Context, cacheKey string) string {
	return cache.Store(cacheAlias, ctx).Get(cacheKey)
}

func (this *BaseModel) SetCache(ctx echo.Context, cacheKey string, data interface{}) {
	if IsNotEmpty(data) {
		bytes, err := json.Marshal(data)
		if err == nil {
			cache.Store(cacheAlias, ctx).Put(cacheKey, string(bytes), CacheSeconds)
		}
	}
}
