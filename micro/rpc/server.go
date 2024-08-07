package rpc

import (
	"context"
	"errors"
	"go.learn.rpc/micro/rpc/message"
	"go.learn.rpc/micro/rpc/serialize"
	"go.learn.rpc/micro/rpc/serialize/json"
	"net"
	"reflect"
)

type Server struct {
	services    map[string]reflectionStub
	serializers map[uint8]serialize.Serializer
}

func NewServer1() *Server {
	ser := &Server{
		services:    make(map[string]reflectionStub, 16),
		serializers: make(map[uint8]serialize.Serializer, 4),
	}
	// 注册一个默认的序列化协议
	ser.RegisterSerializer(&json.Serializer{})
	return ser
}
func (s *Server) RegisterSerializer(serializer serialize.Serializer) {
	s.serializers[serializer.Code()] = serializer
}

func (s *Server) RegisterServer(service Service) {
	s.services[service.Name()] = reflectionStub{
		s:           service,
		value:       reflect.ValueOf(service),
		serializers: s.serializers,
	}
}

func (s *Server) Start(network, addr string) error {
	// Listen的第一个参数network规定通信的协议，tcp还是udp,第二个参数addr里面包含地址和端口
	listener, err := net.Listen(network, addr)
	defer listener.Close()
	if err != nil {
		// 比较常见的就是端口占用
		return err
	}

	for {
		// 使用for循环不断监听，如果有连接到来，则开启一个协程对连接进行处理
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if err = s.handleConn(conn); err != nil {
				conn.Close()
			}
		}()
	}
}

// 我们规定，一个请求包含两部分
// 1.长度字段，用8个字节表示
// 2.请求数据
// 响应数据也是这样规定的
func (s *Server) handleConn(conn net.Conn) error {
	for {
		reqBs, err := ReadMsg(conn)
		if err != nil {
			// 这里是连接err
			return err
		}
		// 获取请求数据
		req := message.DeCodeReq(reqBs)
		if err != nil {
			return err
		}

		// context.BackGround用于创建一个新的上下文
		resp, err := s.Invoke(context.Background(), req)
		if err != nil {
			// 此处是自己的业务逻辑，正常处理逻辑应该是将业务逻辑封装随后返回给调用端,将错误封装在resp.Error里
			resp.Error = []byte(err.Error())
		}
		resp.CalculateHeaderLength()
		resp.CalculateBodyLength()
		_, err = conn.Write(message.EnCodeResp(resp))
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *Server) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	// 还原业务调用，此时已经拿到了service name,method name 和参数了
	service, ok := s.services[req.ServiceName]
	resp := &message.Response{
		MessageId:  req.MessageId,
		Version:    req.Version,
		Compresser: req.Compresser,
		Serializer: req.Serializer,
	}
	if !ok {
		return resp, errors.New("你所调用的服务不存在")
	}
	respData, err := service.invoke(ctx, req)
	resp.Data = respData
	if err != nil {
		return resp, err
	}
	return resp, nil
}

type reflectionStub struct {
	s           Service
	value       reflect.Value
	serializers map[uint8]serialize.Serializer
}

// reflectionStub相当于一个桩，找到一个桩，随后在桩里面解决反射的问题
func (s *reflectionStub) invoke(ctx context.Context, req *message.Request) ([]byte, error) {
	// 反射执行方法，并且执行
	method := s.value.MethodByName(req.MethodName)
	in := make([]reflect.Value, 2)
	// 暂时没有传下标0对应的参数，直接写死
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	// 拿取规定的序列化协议
	serializer, ok := s.serializers[req.Serializer]
	if !ok {
		return nil, errors.New("Server 不支持的序列化协议.")
	}
	err := serializer.DeCode(req.Data, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq
	results := method.Call(in)
	if results[1].Interface() != nil {
		err = results[1].Interface().(error)
	}
	var res []byte
	if results[0].IsNil() {
		// 没有数据可反序列化，直接返回nil数据
		return nil, err
	} else {
		var er error
		// 使用自定义序列化协议编码Data
		res, er = serializer.Encode(results[0].Interface())
		if er != nil {
			// 反序列化出错，返回nil数据
			return nil, er
		}
	}
	return res, err
}
