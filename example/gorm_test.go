package example

import (
	"github.com/arden/gf-plus/example/gorm/service"
	"github.com/gogf/gf/os/glog"
	"testing"
)

func TestGorm(t *testing.T) {
	user, err := service.User.GetByPhone("13590309275")
	if err != nil {
		glog.Error(err)
		return
	}
	println(user.Uname)
}
