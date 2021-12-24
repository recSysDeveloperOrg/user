package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	. "user/constant"
	"user/idl/gen/user"
	"user/model"
)

type RegisterService struct {
}

type RegisterContext struct {
	Ctx     context.Context
	Req     *user.RegisterReq
	Resp    *user.RegisterResp
	ErrCode *ErrorCode
}

var registerService *RegisterService
var registerServiceOnce sync.Once

func NewRegisterService() *RegisterService {
	registerServiceOnce.Do(func() {
		registerService = &RegisterService{}
	})

	return registerService
}

func NewRegisterContext(ctx context.Context, req *user.RegisterReq) *RegisterContext {
	return &RegisterContext{
		Ctx: ctx,
		Req: req,
		Resp: &user.RegisterResp{
			BaseResp: &user.BaseResp{},
		},
	}
}

func (s *RegisterService) Register(ctx *RegisterContext) {
	defer s.buildResponse(ctx)
	if s.checkParams(ctx); ctx.ErrCode != nil {
		return
	}
	s.register(ctx)
}

func (*RegisterService) checkParams(ctx *RegisterContext) {
	if strings.TrimSpace(ctx.Req.User.Name) == "" {
		ctx.ErrCode = BuildErrCode("empty username", RetParamsErr)
		return
	}
	if strings.TrimSpace(ctx.Req.User.Password) == "" {
		ctx.ErrCode = BuildErrCode("empty password", RetParamsErr)
	}
}

func (*RegisterService) register(ctx *RegisterContext) {
	if err := model.NewUserDao().InsertNewUser(ctx.Ctx, &model.User{
		Username: ctx.Req.User.Name,
		Password: ctx.Req.User.Password,
	}); err != nil {
		if errors.Is(err, model.ErrDuplicateName) {
			ctx.ErrCode = BuildErrCode("duplicate user name", RetDuplicateUsernameErr)
			return
		}
		ctx.ErrCode = BuildErrCode(err, RetWriteRepoErr)
	}
}

func (*RegisterService) buildResponse(ctx *RegisterContext) {
	errCode := RetSuccess
	if ctx.ErrCode != nil {
		errCode = ctx.ErrCode
	}

	ctx.Resp.BaseResp.ErrNo, ctx.Resp.BaseResp.ErrMsg = errCode.Code, errCode.Msg
}
