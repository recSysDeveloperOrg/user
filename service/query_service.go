package service

import (
	"context"
	"strings"
	"sync"
	. "user/constant"
	"user/cred"
	"user/idl/gen/user"
	"user/model"
)

type QueryService struct {
}

type QueryContext struct {
	Ctx     context.Context
	Req     *user.QueryReq
	Resp    *user.QueryResp
	ErrCode *ErrorCode

	AccessToken string
	User        *model.UserWithID
}

var queryService *QueryService
var queryServiceOnce sync.Once

var genderString2Gender = map[string]user.Gender{
	"MALE":   user.Gender_GENDER_MALE,
	"FEMALE": user.Gender_GENDER_FEMALE,
}

func NewQueryService() *QueryService {
	queryServiceOnce.Do(func() {
		queryService = &QueryService{}
	})

	return queryService
}

func NewQueryContext(ctx context.Context, req *user.QueryReq) *QueryContext {
	return &QueryContext{
		Ctx: ctx,
		Req: req,
		Resp: &user.QueryResp{
			BaseResp: &user.BaseResp{},
		},
	}
}

func (s *QueryService) Query(ctx *QueryContext) {
	defer s.buildResponse(ctx)
	if s.checkParams(ctx); ctx.ErrCode != nil {
		return
	}
	s.query(ctx)
}

func (*QueryService) checkParams(ctx *QueryContext) {
	if strings.TrimSpace(ctx.Req.AccessToken) == "" {
		ctx.ErrCode = BuildErrCode("empty access token", RetParamsErr)
		return
	}
	if strings.TrimSpace(ctx.Req.RefreshToken) == "" {
		ctx.ErrCode = BuildErrCode("empty refresh token", RetParamsErr)
	}
}

func (s *QueryService) query(ctx *QueryContext) {
	if len(ctx.Req.AccessToken) > 0 {
		s.queryByAccessToken(ctx)
		return
	}
	s.queryByRefreshToken(ctx)
}

func (s *QueryService) queryByAccessToken(ctx *QueryContext) {
	var userWithID *model.UserWithID
	if userWithID = s.tokenString2UserWithID(ctx, ctx.Req.RefreshToken); ctx.ErrCode != nil {
		return
	}
	ctx.User = userWithID
}

func (s *QueryService) queryByRefreshToken(ctx *QueryContext) {
	var userWithID *model.UserWithID
	if userWithID = s.tokenString2UserWithID(ctx, ctx.Req.RefreshToken); ctx.ErrCode != nil {
		return
	}

	if userWithID.LastRefreshToken != ctx.Req.RefreshToken {
		ctx.ErrCode = BuildErrCode("invalid token", RetInvalidTokenErr)
		return
	}
	accessToken, err := cred.IssueJWT(userWithID.ID, cred.ExpireDayAccessToken)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetSysErr)
		return
	}
	ctx.AccessToken = accessToken
	ctx.User = userWithID
}

func (*QueryService) tokenString2UserWithID(ctx *QueryContext, tokenString string) *model.UserWithID {
	userID, err := cred.ParseJWT(tokenString)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetInvalidTokenErr)
		return nil
	}
	userWithID, err := model.NewUserDao().QueryByID(ctx.Ctx, userID)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetReadRepoErr)
		return nil
	}
	if userWithID == nil {
		ctx.ErrCode = BuildErrCode("no such user", RetNoSuchUserErr)
		return nil
	}

	return userWithID
}

func (*QueryService) buildResponse(ctx *QueryContext) {
	errCode := RetSuccess
	if ctx.ErrCode != nil {
		errCode = ctx.ErrCode
	}

	ctx.Resp.BaseResp.ErrNo, ctx.Resp.BaseResp.ErrMsg = errCode.Code, errCode.Msg
	ctx.Resp.AccessToken = ctx.AccessToken
	gender := user.Gender_GENDER_UNDEFINED
	if g, ok := genderString2Gender[ctx.User.Gender]; ok {
		gender = g
	}
	ctx.Resp.User = &user.User{
		Id:     ctx.User.ID,
		Name:   ctx.User.Username,
		Gender: gender,
	}
}
