package scripts

import (
	"context"
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/controllers"
	_ "github.com/bmbstack/ripple/fixtures/forum/controllers"
	. "github.com/bmbstack/ripple/fixtures/forum/initial"
	_ "github.com/bmbstack/ripple/fixtures/forum/models/one"
	_ "github.com/bmbstack/ripple/fixtures/forum/models/two"
	"github.com/bmbstack/ripple/fixtures/forum/proto"
)

func RunServer() {
	Logger.Info("Run server ....")
	controllers.RouteAPI()
	ripple.Default().Run()
}

func RunRpc() {
	Logger.Info("Run rpc server ....")
	ripple.Default().RegisterRpc(proto.ServiceNameOfUser, &UserRpc{}, "")
	ripple.Default().RunRpc()
}

type UserRpc struct {
}

// GetInfo is server rpc method as defined
func (s *UserRpc) GetInfo(ctx context.Context, req *proto.GetInfoReq, reply *proto.GetInfoReply) (err error) {
	// TODO: add business logics
	*reply = proto.GetInfoReply{}
	reply.Name = "tomcat"
	return nil
}
