package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"user/constant"
	"user/idl/gen/user"
)

func TestRegisterService_Register(t *testing.T) {
	name := testRandomString(26)
	ctx := NewRegisterContext(context.Background(), &user.RegisterReq{
		User: &user.User{
			Name:     name,
			Password: "lijunyu",
			Gender:   user.Gender_GENDER_MALE,
		},
	})
	NewRegisterService().Register(ctx)
	assert.Equal(t, constant.RetSuccess.Code, ctx.Resp.BaseResp.ErrNo)
	ctx = NewRegisterContext(context.Background(), &user.RegisterReq{
		User: &user.User{
			Name:     name,
			Password: "lijunyu",
			Gender:   user.Gender_GENDER_MALE,
		},
	})
	NewRegisterService().Register(ctx)
	assert.Equal(t, constant.RetDuplicateUsernameErr.Code, ctx.Resp.BaseResp.ErrNo)
}
