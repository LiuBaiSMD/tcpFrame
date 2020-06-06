package dao

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"tcpFrame/util"
	"time"
)

var rdsConn *redis.Client

var userTokenKey = "userToken"

func InitRedis(Password, redisUrl string, DB int) *redis.Client { //InitTokenRedis
	rdsConn = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: Password, // no password set
		DB:       DB,       // use default DB
	})
	rdsConn.BgRewriteAOF()
	pong, err := rdsConn.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
		return nil
	}
	// Output: PONG <nil>
	return rdsConn
}

func buildTokenKey(userId string) string {
	return userTokenKey + ":" + userId
}

func SaveUserToken(userId, tokenStr string) error {
	//保存用户token与userId
	//mashMember, err := json.Marshal(tokenStr)
	result, err := rdsConn.Set(buildTokenKey(userId), tokenStr, time.Second * time.Duration(3600*2)).Result()
	if err != nil {
		fmt.Println("save user token error:", userId, err)
		return err
	}
	if result != "" {
		fmt.Println("save and update user token: ", userId, tokenStr)
	}

	//fmt.Println("set result: ", result, err)
	return nil
}

func GetuserToken(userId string) (string, error) {
	token, err := rdsConn.Get(buildTokenKey(userId)).Result()
	fmt.Println(util.RunFuncName(), userId, token, err)
	if err != nil {
		return "", err
	}
	fmt.Println("get user token ", userId, ":", token, len(token))
	return token, nil
}

func GetRedisClient() (*redis.Client, error) {
	if rdsConn != nil {
		return rdsConn, nil
	}
	return nil, errors.New("redis连接失败！")
}
