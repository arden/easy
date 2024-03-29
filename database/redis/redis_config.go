package redis

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"time"
)

const (
	DefaultGroupName = "default" // Default configuration group name.
	DefaultRedisPort = 6379      // Default redis port configuration if not passed.
	DefaultMaxActive = 50        // Default redis max active connections.
)

var (
	// Configuration groups.
	configs = gmap.NewStrAnyMap(true)
)

// SetConfig sets the global configuration for specified group.
// If <name> is not passed, it sets configuration for the default group name.
func SetConfig(config Config, name ...string) {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	configs.Set(group, config)
	instances.Remove(group)

	glog.Infof(`SetConfig for group "%s": %+v`, group, config)
}

// SetConfigByStr sets the global configuration for specified group with string.
// If <name> is not passed, it sets configuration for the default group name.
func SetConfigByStr(str string, name ...string) error {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	config, err := ConfigFromStr(str)
	if err != nil {
		return err
	}
	configs.Set(group, config)
	instances.Remove(group)
	return nil
}

// GetConfig returns the global configuration with specified group name.
// If <name> is not passed, it returns configuration of the default group name.
func GetConfig(name ...string) (config Config, ok bool) {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	if v := configs.Get(group); v != nil {
		return v.(Config), true
	}
	return Config{}, false
}

// RemoveConfig removes the global configuration with specified group.
// If <name> is not passed, it removes configuration of the default group name.
func RemoveConfig(name ...string) {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	configs.Remove(group)
	instances.Remove(group)

	glog.Infof(`RemoveConfig: %s`, group)
}

// ConfigFromStr parses and returns config from given str.
// Eg: host:port[,db,pass?maxIdle=x&maxActive=x&idleTimeout=x&maxConnLifetime=x]
func ConfigFromStr(str string) (config Config, err error) {
	array, _ := gregex.MatchString(`([^:]+):*(\d*),{0,1}(\d*),{0,1}(.*)\?(.+?)`, str)
	if len(array) == 6 {
		parse, _ := gstr.Parse(array[5])
		config = Config{
			Host: array[1],
			Port: gconv.Int(array[2]),
			Db:   gconv.Int(array[3]),
			Pass: array[4],
		}
		if config.Port == 0 {
			config.Port = DefaultRedisPort
		}
		if v, ok := parse["maxIdle"]; ok {
			config.MaxIdle = gconv.Int(v)
		}
		if v, ok := parse["maxActive"]; ok {
			config.MaxActive = gconv.Int(v)
		} else {
			config.MaxActive = DefaultMaxActive
		}
		if v, ok := parse["idleTimeout"]; ok {
			config.IdleTimeout = gconv.Duration(v) * time.Second
		}
		if v, ok := parse["maxConnLifetime"]; ok {
			config.MaxConnLifetime = gconv.Duration(v) * time.Second
		}
		if v, ok := parse["tls"]; ok {
			config.TLS = gconv.Bool(v)
		}
		if v, ok := parse["skipVerify"]; ok {
			config.TLSSkipVerify = gconv.Bool(v)
		}
		return
	}
	array, _ = gregex.MatchString(`([^:]+):*(\d*),{0,1}(\d*),{0,1}(.*)`, str)
	if len(array) == 5 {
		config = Config{
			Host: array[1],
			Port: gconv.Int(array[2]),
			Db:   gconv.Int(array[3]),
			Pass: array[4],
		}
		if config.Port == 0 {
			config.Port = DefaultRedisPort
		}
	} else {
		err = gerror.Newf(`invalid redis configuration: "%s"`, str)
	}
	return
}

// ClearConfig removes all configurations and instances of redis.
func ClearConfig() {
	configs.Clear()
	instances.Clear()
}
