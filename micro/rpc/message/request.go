package message

import (
	"bytes"
	"encoding/binary"
)

type Request struct {
	// 头部消息长度
	HeadLength uint32
	// 消息体长度
	BodyLength uint32
	// 消息Id
	MessageId uint32
	// 版本
	Version uint8
	// 压缩算法
	Compresser uint8
	// 序列化协议
	Serializer uint8
	// 服务名
	ServiceName string
	// 方法名
	MethodName string
	// 扩展字段，用于填充元数据
	Meta map[string]string
	// 消息体
	Data []byte
}

func EnCodeReq(req *Request) []byte {
	// 开辟空间
	bs := make([]byte, req.HeadLength+req.BodyLength)
	// 1.写入HeadLength
	binary.BigEndian.PutUint32(bs[:4], req.HeadLength)
	// 2.写入BodyLength
	binary.BigEndian.PutUint32(bs[4:8], req.BodyLength)
	// 3.写入MessageId
	binary.BigEndian.PutUint32(bs[8:12], req.MessageId)
	// 4.写入Version, 因为version就一个字节，所以不需要编码了
	bs[12] = req.Version
	// 5.写入Compresser
	bs[13] = req.Compresser
	// 6.写入Serializer
	bs[14] = req.Serializer
	// 7.写入ServiceName
	cur := bs[15:]
	copy(cur, req.ServiceName)
	cur = cur[len(req.ServiceName):]
	cur[0] = '\n'
	cur = cur[1:]
	// 8.写入MethodName
	copy(cur, req.MethodName)
	cur = cur[len(req.MethodName):]
	cur[0] = '\n'
	cur = cur[1:]
	// 9.写入Meta扩展字段
	// meta扩展字段的key,val内部使用'\r'作为分隔符，每一对之间仍然使用'\n'作为分隔符
	for key, val := range req.Meta {
		copy(cur, key)
		cur = cur[len(key):]
		cur[0] = '\r'
		cur = cur[1:]
		copy(cur, val)
		cur = cur[len(val):]
		cur[0] = '\n'
		cur = cur[1:]
	}
	// 10.写入Data
	copy(cur, req.Data)
	return bs
}

func DeCodeReq(data []byte) *Request {
	req := &Request{}
	req.HeadLength = binary.BigEndian.Uint32(data[:4])
	req.BodyLength = binary.BigEndian.Uint32(data[4:8])
	req.MessageId = binary.BigEndian.Uint32(data[8:12])
	req.Version = data[12]
	req.Compresser = data[13]
	req.Serializer = data[14]
	header := data[15:req.HeadLength]
	// 近似于user_Service\nGetById
	index := bytes.IndexByte(header, '\n')
	// 引入分隔符，切分service name 和 method name
	req.ServiceName = string(header[:index])
	header = header[index+1:]

	// 切出来MethodName
	index = bytes.IndexByte(header, '\n')
	req.MethodName = string(header[:index])
	header = header[index+1:]
	// 切分出来Meta
	index = bytes.IndexByte(header, '\n')
	if index != -1 {
		meta := make(map[string]string, 16)
		for index != -1 {
			pair := header[:index]
			// '\r'的位置
			pairIndex := bytes.IndexByte(pair, '\r')
			key := string(pair[:pairIndex])
			val := string(pair[pairIndex+1:])
			meta[key] = val
			header = header[index+1:]
			index = bytes.IndexByte(header, '\n')
		}
		req.Meta = meta
	}

	// 读取Data
	if req.BodyLength != 0 {
		req.Data = data[req.HeadLength:]
	}
	return req
}

func (req *Request) CalculateHeaderLength() {
	headLength := 15 + len(req.ServiceName) + 1 + len(req.MethodName) + 1
	for key, val := range req.Meta {
		headLength += len(key)
		headLength++
		headLength += len(val)
		headLength++
	}
	req.HeadLength = uint32(headLength)
}
func (req *Request) CalculateBodyLength() {
	req.BodyLength = uint32(len(req.Data))
}
