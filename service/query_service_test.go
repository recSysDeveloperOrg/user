package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"user/idl/gen/user"
)

func TestQueryService_Query(t *testing.T) {
	loginCtx := NewLoginContext(context.Background(), &user.LoginReq{
		Username: "z00",
		Password: "z00",
	})
	NewLoginService().Login(loginCtx)
	accessToken := loginCtx.AccessToken
	ctx := NewQueryContext(context.Background(), &user.QueryReq{
		AccessToken: accessToken,
	})
	NewQueryService().Query(ctx)
	assert.Equal(t, "100000000000000000000001", ctx.Resp.User.Id)

	ctx = NewQueryContext(context.Background(), &user.QueryReq{
		RefreshToken: loginCtx.RefreshToken,
	})
	NewQueryService().Query(ctx)
	assert.Equal(t, "100000000000000000000001", ctx.Resp.User.Id)
	assert.NotEmpty(t, ctx.Resp.AccessToken)
	t.Logf("got new accessToken:%s", ctx.Resp.AccessToken)
}
