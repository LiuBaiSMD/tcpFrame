/*
auth:   wuxun
date:   2020-05-07 20:07
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package util

import (
	"runtime"
)

func RunFuncName()string{
	pc := make([]uintptr,1)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
