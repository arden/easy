package main

import (
	"github.com/arden/gf-plus/example/gorm/repository"
	"github.com/gogf/gf/os/glog"
)

func main()  {
	user, err := repository.User.GetByPhone("13590309275")
	if err != nil {
		glog.Error(err)
		return
	}
	println(user.Uname)
}
