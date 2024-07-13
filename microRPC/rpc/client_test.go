package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setFuncField(t *testing.T) {
	testCases := []struct {
		name    string
		service Service
		wantErr error
	}{
		{
			name:    "nil",
			service: nil,
			wantErr: errors.New("rpc : 不支持service为 nil"),
		},
		{
			name:    "no pointer",
			service: UserService{},
			wantErr: errors.New("rpc : 只支持指向结构体的一级指针"),
		},
		{
			name:    "user service",
			service: &UserService{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := setFuncField(tc.service)
			assert.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}
			resp, err := tc.service.(*UserService).GetById(context.Background(), &GetByUserIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, err)
			t.Log(resp)
		})
	}
}
