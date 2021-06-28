package model

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

type RedisConf struct {
	Host   string `mapstructure:"host" json:"host" yaml:"host"`
	Port   string `mapstructure:"port" json:"port" yaml:"port"`
	Pwd    string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`
	DBName int    `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Head   string `mapstructure:"head" json:"head" yaml:"head"`
}

var Pool *redis.Client

var ctx = context.Background()

//初始化
func init() {
	Pool = InitRedis(ServerInfo.RedisConf)
}

// redis初始化
func InitRedis(redisMsg RedisConf) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     redisMsg.Host + ":" + redisMsg.Port,
		Password: redisMsg.Pwd,
		DB:       redisMsg.DBName,
	})
	err := client.Ping().Err()
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	Pool = client
	fmt.Println("连接 redis 成功")
	return client
}

//清理缓存
func Delcash() error {

	//拿到key头在集合中的数量
	num, err := Pool.SCard(ServerInfo.RedisConf.Head).Result()
	if err != nil {
		return err
	}
	var i int64
	for i = 0; i < num; i++ {

		//删除一条数据返回被删除的元素，逐个删除，但这个会返回对应元素
		red_key, err := Pool.SPop(ServerInfo.RedisConf.Head).Result()

		if err != nil {
			return err
		}

		if Pool.Del(red_key).Err() != nil {
			return err
		}

	}

	return err
}
