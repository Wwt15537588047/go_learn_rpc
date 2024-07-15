package serialize

type Serializer interface {
	Code() uint8
	Encode(val any) ([]byte, error)
	// 规定val必须是指向结构体的指针
	DeCode(data []byte, val any) error
}

// 上述序列化协议中能否将val设计为泛型，不能
// 客户端发送数据的时候通过反射的数据其类型为interface{}，并不能确定其具体类型，无法传入泛型具体实现
