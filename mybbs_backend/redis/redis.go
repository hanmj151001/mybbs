package redis

// base km project
// ref: https://platgit.mihoyo.com/plat/go/operation/km/-/blob/74fff4c0707794f4f7d17d7a32c5dd054ac4b547/tools/redis.go

import (
	"collabtool/fbi/utils/merror"
	"context"
	"errors"
	"reflect"
	"time"

	"gopkg.mihoyo.com/takumi/log"

	rd "github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"gopkg.mihoyo.com/takumi/redis"
)

const (
	NodeName  = "fbi"
	PrefixStr = "fbi:" // 默认全局缓存前缀
)

var (
	Nil = redis.Nil
)

func GetClient() redis.Client {
	return redis.GetClient(NodeName)
}

// SetCache 缓存数据
func SetCache(ctx context.Context, key string, val interface{}, expire time.Duration, isMarshal ...bool) error {
	redisCli := redis.GetClient(NodeName)
	if len(isMarshal) == 0 || isMarshal[0] {
		var err error
		val, err = jsoniter.MarshalToString(val)
		if err != nil {
			return warpError(ctx, err)
		}
	}

	err := redisCli.Set(ctx, FormatKey(key), val, expire).Err()
	return warpError(ctx, err)
}

func GetCache(ctx context.Context, key string, val interface{}, isUnmarshal ...bool) error {
	redisCli := redis.GetClient(NodeName)
	res, err := redisCli.Get(ctx, FormatKey(key)).Result()
	if err != nil {
		if err != Nil {
			err = warpError(ctx, err)
		}

		return err
	}

	if len(isUnmarshal) == 0 || isUnmarshal[0] {
		return warpError(ctx, jsoniter.UnmarshalFromString(res, val))
	}

	if sPtr, ok := val.(*string); ok {
		*sPtr = res
	}

	return nil
}

// CacheGet 带缓存的get方法，优先从redis获取，失败后从f函数获取，然后更新缓存。注意：val必须是一个指针
// tips: 当f的返回值为nil的时候是不会缓存的
func CacheGet(ctx context.Context, key string, val interface{}, expires time.Duration,
	f func(ctx context.Context) (interface{}, error)) error {
	vf := reflect.ValueOf(val)
	if vf.Kind() != reflect.Ptr {
		return errors.New("val is not a ptr")
	}
	// 首先尝试从缓存读取数据
	err := GetCache(ctx, key, val)
	if err == nil {
		return nil
	}
	getVal, err := f(ctx)
	if err != nil {
		return err
	}
	if getVal == nil {
		return nil
	}
	vf.Elem().Set(reflect.ValueOf(getVal))
	// 加入缓存
	err = SetCache(ctx, key, getVal, expires)
	if err != nil {
		return err
	}
	return nil
}

func DelCache(ctx context.Context, key string) error {
	redisCli := redis.GetClient(NodeName)
	_, err := redisCli.Del(ctx, FormatKey(key)).Result()
	if err != nil {
		return warpError(ctx, err)
	}
	return nil
}

func Exists(ctx context.Context, key string) (bool, error) {
	redisCli := redis.GetClient(NodeName)
	res, err := redisCli.Exists(ctx, FormatKey(key)).Result()
	if err != nil {
		return false, warpError(ctx, err)
	}

	if res == 1 {
		return true, nil
	}

	return false, nil
}

// HDel hash表del
func HDel(ctx context.Context, key, field string) error {
	redisCli := redis.GetClient(NodeName)
	err := redisCli.HDel(ctx, FormatKey(key), field).Err()
	return warpError(ctx, err)
}

// HGet hash表get
func HGet(ctx context.Context, key, field string, val interface{}, isUnmarshal ...bool) error {
	redisCli := redis.GetClient(NodeName)
	res, err := redisCli.HGet(ctx, FormatKey(key), field).Result()
	if err != nil {
		if err != Nil {
			err = warpError(ctx, err)
		}

		return err
	}

	if len(isUnmarshal) == 0 || isUnmarshal[0] {
		return warpError(ctx, jsoniter.UnmarshalFromString(res, val))
	}

	if sPtr, ok := val.(*string); ok {
		*sPtr = res
	}

	return nil
}

// HSet hash表set
func HSet(ctx context.Context, key, field string, val interface{}, isMarshal ...bool) error {
	redisCli := redis.GetClient(NodeName)
	if len(isMarshal) == 0 || isMarshal[0] {
		var err error
		val, err = jsoniter.MarshalToString(val)
		if err != nil {
			return warpError(ctx, err)
		}
	}
	err := redisCli.HSet(ctx, FormatKey(key), field, val).Err()
	return warpError(ctx, err)
}

// HGetRaw hash表get
func HGetRaw(ctx context.Context, key, field string) string {
	redisCli := redis.GetClient(NodeName)
	results := redisCli.HGet(ctx, FormatKey(key), field).Val()
	return results
}

// HSetRaw hash表set
func HSetRaw(ctx context.Context, key, field string, value interface{}) error {
	redisCli := redis.GetClient(NodeName)
	err := redisCli.HSet(ctx, FormatKey(key), field, value).Err()
	return warpError(ctx, err)
}

// HGetAll hash表getAll
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	redisCli := redis.GetClient(NodeName)
	return redisCli.HGetAll(ctx, FormatKey(key)).Result()
}

// HMGet hash表getAll
func HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	redisCli := redis.GetClient(NodeName)
	res, err := redisCli.HMGet(ctx, FormatKey(key), fields...).Result()
	if err != nil {
		return nil, err
	}

	mRes := make([]interface{}, 0)
	for _, item := range res {
		if item != nil {
			mRes = append(mRes, item)
		}
	}

	return mRes, nil
}

// HMSet hash表批量添加
func HMSet(ctx context.Context, key string, data map[string]interface{}, isMarshal ...bool) error {
	redisCli := redis.GetClient(NodeName)
	if len(isMarshal) == 0 || isMarshal[0] {
		for k, v := range data {
			val, err := jsoniter.MarshalToString(v)
			if err != nil {
				return warpError(ctx, err)
			}

			data[k] = val
		}
	}

	err := redisCli.HMSet(ctx, FormatKey(key), data).Err()
	err = warpError(ctx, err)
	return err
}

// HKeys 获取keys
func HKeys(ctx context.Context, key string) ([]string, error) {
	redisCli := redis.GetClient(NodeName)
	return redisCli.HKeys(ctx, FormatKey(key)).Result()
}

// Expire
func Expire(ctx context.Context, key string, expire time.Duration) bool {
	redisCli := redis.GetClient(NodeName)
	b := redisCli.Expire(ctx, FormatKey(key), expire).Val()

	return b
}

func HSetAndExpire(ctx context.Context, key, account, value string, expire time.Duration) (err error) {
	redisCli := redis.GetClient(NodeName)

	formatKey := FormatKey(key)
	err = redisCli.HSet(ctx, formatKey, account, value).Err()
	if err != nil {
		return err
	}

	err = redisCli.Expire(ctx, formatKey, expire).Err()
	return err
}

func BRPop(ctx context.Context, key string, timeout time.Duration) (string, error) {
	redisCli := redis.GetClient(NodeName)

	formatKey := FormatKey(key)
	res, err := redisCli.BRPop(ctx, timeout, formatKey).Result()
	if err != nil {
		return "", err
	}

	resValue := ""
	if len(res) == 2 && res[0] == formatKey {
		resValue = res[1]
	}

	return resValue, nil
}

func RPop(ctx context.Context, key string) (string, error) {
	redisCli := redis.GetClient(NodeName)

	formatKey := FormatKey(key)
	res, err := redisCli.RPop(ctx, formatKey).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func LPush(ctx context.Context, key, value string) error {
	if value == "" { // 避免无意义的空串
		return nil
	}

	redisCli := redis.GetClient(NodeName)
	formatKey := FormatKey(key)

	_, err := redisCli.LPush(ctx, formatKey, value).Result()

	return err
}

func LPushN(ctx context.Context, key string, values ...string) error {
	if len(values) == 0 { // 避免空数组，造成：wrong number of arguments for 'lpush' command 错误
		return nil
	}

	redisCli := redis.GetClient(NodeName)
	formatKey := FormatKey(key)

	_, err := redisCli.LPush(ctx, formatKey, values).Result()

	return err
}

func RPush(ctx context.Context, key, value string) error {
	if value == "" { // 避免无意义的空串
		return nil
	}

	redisCli := redis.GetClient(NodeName)
	formatKey := FormatKey(key)

	_, err := redisCli.RPush(ctx, formatKey, value).Result()
	return err
}

func ZAdd(ctx context.Context, key string, score float64, member interface{}) error {
	redisCli := redis.GetClient(NodeName)
	formatKey := FormatKey(key)

	m := rd.Z{
		Score:  score,
		Member: member,
	}

	_, err := redisCli.ZAdd(ctx, formatKey, m).Result()

	return err
}

func ZRangeByScore(ctx context.Context, key string, min, max string, offset, count int64) ([]string, error) {
	redisCli := redis.GetClient(NodeName)
	formatKey := FormatKey(key)

	opt := rd.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}

	res, err := redisCli.ZRangeByScore(ctx, formatKey, opt).Result()

	return res, err
}

func ZRemRangeByScore(ctx context.Context, key, min, max string) error {
	redisCli := redis.GetClient(NodeName)
	formatKey := FormatKey(key)

	_, err := redisCli.ZRemRangeByScore(ctx, formatKey, min, max).Result()

	return err
}

func ZRem(ctx context.Context, key string, member interface{}) error {
	redisCli := redis.GetClient(NodeName)
	formatKey := FormatKey(key)

	_, err := redisCli.ZRem(ctx, formatKey, member).Result()

	return err
}

// Incr 自增
func Incr(ctx context.Context, key string) (result int64, err error) {
	redisCli := redis.GetClient(NodeName)

	result, err = redisCli.Incr(ctx, FormatKey(key)).Result()
	return result, warpError(ctx, err)
}

// FormatKey 格式化key
func FormatKey(key string) string {
	return PrefixStr + key
}

func warpError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	log.Stack(ctx, "redis error",
		merror.ErrorField(err),
	)

	return err
}
