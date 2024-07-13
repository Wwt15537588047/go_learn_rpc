package test

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestGetUser(t *testing.T) {
	// 创建gomock控制器，用来记录后续操作信息
	mockCtl := gomock.NewController(t)
	// 调用mock文件中的NewMockMyInter方法，创建一个MyInter接口的mock示例
	mockMyInter := NewMockMyInter(mockCtl)

	// EXPECT()接口设置预期返回值
	mockMyInter.EXPECT().GetName(1).Return("SUCCESS")
	// 将mock的MyInteger传入GetUser函数
	resp := GetUser(mockMyInter, 1)
	if resp == "SUCCESS" {
		t.Log("right")
	} else {
		t.Error("error")
	}
}
