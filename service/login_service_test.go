package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"user/constant"
	"user/idl/gen/user"
)

func TestLoginService_Login(t *testing.T) {
	ctx := NewLoginContext(context.Background(), &user.LoginReq{
		Username: "z00",
		Password: "z00",
	})
	NewLoginService().Login(ctx)
	assert.Equal(t, constant.RetSuccess.Code, ctx.Resp.BaseResp.ErrNo)
	assert.NotEmpty(t, ctx.Resp.AccessToken)
	assert.NotEmpty(t, ctx.Resp.RefreshToken)
}
