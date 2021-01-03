package main

import (
	"github.com/arden/gf-plus/plus"
	"github.com/gogf/gf/os/glog"
)

func main()  {
	var user User
	gormDB := plus.Gorm()
	db := gormDB.First(&user,23166898)
	if db.Error != nil {
		glog.Error(db.Error.Error())
		return
	}
	println(user.Uname)

	var user2 User
	gormDB2 := plus.Gorm()
	db = gormDB2.First(&user2,  1120442)
	if db.Error != nil {
		glog.Error(db.Error.Error())
		return
	}
	println(user2.Uname)
}
