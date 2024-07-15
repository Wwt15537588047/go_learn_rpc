package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseEncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		resp *Response
	}{
		{
			name: "normal",
			resp: &Response{
				MessageId: 12,
				//Version:    13,
				Compresser: 32,
				//Serializer: 7,
				Error: []byte("Error"),
				Data:  []byte("Hello world"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.resp.CalculateHeaderLength()
			tc.resp.CalculateBodyLength()
			data := EnCodeResp(tc.resp)
			resp := DeCodeResp(data)
			assert.Equal(t, tc.resp, resp)
		})
	}
}
