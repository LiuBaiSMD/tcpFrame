/*
auth:   wuxun
date:   2019-12-09 20:39
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package handle

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"tcpPractice/dao"
	heartbeat "tcpPractice/proto"
)

var upGrader = websocket.Upgrader{
	//对请求头进行检查
	//CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clientRes heartbeat.Request
	serverRsp heartbeat.Response
	msgSeqId uint64 = 0
	USERID uint64 = 666
	CLIENTID uint64 = 678

)

func UserAuth(w http.ResponseWriter, r *http.Request) {
	//用户通过账号密码获取token
	userId := r.Header.Get("Name")
	passwd := r.Header.Get("Passwd")
	fmt.Println("Name: ", userId, "passwd: ", passwd)
	tokenString, err := GetTokenReal(userId, userId)
	fmt.Println("getToken: ", tokenString)
	if err!=nil{
		log.Fatal(err)
	}
	err1 := dao.SaveUserToken(userId, tokenString)
	if err != nil {
		fmt.Fprint(w, err1)
	}
	fmt.Fprint(w, tokenString)
}

func checkToken(userId, userToken string)bool{
	redisUToken, _ := dao.GetuserToken(userId)
	fmt.Println("checkToken" , userToken, redisUToken, len(userToken), len(redisUToken))
	if redisUToken == userToken && len(userToken)>0{
		return true
	}
	return false
}

func MsgAssemblerReader(data string, userId uint64) []byte {
	msgSeqId += 1
	retPb := &heartbeat.Response{
		ClientId: CLIENTID,
		UserId:   userId,
		MsgId:    msgSeqId,
		SessionId: 1000,
		Data:     data,
	}
	byteData, err := proto.Marshal(retPb)
	if err != nil {
		log.Fatal("pb marshaling error: ", err)
	}
	return byteData
}
