package fsyn

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type Redis struct {
	redisobj redis.Conn
}

func (redisclient *Redis) HgetAll(key string) []string {
	v, err := redis.Strings(redisclient.redisobj.Do("HGETALL", key))

	var record = make([]string, len(v)/2)

	if err != nil {
		fmt.Println(err)
		return v
	}

	j := 0

	for i := 0; i < len(v); i++ {
		if (i % 2) == 0 {
			record[j] = v[i]
			j++
		}
	}

	return record
}

func (redisclient *Redis) Hdel(key string, field string) int64 {
	result, err := redis.Int64(redisclient.redisobj.Do("HDEL", key, field))
	if err != nil {
		fmt.Println(err)
		return result
	}

	return result
}

func (redisclient *Redis) Hget(key string, field string) string {
	v, err := redis.String(redisclient.redisobj.Do("HGET", key, field))
	if err != nil {
		fmt.Println(err)
		return v
	}
	return v
}

func (redisclient *Redis) Hset(key string, field string, value string) string {
	result, err := redis.String(redisclient.redisobj.Do("HMSET", key, field, value))
	if err != nil {
		return result
	}
	return result
}

func NewRedis(ip string, port int) (redisclient *Redis, err error) {
	redisclient = &Redis{}
	var addr = ip + ":" + strconv.Itoa(port)
	redisclient.redisobj, err = redis.DialTimeout("tcp", addr, 0, 1*time.Second, 1*time.Second)
	return redisclient, err
}
