package example

import (
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestGorm(t *testing.T) {
	g.Config().SetFileName("config.toml")
	value := g.Config().GetString("redis")
	println(value)
	//user, err := service.User.GetByPhone("13590309275")
	//if err != nil {
	//	glog.Error(err)
	//	return
	//}
	//println(user.Uname)
}
