package rpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInitClientProxy(t *testing.T) {
	server := NewServer1()
	server.RegisterServer(&UserServiceServer{})
	go func() {
		err := server.Start("tcp", ":8081")
		t.Log(err)
	}()
	time.Sleep(3 * time.Second)
	usClient := &UserService{}
	err := InitClientProxy(":8081", usClient)
	require.NoError(t, err)
	resp, err := usClient.GetById(context.Background(), &GetByUserIdReq{Id: 123})
	require.NoError(t, err)
	assert.Equal(t, &GetByUserIdResp{
		Msg: "hello world",
	}, resp)
}
