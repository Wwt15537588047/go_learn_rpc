package message

import (
	"encoding/binary"
)

type Response struct {
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
	// 错误
	Error []byte
	// 消息体
	Data []byte
}

func EnCodeResp(resp *Response) []byte {
	bs := make([]byte, resp.HeadLength+resp.BodyLength)
	// 1.HeadLength的处理
	binary.BigEndian.PutUint32(bs[:4], resp.HeadLength)
	// 2.写入BodyLength
	binary.BigEndian.PutUint32(bs[4:8], resp.BodyLength)
	// 3.写入MessageId
	binary.BigEndian.PutUint32(bs[8:12], resp.MessageId)
	// 4.写入Version
	bs[12] = resp.Version
	// 5.写入Compresser
	bs[13] = resp.Compresser
	// 6.写入Serializer
	bs[14] = resp.Serializer
	copy(bs[15:resp.HeadLength], resp.Error)
	copy(bs[resp.HeadLength:], resp.Data)
	return bs
}

func DeCodeResp(data []byte) *Response {
	resp := &Response{}
	resp.HeadLength = binary.BigEndian.Uint32(data[:4])
	resp.BodyLength = binary.BigEndian.Uint32(data[4:8])
	resp.MessageId = binary.BigEndian.Uint32(data[8:12])
	resp.Version = data[12]
	resp.Compresser = data[13]
	resp.Serializer = data[14]
	if resp.HeadLength != 15 {
		resp.Error = data[15:resp.HeadLength]
	}
	if resp.BodyLength != 0 {
		resp.Data = data[resp.HeadLength:]
	}
	return resp
}

func (resp *Response) CalculateHeaderLength() {
	headerLength := 15 + len(resp.Error)
	resp.HeadLength = uint32(headerLength)
}

func (resp *Response) CalculateBodyLength() {
	resp.BodyLength = uint32(len(resp.Data))
}
