package rpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// InitClientProxy 要为GetById之类的函数类型的字段赋值
func InitClientProxy(service Service) error {
	return setFuncField(service, nil)
}

func setFuncField(service Service, p Proxy) error {
	if service == nil {
		return errors.New("rpc : 不支持service为 nil")
	}
	val := reflect.ValueOf(service)
	typ := val.Type()
	// 只支持指向结构体的一级指针
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("rpc : 只支持指向结构体的一级指针")
	}

	val = val.Elem()
	typ = typ.Elem()
	for i := 0; i < typ.NumField(); i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		if fieldVal.CanSet() {
			// 这里才是真正的将本地调用捕捉到的地方
			fn := func(args []reflect.Value) (results []reflect.Value) {
				// args[0] 是context
				ctx := args[0].Interface().(context.Context)
				// args[1] 是参数
				// 构建请求参数,如果获取请求名称，①获取类型名，但是类型名会冲突；②包名+类型名；
				// ③让所有调用实现一个接口，返回调用的名称,此时不需要关心命名空间
				req := &Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Args:        args[1].Interface(),
				}

				// 发起调用
				resp, err := p.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(err)}
				}

				fmt.Println(resp)
				return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(errors.New("nil"))}
			}
			// 设置值给GetById
			finVal := reflect.MakeFunc(fieldTyp.Type, fn)
			fieldVal.Set(finVal)
		}
	}
	return nil
}
