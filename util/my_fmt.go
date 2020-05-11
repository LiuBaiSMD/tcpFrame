/*
auth:   wuxun
date:   2020-05-07 20:07
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package util

import (
	"bytes"
	"runtime"
)

func RunFuncName()string{
	pc := make([]uintptr,1)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}