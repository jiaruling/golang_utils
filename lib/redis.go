package lib

import "github.com/go-redis/redis/v8"

// https://redis.uptrace.dev/
// https://github.com/go-redis/redis

var redisMap map[string]*redis.Client

type Redis struct {
	Name     string
	Addr     string
	Password string
	DB       int
}

func NewRedis(addr, password string, db int) *Redis {
	if redisMap == nil {
		redisMap = make(map[string]*redis.Client)
	}
	return &Redis{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
}

func (r *Redis) InitRedis() {
	if redisMap == nil {
		redisMap = make(map[string]*redis.Client)
	}
	if r.Name == "" {
		r.Name = "default"
	}
	if _, ok := redisMap[r.Name]; ok {
		return
	}
	redisMap[r.Name] = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password, // no password set
		DB:       r.DB,       // use default DB
	})
}

func GetRedis(name... string) *redis.Client {
	if redisMap == nil {
		return nil
	}
	var n string
	if len(name) > 0 {
		n = name[0]
	} else {
		n = "default"
	}
	if _, ok := redisMap[n]; ok {
		return redisMap[n]
	} else {
		return nil
	}
}

func DestroyRedis(name... string) {
	if redisMap == nil {
		return
	}
	var n string
	if len(name) > 0 {
		n = name[0]
	} else {
		n = "default"
	}
	if _, ok := redisMap[n]; ok {
		_ = redisMap[n].Close()
		delete(redisMap, n)
	}
}

func DestroyRedisAll() {
	if redisMap == nil {
		return
	}
	for k, r := range redisMap {
		_ = r.Close()
		delete(redisMap, k)
	}
}
