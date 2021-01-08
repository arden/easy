package redis

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/errors/gerror"
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