package redis

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gconv"
	"time"
)

// GetObject get object from key.
// It returns an error if occurs.
func (db *Redis) GetObject(ctx context.Context, key string, value interface{}) error {
	cmd := db.client.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}
	val := cmd.Val()
	err := json.Unmarshal([]byte(val), value)
	if err != nil {
		return err
	}
	return nil
}

// AddObject add object for a key.
// It returns an error if occurs.
func (db *Redis) AddObject(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := db.client.Set(ctx, key, j, duration)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

// DeleteObject delete object from key.
// It returns an error if occurs.
func (db *Redis) DeleteObject(ctx context.Context, key string) error {
	cmd := db.client.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

// GetAllHMObjects get all objects from a hash map table.
// It returns the objects and an error if occurs.
func (db *Redis) GetAllHMObjects(ctx context.Context, entity string) (map[string]string, error) {
	cmd := db.client.HGetAll(ctx, entity)
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	return cmd.Val(), nil
}

// GetHMObject get an object from a hash map table.
// It returns an error if occurs.
func (db *Redis) GetHMObject(ctx context.Context, entity, key string, value interface{}) error {
	cmd := db.client.HMGet(ctx, entity, key)
	if err := cmd.Err(); err != nil {
		return err
	}
	val, ok := cmd.Val()[0].(string)
	if !ok {
		return gerror.Newf("not found for key: %v", key)
	}
	err := json.Unmarshal([]byte(val), value)
	if err != nil {
		return err
	}
	return nil
}

// AddHMObject add an object to a hash map table.
// It returns an error if occurs.
func (db *Redis) AddHMObject(ctx context.Context, entity, key string, value interface{}) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := db.client.HMSet(ctx, entity, map[string]interface{}{key: j})
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

// DeleteHMObject delete an object from a hash map table.
// It returns an error if occurs.
func (db *Redis) DeleteHMObject(ctx context.Context, entity, key string) error {
	cmd := db.client.HDel(ctx, entity, key)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

// IsReady verify the database is ready.
// It returns true if ready.
func (db *Redis) IsReady(ctx context.Context) bool {
	if db.client == nil {
		return false
	}
	return db.client.Ping(ctx).Err() == nil
}

/*
ScanIter 获取匹配指定pattern的所有redis key  替代keys方法
*/
func (db *Redis) ScanIter(pattern string, count int64) *ScansCmd {
	scansCmd := NewScansCmd()
	cursor := uint64(0)
	for {
		scanCmd := db.Scan(context.Background(), cursor, pattern, count)
		scansCmd.addScanCmd(scanCmd)
		_, cursor, err := scanCmd.Result()
		if err != nil || cursor == 0 {
			return scansCmd
		}
	}
}

/*
SScanIter 获取集合中匹配指定pattern的所有元素  替代sismembers
*/
func (db *Redis) SScanIter(key, match string, count int64) *ScansCmd {
	scansCmd := NewScansCmd()
	cursor := uint64(0)
	for {
		scanCmd := db.SScan(context.Background(), key, cursor, match, count)
		scansCmd.addScanCmd(scanCmd)
		_, cursor, err := scanCmd.Result()
		if err != nil || cursor == 0 {
			return scansCmd
		}
	}
}

/*
ZScanIter 获取有序集合中匹配指定pattern的所有元素
*/
func (db *Redis) ZScanIter(key, match string, count int64) *ScansCmd {
	scansCmd := NewScansCmd()
	cursor := uint64(0)
	for {
		scanCmd := db.ZScan(context.Background(), key, cursor, match, count)
		scansCmd.addScanCmd(scanCmd)
		_, cursor, err := scanCmd.Result()
		if err != nil || cursor == 0 {
			return scansCmd
		}
	}
}

/*
HScanIter 获取字典中匹配指定pattern的所有字段
*/
func (db *Redis) HScanIter(key, match string, count int64) *ScansCmd {
	scansCmd := NewScansCmd()
	cursor := uint64(0)
	for {
		scanCmd := db.HScan(context.Background(), key, cursor, match, count)
		scansCmd.addScanCmd(scanCmd)
		_, cursor, err := scanCmd.Result()
		if err != nil || cursor == 0 {
			return scansCmd
		}
	}
}

/*
HSetEX 执行hset命令并设置过期时间 单位: 秒
*/
func (db *Redis) HSetEX(key, field string, value interface{}, expiration time.Duration) *Cmd {
	keys := []string{key}
	ex := durationToIntSecond(expiration)
	return db.Eval(context.Background(), hsetScript, keys, ex, field, value)
}

/*
HMSetEX 执行hmset命令并设置过期时间 单位: 秒
*/
func (db *Redis) HMSetEX(key string, fields map[string]interface{}, expiration time.Duration) *Cmd {
	keys := []string{key}
	ex := durationToIntSecond(expiration)
	args := []interface{}{ex}
	for k, v := range fields {
		if v != nil {
			args = append(args, k, v)
		}
	}
	return db.Eval(context.Background(), hmsetScript, keys, args...)
}

/*
ZAddRemByRank 向zset中插入成员并剪切，并截取只保留分数最高的length个成员
*/
func (db *Redis) ZAddRemByRank(key string, length int64, members ...Z) *Cmd {
	keys := []string{key}
	args := []interface{}{0, -(length + 1)}
	for _, member := range members {
		args = append(args, member.Score, member.Member)
	}
	return db.Eval(context.Background(), zaddRemByRankScript, keys, args...)
}

/*
LPushTrim 从左边向list插入元素，并截取只保留左起length个元素
*/
func (db *Redis) LPushTrim(key string, length int64, values ...interface{}) *Cmd {
	keys := []string{key}
	args := []interface{}{0, length - 1}
	args = append(args, values...)
	return db.Eval(context.Background(), lpushTrimScript, keys, args...)
}

/*
RPushTrim 从右边向list插入元素，并截取只保留右起length个元素
*/
func (db *Redis) RPushTrim(key string, length int64, values ...interface{}) *Cmd {
	keys := []string{key}
	args := []interface{}{-length, -1}
	args = append(args, values...)
	return db.Eval(context.Background(), rpushTrimScript, keys, args...)
}

/*
DurationToIntSecond 将time.Duration类型转换为值为秒数的int类型
*/
func durationToIntSecond(duration time.Duration) int {
	return int(duration) / 1e9
}

func (db *Redis) HMSetMapInfo(ctx context.Context, key string, values map[string]interface{}) *BoolCmd {

	for k , v := range values {
		values[k] = gconv.String(v)
	}

	return db.client.HMSet(context.Background(), key, values)
}