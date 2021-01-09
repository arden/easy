package redis

import "github.com/go-redis/redis/v8"

// 起别名
type (
	// Options  实例化redis client时的参数结构类型
	Options     = redis.Options
	// BoolCmd bool
	BoolCmd     = redis.BoolCmd
	// Cmd cmd
	Cmd         = redis.Cmd
	// Z z
	Z           = redis.Z
	// ScanCmd  scan cmd
	ScanCmd     = redis.ScanCmd
)

/*
ScansCmd 定义ScansCmd结构类型, 用于scaniter方法的返回值类型
*/
type ScansCmd struct {
	scanCmds []*ScanCmd
}

/*
Val 获取命令执行的结果
*/
func (c *ScansCmd) Val() (keys []string) {
	for _, scanCmd := range c.scanCmds {
		page, _ := scanCmd.Val()
		keys = append(keys, page...)
	}
	return keys
}

/*
Result 获取命令执行的结果
*/
func (c *ScansCmd) Result() (keys []string, err error) {
	for _, scanCmd := range c.scanCmds {
		page, _, e := scanCmd.Result()
		if e != nil {
			err = e
		}
		keys = append(keys, page...)
	}
	return keys, err
}

/*
String 以字符串形式展示每个命令的执行结果
*/
func (c *ScansCmd) String() (strSlice []string) {
	for _, scanCmd := range c.scanCmds {
		strSlice = append(strSlice, scanCmd.String())
	}
	return strSlice
}

func (c *ScansCmd) addScanCmd(scanCmd *ScanCmd) {
	c.scanCmds = append(c.scanCmds, scanCmd)
}

/*
NewScansCmd 实例化
*/
func NewScansCmd() *ScansCmd {
	return &ScansCmd{}
}
