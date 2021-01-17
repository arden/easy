package database

import (
	"fmt"
	"github.com/arden/easy/app"
	easyGorm "github.com/arden/easy/database/gorm"
	"gorm.io/gorm"
)

const (
	frameCoreComponentNameDatabase = "easy.core.component.database"
)

// Database returns an instance of database ORM object
// with specified configuration group name.
func GetGorm(name ...string) *gorm.DB {
	group := easyGorm.DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameDatabase, group)
	db := app.Instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		gormDB := easyGorm.Instance(name...)
		return gormDB
	})
	if db != nil {
		return db.(*gorm.DB)
	}
	return nil
}