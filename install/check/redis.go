package check

import (
	"fmt"
	"gatekeeper/install/tool"
	"github.com/go-redis/redis"
)


type Redis struct{
	Host 	 string
	Port 	 string
	Pwd	 	 string
}


var (
	RedisClient Redis
)


func InitRedis() error{
	host, err := tool.Input("please enter redis host (default:127.0.0.1)", "127.0.0.1")
	if err != nil{
		return err
	}

	port, err := tool.Input("please enter redis port (default:6379)", "6379")
	if err != nil{
		return err
	}

	pwd, err := tool.Input("please enter redis pwd (default:null)", "")
	if err != nil{
		return err
	}

	redisClient := Redis{
		Host: host,
		Port: port,
		Pwd: pwd,
	}
	RedisClient = redisClient
	tool.LogInfo.Println(fmt.Sprintf("redis connect info host:[%s] port:[%s] pwd:[%s]", host, port, pwd))
	err = redisClient.Init();if err !=nil{
		tool.LogError.Println(err)
		return err
	}
	return nil
}

func (r *Redis) Init() error{
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Pwd,
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		tool.LogWarning.Println(err)
		return InitRedis()
	}
	tool.LogInfo.Println("connect redis success")
	return nil
}