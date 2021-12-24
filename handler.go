package main

import (
	"context"
	"log"
	"user/idl/gen/user"
	"user/service"
)

type Handler struct {
	*user.UnimplementedUserServiceServer
}

func (*Handler) Login(ctx context.Context, req *user.LoginReq) (*user.LoginResp, error) {
	log.Printf("req:%+v", req)
	lCtx := service.NewLoginContext(ctx, req)
	service.NewLoginService().Login(lCtx)
	log.Printf("resp:%+v", lCtx.Resp)

	return lCtx.Resp, nil
}

func (*Handler) Register(ctx context.Context, req *user.RegisterReq) (*user.RegisterResp, error) {
	log.Printf("req:%+v", req)
	rCtx := service.NewRegisterContext(ctx, req)
	service.NewRegisterService().Register(rCtx)
	log.Printf("resp:%+v", rCtx.Resp)

	return rCtx.Resp, nil
}

func (*Handler) Query(ctx context.Context, req *user.QueryReq) (*user.QueryResp, error) {
	log.Printf("req:%+v", req)
	qCtx := service.NewQueryContext(ctx, req)
	service.NewQueryService().Query(qCtx)
	log.Printf("resp:%+v", qCtx.Resp)

	return qCtx.Resp, nil
}
