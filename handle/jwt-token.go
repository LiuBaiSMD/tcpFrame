/*
auth:   wuxun
date:   2020-01-15 11:20
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package handle

import (
	heartbeat "tcpPractice/proto"
	"tcpPractice/proto"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
)

var SecretKey = "abcdefg"
type JwtTokenCreator struct {
}

type jwtCustomClaims struct {
	jwt.StandardClaims

	// 追加自己需要的信息
	Uid string `json:"uid"`
	Name string `json:"name"`
}

func (t *JwtTokenCreator) GetToken(ctx context.Context, req *heartbeat.TokenRequest, rsp *heartbeat.TokenResponse)error{
	log.Print("Received TokenCreator.TokenRequest request")
	fmt.Println(req)
	name := req.Name
	id := req.Uid
	tokenString, err := GetTokenReal(id, name)
	if err!=nil{
		return nil
	}
	rsp.Token = string(tokenString)

	return nil
}

func GetTokenReal(userId string, name string)(string, error){
	claims := &jwtCustomClaims{
		StandardClaims:jwt.StandardClaims{
		},
		Uid:userId,
		Name:name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err!=nil{
		return "", err
	}

	return string(tokenString), nil
}