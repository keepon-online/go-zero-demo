package main

import (
	"flag"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/zeromicro/zero-contrib/zrpc/registry/nacos"

	"go-zero-demo/nacos/mall/user/rpc/internal/config"
	"go-zero-demo/nacos/mall/user/rpc/internal/server"
	"go-zero-demo/nacos/mall/user/rpc/internal/svc"
	"go-zero-demo/nacos/mall/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewUserServer(ctx)

	serverRpc := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// 注册服务
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           50000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	opts := nacos.NewNacosConfig("user.rpc", c.ListenOn, sc, cc)
	_ = nacos.RegisterService(opts)
	serverRpc.Start()
}
