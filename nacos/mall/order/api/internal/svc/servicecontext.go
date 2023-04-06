package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-demo/nacos/mall/order/api/internal/config"
	"go-zero-demo/nacos/mall/order/api/internal/middleware"
	"go-zero-demo/nacos/mall/user/rpc/userclient"
)

type ServiceContext struct {
	Config           config.Config
	UserRpc          userclient.User
	CasbinMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		UserRpc:          userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		CasbinMiddleware: middleware.NewCasbinMiddleware(c).Handle,
	}
}
