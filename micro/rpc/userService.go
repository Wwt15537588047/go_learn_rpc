package rpc

import (
	"context"
	"log"
)

type GetByUserIdReq struct {
	Id int
}
type GetByUserIdResp struct {
	Msg string
}

type UserService struct {
	// 用反射来赋值
	// 下面这个不叫方法，只能叫做类型是函数的字段，不是方法（它不是定义在UserService上的方法，本质上是一个字段）
	GetById func(ctx context.Context, req *GetByUserIdReq) (*GetByUserIdResp, error)
}

func (u UserService) Name() string {
	return "user_service"
}

// 上面的 UserService只是一个定义，相当于注册中心的作用，下面的UserServiceServer是服务端的一个具体的实现
type UserServiceServer struct {
	Err error
	Msg string
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByUserIdReq) (*GetByUserIdResp, error) {
	log.Println(req)
	return &GetByUserIdResp{
		Msg: u.Msg,
	}, u.Err
}
func (u *UserServiceServer) Name() string {
	return "user_service"
}
