package common

import (
	"encoding/binary"
	"net"
)

const NumOfLengthBytes = 8

func ReadMsg(conn net.Conn) ([]byte, error) {
	// ④读取接收数据的长度
	lenBs := make([]byte, NumOfLengthBytes)
	_, err := conn.Read(lenBs)
	if err != nil {
		return nil, err
	}
	// ⑤构造并读取接收数据
	length := binary.BigEndian.Uint64(lenBs)
	data := make([]byte, length)
	_, err = conn.Read(data)
	return data, err
}

func WriteMsg(conn net.Conn, data []byte) error {
	req := EncodeMsg(data)
	_, err := conn.Write(req)
	return err
}

// 编码消息
func EncodeMsg(data []byte) []byte {
	// 计算请求数据的长度
	reqLen := len(data)
	// 构建请求数据 data =reqLen的64位表示 + respData
	reqData := make([]byte, reqLen+NumOfLengthBytes)
	// ①把长度写进去前八个字节
	binary.BigEndian.PutUint64(reqData[:NumOfLengthBytes], uint64(reqLen))
	// ②写入数据
	copy(reqData[NumOfLengthBytes:], data)
	return reqData
}
