package rpc

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.learn.rpc/micro/rpc/message"
	"go.learn.rpc/micro/rpc/serialize/json"
	"testing"
)

func Test_setFuncField(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) Proxy
		service Service
		wantErr error
	}{
		{
			name: "nil",
			mock: func(ctrl *gomock.Controller) Proxy {
				return NewMockProxy(ctrl)
			},
			service: nil,
			wantErr: errors.New("rpc : 不支持service为 nil"),
		},
		{
			name: "no pointer",
			mock: func(ctrl *gomock.Controller) Proxy {
				return NewMockProxy(ctrl)
			},
			service: UserService{},
			wantErr: errors.New("rpc : 只支持指向结构体的一级指针"),
		},
		{
			name: "user service",
			mock: func(ctrl *gomock.Controller) Proxy {
				p := NewMockProxy(ctrl)
				p.EXPECT().Invoke(gomock.Any(), &message.Request{
					ServiceName: "user_service",
					MethodName:  "GetById",
					Data:        []byte(`{"Id" : "123"}"`),
				}).Return(&message.Response{}, nil)
				return p
			},
			service: &UserService{},
		},
	}
	s := json.Serializer{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			err := setFuncField(tc.service, tc.mock(ctrl), s)
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
