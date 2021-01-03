package main

import (
	"github.com/arden/gf-plus/example/gorm/service"
	"github.com/gogf/gf/os/glog"
)

func main()  {
	user, err := service.User.GetByPhone("13590309275")
	if err != nil {
		glog.Error(err)
		return
	}
	println(user.Uname)
}
