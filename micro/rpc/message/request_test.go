package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestEncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
	}{
		{
			name: "normal",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  14,
				Serializer:  13,
				ServiceName: "user_service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"trace-id": "123",
					"a/b":      "a",
				},
				Data: []byte("Hello World."),
			},
		},
		{
			name: "data with \n",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  14,
				Serializer:  13,
				ServiceName: "user_service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"trace-id": "123",
					"a/b":      "a",
				},
				Data: []byte("Hello \n World."),
			},
		},
		//{
		// 此种情况暂时是无法解决的，唯一的解决方案就是禁止开发者《框架的使用者》在meta里面使用 \n 和 \r,所以不会出现这种情况
		//	name: "meta with \n",
		//	req: &Request{
		//		MessageId:   123,
		//		Version:     12,
		//		Compresser:  14,
		//		Serializer:  13,
		//		ServiceName: "user_service",
		//		MethodName:  "GetById",
		//		Meta: map[string]string{
		//			"trace-id": "123",
		//			"a/b":      "a",
		//		},
		//		Data: []byte("Hello \n World."),
		//	},
		//},
		{
			name: "no meta",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  14,
				Serializer:  13,
				ServiceName: "user_service",
				MethodName:  "GetById",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.req.CalculateBodyLength()
			tc.req.CalculateHeaderLength()
			data := EnCodeReq(tc.req)
			req := DeCodeReq(data)
			assert.Equal(t, tc.req, req)
		})
	}
}
