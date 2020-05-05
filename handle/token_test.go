/*
auth:   wuxun
date:   2020-01-13 15:15
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/


package handle_test

import (
    "fmt"
    "testing"
    "tcpPractice/handle"
)


func Test_GetTokenReal(t *testing.T) {
    tokenString, err:=handle.GetTokenReal("12345", "12345")
    fmt.Println("get token: ", tokenString, err)
}
