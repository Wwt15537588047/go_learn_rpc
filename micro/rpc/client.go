package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/silenceper/pool"
	"go.learn.rpc/micro/rpc/message"
	"net"
	"reflect"
	"time"
)

type Client struct {
	// 使用连接池优化，替代之前的addr string，连接池中会使用到addr
	pool pool.Pool
}

func NewClient(addr string) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap:  1,
		MaxCap:      30,
		MaxIdle:     10,
		IdleTimeout: time.Minute,
		Factory: func() (interface{}, error) {
			// 创建一个新的TCP连接，连接超时时间为3s
			return net.DialTimeout("tcp", addr, 3*time.Second)
		},
		Close: func(i interface{}) error {
			// 对于连接池中的连接执行的关闭实现
			return i.(net.Conn).Close()
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		pool: p,
	}, nil
}

// InitClientProxy 要为GetById之类的函数类型的字段赋值
// InitClientProxy 的作用就是捕获本地调用，构建请求参数：服务名、方法名、调用参数，随后发起调用
func InitClientProxy(addr string, service Service) error {
	client, err := NewClient(addr)
	if err != nil {
		return err
	}
	return setFuncField(service, client)
}

func setFuncField(service Service, p Proxy) error {
	if service == nil {
		return errors.New("rpc : 不支持service为 nil")
	}
	val := reflect.ValueOf(service)
	typ := val.Type()
	// 只支持指向结构体的一级指针
	// Kind()返回typ变量的类型级别，若其为指针，则返回类型级别为指针
	// Elum()返回对应元素的指针解引用，type.Elem()返回指针解引用也即对应的具体类型值，.Kind()返回类型级别
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
				// GetById函数返回值，第一个值类型为Response，第二个值类型为error
				retVal := reflect.New(fieldTyp.Type.Out(0).Elem())
				// args[0] 是context
				ctx := args[0].Interface().(context.Context)
				// args[1] 是参数
				// 构建请求参数,如何获取请求名称，①获取类型名，但是类型名会冲突；②包名+类型名；
				// ③让所有调用实现一个接口，返回调用的名称,此时不需要关心命名空间
				reqData, err := json.Marshal(args[1].Interface())
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				req := &message.Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Data:        reqData,
				}
				req.CalculateHeaderLength()
				req.CalculateBodyLength()

				// 发起调用
				resp, err := p.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}

				var retErr error
				if len(resp.Error) > 0 {
					// 服务端传回来的error
					retErr = errors.New(string(resp.Error))
				}

				if len(resp.Data) > 0 {
					err = json.Unmarshal(resp.Data, retVal.Interface())
					if err != nil {
						//反序列化失败
						return []reflect.Value{retVal, reflect.ValueOf(err)}
					}
				}
				var retErrVal reflect.Value
				if retErr == nil {
					retErrVal = reflect.Zero(reflect.TypeOf(new(error)).Elem())
				} else {
					retErrVal = reflect.ValueOf(retErr)
				}
				return []reflect.Value{retVal, retErrVal}
			}
			// 设置值给GetById
			finVal := reflect.MakeFunc(fieldTyp.Type, fn)
			fieldVal.Set(finVal)
		}
	}
	return nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	// 使用Json序列化数据
	data := message.EnCodeReq(req)
	// 客户端发送请求
	resp, err := c.Send(data)
	if err != nil {
		return message.DeCodeResp(resp), err
	}

	return message.DeCodeResp(resp), nil
}

func (c *Client) Send(data []byte) ([]byte, error) {
	// 从连接池中获取连接
	val, err := c.pool.Get()
	if err != nil {
		return nil, err
	}
	conn := val.(net.Conn)
	defer func() {
		_ = conn.Close()
	}()
	//err = WriteMsg(conn, data)
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	return ReadMsg(conn)
}
