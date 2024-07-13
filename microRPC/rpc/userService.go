package rpc

import "context"

type GetByUserIdReq struct {
	Id int
}
type GetByUserIdResp struct {
}

type UserService struct {
	// 用反射来赋值
	// 下面这个不叫方法，只能叫做类型是函数的字段，不是方法（它不是定义在UserService上的方法，本质上是一个字段）
	GetById func(ctx context.Context, req *GetByUserIdReq) (*GetByUserIdResp, error)
}

func (u UserService) Name() string {
	//TODO implement me
	return "user_service"
}
