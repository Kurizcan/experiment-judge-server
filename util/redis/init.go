package redis

import (
	"github.com/go-redis/redis"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type ClientRedis struct {
	Object *redis.Client
}

var Client *ClientRedis

func new() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
	})
	if client != nil {
		log.Infof("redis client create success")
	}
	return client
}

func (c *ClientRedis) Init() {
	Client = &ClientRedis{Object: new()}
}

func (c *ClientRedis) Close() {
	Client.Object.Close()
}

func (c *ClientRedis) Del(key string) error {
	// 删除多个key, Del函数支持删除多个key
	err := Client.Object.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientRedis) HDel(key, filed string) error {
	// 删除一个字段id
	err := Client.Object.HDel(key, filed).Err()
	if err != nil {
		return err
	}
	return nil
}
