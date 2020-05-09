// @Author: liubai
// @Date: 2020/5/2 10:18 下午
// @Desc: test transport struct data

package datas

type BaseData struct {
	Action string
	UserId int
	BData []byte
}

type Request struct{
	Action string
	Name string
	PWD string
	UserId int
	BData []byte
}

type Respone struct {
	Action string
	Name string
	PWD string
	Code int
	UserId int
	BData []byte
}

type LoginRequest struct{
	Action string
	Name string
	PWD string
	UserId int
	BData []byte
}

type LoginRespone struct {
	Action string
	Name string
	PWD string
	Code int
	UserId int
	BData []byte
}

