package conn_router

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

//业务id和 server的对于关系，分布式存储
//基于redis

var redisClient *redis.Client
var clientOnce = &sync.Once{}

func getClient() *redis.Client {

	clientOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "172.16.9.248:6379", //cclehui_test
			Password: "",                  // no password set
			DB:       0,                   // use default DB
		})

		//pong, err := client.Ping().Result()
		if _, err := redisClient.Ping().Result(); err != nil {
			panic(fmt.Sprintf("router getClient error:%v", err))
		}

	})

	return redisClient
}

func AddServerRoute(yewuId string, serverAddr string) error {
	client := getClient()

	err := client.Set(yewuId, serverAddr, 0).Err()

	return err
}

func GetServerAddr(yewuId string) (string, error) {
	client := getClient()

	result, err := client.Get(yewuId).Result()
	if err == nil {
		return result, nil
	}

	if err == redis.Nil {
		return "", nil
	} else {
		return "", errors.New(fmt.Sprintf("get serverAddr %s, error:%v", yewuId, err))
	}
}

func RemoveServerRoute(yewuId string) error {
	client := getClient()

	err := client.Del(yewuId).Err()

	return err
}
