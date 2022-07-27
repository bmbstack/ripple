package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/internal/controllers"
	. "github.com/bmbstack/ripple/fixtures/forum/internal/initial"
	_ "github.com/bmbstack/ripple/fixtures/forum/internal/models/one"
	_ "github.com/bmbstack/ripple/fixtures/forum/internal/models/two"
	"github.com/bmbstack/ripple/fixtures/forum/internal/rpc"
	"github.com/bmbstack/ripple/fixtures/forum/proto"
)

func RunServer() {
	Logger.Info("Run server ....")
	controllers.RouteAPI()
	ripple.Default().Run()
}

func RunRpc() {
	Logger.Info("Run rpc server ....")
	ripple.Default().RegisterRpc(proto.ServiceNameOfStudent, &rpc.StudentRpc{}, "")
	ripple.Default().RegisterRpc(proto.ServiceNameOfTeacher, &rpc.TeacherRpc{}, "")
	ripple.Default().RunRpc()
}
