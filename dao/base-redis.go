package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

var rdsConn *redis.Client

var userTokenKey = "userToken"

func InitRedis(Password, redisUrl string, DB int) *redis.Client { //InitTokenRedis
	rdsConn = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: Password, // no password set
		DB:       DB,  // use default DB
	})
	rdsConn.BgRewriteAOF()
	pong, err := rdsConn.Ping().Result()
	if err != nil{
		fmt.Println(pong, err)
		return nil
	}
	// Output: PONG <nil>
	return rdsConn
}


func SaveUserToken(userId , tokenStr string)error{
	//保存用户token与userId
	fmt.Println("set user conn ", userId)
	mashMember, err := json.Marshal(tokenStr)
	result, err := rdsConn.HSet(userTokenKey, string(userId), mashMember).Result()
	if err != nil{
		return err
	}
	fmt.Println("set result: ", result)
	return nil
}

func GetuserToken(userId string)(string, error){
	result, err := rdsConn.HGet(userTokenKey, userId).Result()
	if err != nil{
		return "", err
	}
	var token string
	json.Unmarshal([]byte(result), &token)
	fmt.Println("get user token ", userId, ":", token, len(token))
	return token, nil
}

func GetRedisClient()(*redis.Client, error){
	if rdsConn!=nil{
		return rdsConn, nil
	}
	return nil, errors.New("redis连接失败！")
}