package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	. "user/constant"
	"user/cred"
	"user/idl/gen/user"
	"user/model"
)

type LoginService struct {
}

type LoginContext struct {
	Ctx     context.Context
	Req     *user.LoginReq
	Resp    *user.LoginResp
	ErrCode *ErrorCode

	AccessToken  string
	RefreshToken string
}

var loginService *LoginService
var loginServiceOnce sync.Once

func NewLoginService() *LoginService {
	loginServiceOnce.Do(func() {
		loginService = &LoginService{}
	})

	return loginService
}

func NewLoginContext(ctx context.Context, req *user.LoginReq) *LoginContext {
	return &LoginContext{
		Ctx: ctx,
		Req: req,
		Resp: &user.LoginResp{
			BaseResp: &user.BaseResp{},
		},
	}
}

func (s *LoginService) Login(ctx *LoginContext) {
	defer s.buildResponse(ctx)
	if s.checkParams(ctx); ctx.ErrCode != nil {
		return
	}
	s.login(ctx)
}

func (s *LoginService) login(ctx *LoginContext) {
	loginUser, err := model.NewUserDao().QueryByUsernamePassword(ctx.Ctx, ctx.Req.Username, ctx.Req.Password)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetReadRepoErr)
		return
	}
	if loginUser == nil {
		ctx.ErrCode = BuildErrCode(ctx.Req.Username, RetNoSuchUserErr)
		return
	}

	accessToken, err := cred.IssueJWT(loginUser.ID, cred.ExpireDayAccessToken)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetSysErr)
		return
	}
	ctx.AccessToken = accessToken

	refreshToken, err := cred.IssueJWT(loginUser.ID, cred.ExpireDayRefreshToken)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetSysErr)
		return
	}
	ctx.RefreshToken = refreshToken
	affectRows, err := model.NewUserDao().UpdateRefreshToken(ctx.Ctx, loginUser.ID, refreshToken)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetWriteRepoErr)
		return
	}
	if affectRows == 0 {
		ctx.ErrCode = BuildErrCode(fmt.Sprintf("no user found:%s", loginUser.ID), RetNoSuchUserErr)
	}
}

func (*LoginService) checkParams(ctx *LoginContext) {
	if strings.TrimSpace(ctx.Req.Username) == "" {
		ctx.ErrCode = BuildErrCode("empty username", RetParamsErr)
		return
	}
	if strings.TrimSpace(ctx.Req.Password) == "" {
		ctx.ErrCode = BuildErrCode("empty password", RetParamsErr)
		return
	}
}

func (*LoginService) buildResponse(ctx *LoginContext) {
	errCode := RetSuccess
	if ctx.ErrCode != nil {
		errCode = ctx.ErrCode
	}

	ctx.Resp.BaseResp.ErrNo, ctx.Resp.BaseResp.ErrMsg = errCode.Code, errCode.Msg
	ctx.Resp.AccessToken, ctx.Resp.RefreshToken = ctx.AccessToken, ctx.RefreshToken
}
