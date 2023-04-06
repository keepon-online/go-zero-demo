package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/nacos/mall/order/api/internal/config"
	"go-zero-demo/nacos/mall/order/api/internal/types"
	"log"
	"net/http"
	"strings"
)

type CasbinMiddleware struct {
	Config config.Config
}

func NewCasbinMiddleware(c config.Config) *CasbinMiddleware {
	return &CasbinMiddleware{
		Config: c,
	}
}

func (m *CasbinMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	logx.Info("example middle")
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		err := check(req.Id, r.RequestURI, r.Method)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// Passthrough to next handler if need

		next(w, r)
	}
}

func check(sub, obj, act string) error {
	// 使用MySQL数据库初始化一个gorm适配器
	a, err := gormadapter.NewAdapter("mysql", "root:123456@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}

	m, err := model.NewModelFromString(`
										[request_definition]
											r = sub, obj, act
											
											[policy_definition]
											p = sub, obj, act
											
											[role_definition]
											g = _, _
											
											[policy_effect]
											e = some(where (p.eft == allow))
											
											[matchers]
											m = r.sub == p.sub && ParamsMatch(r.obj,p.obj) && r.act == p.act
								`)
	if err != nil {
		log.Fatalf("error: model: %s", err)
	}

	e, err := casbin.NewEnforcer(m, a)

	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
	}
	e.AddFunction("ParamsMatch", ParamsMatchFunc)
	//sub := "alice" // 想要访问资源的用户。
	//obj := "data1" // 将被访问的资源。
	//act := "read"  // 用户对资源执行的操作。
	ok, err := e.Enforce(sub, obj, act)
	if err != nil {
		// 处理err
	}
	if ok == true {
		// 允许alice读取data1
		fmt.Println("允许alice读取data1")
		return nil
	} else {
		// 拒绝请求，抛出异常
		fmt.Println("不允许alice读取data1")
		return errors.New("无权限")
	}

	// 您可以使用BatchEnforce()来批量执行一些请求
	// 这个方法返回布尔切片，此切片的索引对应于二维数组的行索引。
	// 例如results[0] 是{"alice", "data1", "read"}的结果

}
func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)
	return ParamsMatch(name1, name2), nil
}
func ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key2)
}
