[![GoDoc](https://pkg.go.dev/badge/github.com/bmbstack/ripple?status.svg)](https://pkg.go.dev/github.com/bmbstack/ripple?tab=doc)
[![Release](https://img.shields.io/github/release/bmbstack/ripple.svg?style=flat-square)](https://github.com/bmbstack/ripple/releases)

# Ripple

a lightweight web framework for Go(base on [Echo](https://github.com/labstack/echo))

## Installation

```shell
go get github.com/bmbstack/ripple@latest
go install github.com/bmbstack/ripple/cmd/ripple@latest
go install github.com/bmbstack/ripple/protoc/protoc-gen-gofast@latest
ripple new rippleApp
cd $GOPATH/src/rippleApp
go mod init
go mod tidy
go mod vendor
go run cmd/main.go --env dev s
```

## Upgrade

```shell
go get github.com/bmbstack/ripple@latest
go install github.com/bmbstack/ripple/cmd/ripple@latest
go install github.com/bmbstack/ripple/protoc/protoc-gen-gofast@latest
go mod tidy
go mod vendor
go run cmd/main.go --env dev s
```

If you use nacos, we recommend you:

```shell
go get github.com/smallnest/rpcx@v1.7.3
```

Then, Open the url: [http://127.0.0.1:8090](http://127.0.0.1:8090)

## Command

```
NAME:
   ripple - Command line tool to managing your Ripple application

USAGE:
   ripple [global options] command [command options] [arguments...]

VERSION:
   1.2.2

AUTHOR:
   wangmingjob <wangmingjob@icloud.com>

COMMANDS:
   new      Create a Ripple application
            desc: ripple new appName, however this appName can be empty, will be generated in the current directory
            ripple new
            ripple new app
   run, r   Run the Ripple application
   gen, g   Auto generate code, *.proto => *.pb.go *.rpc.go rpc.client.go; *.dto.go => *.controller.go && *.service.go
            desc: ripple g path component name/pbPath (path: dir/file; component: ''/proto/controller/service, name: component name, pbPath: *.pb.go path)
            ripple g
            ripple g packages/app
            ripple g packages/app proto
            ripple g packages/app controller
            ripple g packages/app service
            ripple g packages/app service product
            ripple g packages/app ecode
            ripple g packages/app/proto/user.proto
            ripple g packages/app/internal/dto/user.dto.go
            ripple g packages/app2 rpc.client packages/app1/proto/user.pb.go
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --env value    执行环境 (开发环境dev(只有ripple作者会使用)、线上环境prod) (default: prod)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

Note: `ripple g` or `ripple g path` or `ripple g path component` will generate some code, and the path must be the parent directory of go.mod, that is, the main directory of the project. For example, the demo project `fixtures/form`, and you will execute command:

```shell
ripple g fixture/form
```

or

```shell
cd fixture/form
ripple g
```

* *.proto => \*.pb.go, *.rpc.go(internal/rpc)
* *.dto.go => *.controller.go(internal/controllers/v1), *.service.go(internal/service)

*.proto example

```protobuf
syntax = "proto3";

// @RippleRpc
// @NacosGroup DEFAULT_GROUP
// @NacosCluster ripple
package proto;

// The Student service definition.
service Student {
  rpc Learn (LearnReq) returns (LearnReply) {}
}

message LearnReq {
  uint64 id = 1;
}

message LearnReply {
  string name = 1;
}

```

`Note`: @RippleRpc, @NacosGroup, @NacosCluster

`*.dto.go` example:

```go
// ReqStudentLearn
// @RippleApi
// @Uri /student/learn
// @Method POST
type ReqStudentLearn struct {
    ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespStudentLearn struct {
    Name string `json:"name"`
}

```

`Note`: @RippleApi, @Uri, @Method

## Command: ripple g

```shell
# generate all files, contains *.pb.go, *.controller.go, *.service.go, *.rpc.go(source file: *.dto.go, *.proto)
ripple g
# generate all files in this directory, contains *.pb.go, *.controller.go, *.service.go, *.rpc.go(source file: *.dto.go, *.proto)
ripple g packages/app
# just generate *.pb.go(source file: *.proto)
ripple g packages/app proto
# just generate *.controller.go(source file: *.dto.go)
ripple g packages/app controller
# just generate *.service.go(source file: *.dto.go)
ripple g packages/app service
# just generate *.service.go(source: serviceName, this product is a serviceName)
ripple g packages/app service product
# just generate *.ecode.go
ripple g packages/app ecode
# just generate *.pb.go(souce file: *.proto)
ripple g packages/app/proto/user.proto
# just generate *.controller.go, *.service.go(souce file: *.dto.go)
ripple g packages/app/internal/dto/user.dto.go
# just generate rpcclient/rpc.client.go in app2(souce file: app1's *.pb.go)
ripple g packages/app2 rpc.client packages/app1/proto/user.pb.go
```

This is the structure of the `rippleApp` list application that will showcase how you can build web apps with `ripple`:

```shell
.
├── Makefile
├── cmd                               // 程序入口
│   └── main.go
├── config                            // 配置文件
│   ├── config.dev.yaml
│   ├── config.prod.yaml
│   └── config.test.yaml
├── frontend                          // 前端页面
│   ├── static
│   └── templates
├── go.mod
├── go.sum
├── internal
│   ├── controllers                  // `ripple gen`固定输出目录 controllers/v1
│   ├── dto                          // `ripple gen`的输入源 *.dto.go
│   ├── helper
│   ├── initial
│   ├── rpc                          // `ripple gen`固定输出目录
│   ├── rpcclient                    // `ripple gen`固定输出目录
│   ├── scripts
│   └── services                     // `ripple gen`固定输出目录
└── proto                            // `ripple gen`固定输入源 *.proto
    ├── user.pb.go
    └── user.proto

```

## Rpc client call, eg (fixture/form):

```go
sc := services.GetStudentClient()
reply, _ := sc.Learn(context.Background(), &proto.LearnReq{Id: 1})
```

## Rpc server register, eg (fixture/form):

```go
ripple.Default().RegisterRpc(proto.ServiceNameOfStudent, &rpc.StudentRpc{}, "")
ripple.Default().RunRpc()

// service impl
type StudentRpc struct {
}

// Learn is server rpc method as defined
func (this *StudentRpc) Learn(ctx context.Context, req *proto.LearnReq, reply *proto.LearnReply) (err error) {
	// TODO: add some code
	*reply = proto.LearnReply{}
	reply.Name = "student learn function"
	return nil
}
```

## Cache (Redis)

```go
cache := ripple.Default().GetCache("cacheAlias")
result, err := cache.HGetAll(context.Background(), "key").Result()
```

or

```go
cache := ripple.Default().GetCache("cacheAlias")
client := cache.Client().(*redis.Client)
result, err := client.HGetAll(context.Background(), "key").Result()
```

## Logger

### SLS(aliyun log)
```go
ripple.Default().AddLogType(ripple.LogTypeSLS)
```
or
### CLS(tecent cloud log)
```go
ripple.Default().AddLogType(ripple.LogTypeCLS)
```

Use logger like this
```go
import "github.com/bmbstack/ripple/logger"

logger.With(map[string]interface{}{
    "userId": 101,
    "traceId": "lskajdfouiaadgvv",
}).Info("hello, tom")

logger.Info("hello, jack")
```
## Please close Resource, when the program exits

```go
rpcclient.CloseAll()
ripple.Default().StopRpc()
ripple.Default().CloseOrm()
ripple.Default().CloseCache()
```

## Which features to include in a framework

A framework has 3 parts. A router receiving a request and directing it to a handler, a middleware system to add reusable pieces of software before and after the handler, and the handler processing the request and writing the response.

- Router
- Middlewares
- Handler processing

<img src="https://raw.githubusercontent.com/bmbstack/ripple/master/screenshots/framework.png" width="720" height="576" />

Middlewares handle:

- error/panic
- logging
- security
- sessions
- cookies
- body parsing

## Features

* [x] MySQL and Foundation database support
* [x] Modular (you can choose which components to use)
* [x] Middleware support, compatible Middleware works out of the box
* [x] Lightweight
* [x] Multiple configuration files support (currently ```.env```)

## Overview

`ripple` is a lightweight framework. It is based on the principles of simplicity, relevance and elegance.

- Simplicity. The design is simple, easy to understand and doesn't introduce many layers between you and the standard library. It is a goal of the project that users should be able to understand the whole framework in a single day.
- Relevance. `ripple` doesn't assume anything. We focus on things that matter, this way we are able to ensure easy maintenance and keep the system well-organized, well-planned and sweet.
- Elegance. `ripple` uses golang best practises. We are not afraid of heights, it's just that we need a parachute in our backpack. The source code is heavily documented, any functionality should be well explained and well tested.

## Configurations

ripple support .env configurations files. In our rippleApp app, we put the configuration files in the current project directory. I have included all three formats for clarity, you can be just fine with either one of them.

ripple searches for a file named .env in the config directory. The first to be found is the one to be used.

This is the content of ```config.dev.yaml/config.test.yaml/config.prod.yaml``` file:

```shell
domain: "127.0.0.1:8090"
static: "frontend/static"
templates: "frontend/templates"
autoMigrate: false
bindAllTag: false
databases: [
  {
    "alias": "one",
    "dialect": "mysql",
    "host": "127.0.0.1",
    "port": 3306,
    "name": "one",
    "username": "root",
    "password": "123456"
    "maxIdleConns": 200,
    "maxOpenConns": 1000
  },
  {
    "alias": "two",
    "dialect": "mysql",
    "host": "127.0.0.1",
    "port": 3306,
    "name": "two",
    "username": "root",
    "password": "123456"
    "maxIdleConns": 200,
    "maxOpenConns": 1000
  }
]
caches: [
  {
    "alias": "one",
    "section": "one",
    "adapter": "redis",
    "host": "127.0.0.1",
    "port": 6379,
    "password": "123456"
  },
  {
    "alias": "two",
    "section": "two",
    "adapter": "redis",
    "host": "127.0.0.1",
    "port": 6379,
    "password": "123456"
  }
]
nacos:
  host: "my.nacos.com"
  port: 8848
  namespaceId: "public"
  cluster: "ripple_user"
  group: "DEFAULT_GROUP"
  failMode: "failover"
  selectMode: "roundRobin"
  clientPoolSize: 10
  cacheDir: "./cache"
  logDir: "./log"
  server: "127.0.0.1:18090"
sls:
  accessKeyId: "xxxxxxx"
  accessKeySecret: "xxxxxxx"
  endpoint: "cn-beijing-intranet.log.aliyuncs.com"
  allowLogLevel: "info"  # debug,info,warn,error
  closeStdout: false
  project: "xxxxxxx"
  logstore: "xxxxxxx"
  topic: "topic"
  source: "source"
cls:
  accessKeyId: "xxxxxxx"
  accessKeySecret: "xxxxxxx"
  endpoint: "ap-beijing.cls.tencentcs.com"
  allowLogLevel: "info"  # debug,info,warn,error
  closeStdout: false
  topic: "b246af67-dab3-408b-b802-ba150488ffbf"
```

## Models

ripple uses the [gorm](https://gorm.io/gorm) library as its Object Relational Mapper, so you won't need to learn anything fancy. In our rippleApp app, we need to define a User model that will be used to store our todo details.

In the file models/user.go we define our rippleApp model like this

```go
package models

import (
	"github.com/bmbstack/ripple"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login       string `sql:"size:255;not null"`
	Password    string `sql:"size:255;not null"`
	Email       string `sql:"size:255"`
	Avatar      string `sql:"size:255"`
	Github      string `gorm:"column:github"`
	QQ          string `gorm:"column:qq"`
	Weibo       string `gorm:"column:weibo"`
	Weixin      string `gorm:"column:weixin"`
	HomePage    string `gorm:"column:home_page"`
	Tagline     string
	Description string
	Location    string
}

func init() {
	ripple.Default().RegisterModels(&User{})
}

```

Notice that we need to register our model by calling ripple.RegisterModels(&User{}) in the init function otherwise ripple won't be aware of the model.

ripple will automatically create the table users if it doesn't exist.

Don't be confused by the schema tag, I just added them since we will use the schema package to decode form values(this has nothing to do with ripple, you can use whatever form library you fancy.)

## Controllers

ripple controllers are structs that implement the Controller interface. To help make ripple usable, Structs must implement the Controller interface.

```go
type Controller interface {
Path() string
}
```

This creates a new Echo group at the Controller#Path, in our example /posts, with all the defined actions.

```shell
 GET /posts     => #ActionIndex
 GET /posts/:id => #ActionShow
```

Our rippleApp Controller is in the controllers/home.go

```go
package controllers

import (
	"net/http"
	"github.com/bmbstack/ripple"
	"github.com/labstack/echo/v4"
)

type HomeController struct {
	Index  echo.HandlerFunc `controller:"GET /"`
	Html   echo.HandlerFunc `controller:"GET html"`
	String echo.HandlerFunc `controller:"GET string"`
}

func init() {
	ripple.Default().RegisterController(&HomeController{})
}

func (this HomeController) Path() string {
	return "/"
}

func (this HomeController) ActionIndex(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "home/index.html", map[string]interface{}{
		"title": "Hello, forum is a Ripple application ",
	})
}

func (this HomeController) ActionHtml(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "home/html.html", map[string]interface{}{
		"title": "Hello, this is a html template",
	})
}

func (this HomeController) ActionString(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, this is a string")
}

```

Note that we registered our controller by calling ripple.RegisterController($HomeController{}) in the init function so as to make ripple aware of our controller. See Routing section below for more explanation of what the controller is doing.

## Templates

ripple templates are golang templates(use [pongo2](https://github.com/flosch/pongo2)). This is the content of frontend/templates/home/index.html:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title></title>
</head>
<body>
<h3>{{title}}</h3>

<h3>Routes</h3>
<ul>
    <li>
        <a href="http://localhost:8090/html" target="_blank">Html template Page</a>
    </li>
    <li>
        <a href="http://localhost:8090/string" target="_blank">String Page</a>
    </li>
</ul>
</body>
</html>
```