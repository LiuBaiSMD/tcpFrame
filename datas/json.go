// @Author: liubai
// @Date: 2020/5/2 10:18 下午
// @Desc: test transport struct data

package datas

type Request struct{
	Action string
	Name string
	PWD string
	UserId int
}

type Respone struct {
	Action string
	Name string
	PWD string
	Code int
	UserId int
}
