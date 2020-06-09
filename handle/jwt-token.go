/*
auth:   wuxun
date:   2020-01-15 11:20
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package handle

import (
	"github.com/dgrijalva/jwt-go"
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