/*
auth:   wuxun
date:   2019-12-09 20:39
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package handle

import (
	"fmt"
	heartbeat "tcpPractice/proto"
	"tcpPractice/conns"
	"tcpPractice/dao"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
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

func Login(w http.ResponseWriter, r *http.Request) {
	//
	conn, err := upGrader.Upgrade(w, r, nil)
	var userIdItf interface{} = r.Header.Get("UserId")
	var userToken = r.Header.Get("UserToken")

	canLogin := checkToken(userIdItf.(string), userToken)
	if !canLogin{
		fmt.Printf("token login failed!")
		conn.Close()
	}
	userId, _ := strconv.Atoi(userIdItf.(string))
	fmt.Println("get login: ", userId, userToken)
	if err != nil {
		log.Printf("upgrade: %s", err)
		return
	}
	connID := conns.GetLastestConnID()   //获取最新的connID，进行连接排队
	connClient := conns.NewClient(userId, conn, connID)
	conns.PushChan(userId, connClient)
	ListenMessage(connClient)
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

func ListenMessage(c *conns.ClientConn){
	done := make(chan struct{})
	clientRes := heartbeat.Response{}

	msgReviecer := make(chan uint64, 1)
	msgRecData := make(chan string, 1)
	go func() {
		defer close(done)
		for {
			_, message, err := c.GetConn().ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if err := proto.Unmarshal(message, &clientRes); err != nil {
				log.Printf("proto unmarshal: %v", err)
			}
			log.Println("recv: ", clientRes.Data, clientRes.UserId)
			msgReviecer <- clientRes.UserId
			msgRecData <- clientRes.Data
		}
	}()

	go func(){
		for{
			receiver := <- msgReviecer
			d := <-msgRecData
			recConn := conns.GetConnByUId(10-int(receiver))
			if recConn == nil{
				continue
			}
			fmt.Println("sendId: ", int(receiver), "rec: ", 10-int(receiver), recConn.GetUserID())

			err := recConn.GetConn().WriteMessage(websocket.BinaryMessage, MsgAssemblerReader("recevied from: "+d, uint64(c.GetUserID())))
			if err != nil{
				log.Println("write close:", err)
			}
		}
	}()

}

func Chat2User(userId int, message string){

}