package database

import (
	"fmt"
	"github.com/arden/easy/app"
	easyRedis "github.com/arden/easy/database/redis"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

const (
	frameCoreComponentNameRedis = "easy.core.component.redis"
	configNodeNameRedis         = "redis"
)

// Redis returns an instance of redis client with specified configuration group name.
func GetRedis(name ...string) *easyRedis.Redis {
	config := gins.Config()
	group := easyRedis.DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameRedis, group)
	result := app.Instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		// If already configured, it returns the redis instance.
		if _, ok := easyRedis.GetConfig(group); ok {
			return easyRedis.Instance(group)
		}
		// Or else, it parses the default configuration file and returns a new redis instance.
		var m map[string]interface{}
		if _, v := gutil.MapPossibleItemByKey(gins.Config().GetMap("."), configNodeNameRedis); v != nil {
			m = gconv.Map(v)
		}
		if len(m) > 0 {
			if v, ok := m[group]; ok {
				redisConfig, err := easyRedis.ConfigFromStr(gconv.String(v))
				if err != nil {
					panic(err)
				}
				return easyRedis.New(redisConfig)
			} else {
				panic(fmt.Sprintf(`configuration for redis not found for group "%s"`, group))
			}
		} else {
			filePath, _ := config.GetFilePath()
			panic(fmt.Sprintf(`incomplete configuration for redis: "redis" node not found in config file "%s"`, filePath))
		}
		return nil
	})
	if result != nil {
		return result.(*easyRedis.Redis)
	}
	return nil
}
