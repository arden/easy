package plus

import (
	"fmt"
	gfplusGorm "github.com/arden/gf-plus/gorm"
	"gorm.io/gorm"
)

const (
	frameCoreComponentNameDatabase = "gfplus.core.component.database"
)

// Database returns an instance of database ORM object
// with specified configuration group name.
func GetGorm(name ...string) *gorm.DB {
	group := gfplusGorm.DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameDatabase, group)
	db := instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		gormDB := gfplusGorm.Instance(name...)
		return gormDB
	})
	if db != nil {
		return db.(*gorm.DB)
	}
	return nil
}