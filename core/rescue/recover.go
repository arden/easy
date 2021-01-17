package rescue

import "github.com/gogf/gf/os/glog"

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		glog.Error(p)
	}
}
