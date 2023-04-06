package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
	"go-zero-demo/nacos/mall/order/api/internal/config"
	"go-zero-demo/nacos/mall/order/api/internal/handler"
	"go-zero-demo/nacos/mall/order/api/internal/svc"
)

var configFile = flag.String("f", "nacos/mall/order/api/etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
