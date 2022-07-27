package main

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/guess"
	"go/token"
)

func main() {
	code := `package main

import (
	"github.com/labastack/echo/v4"
)

var a echo.Tracer
`
	packageNames := map[string]string{"github.com/labastack/echo/v4": "echo"}
	packageNameResolver := guess.WithMap(packageNames)
	resolver := goast.WithResolver(packageNameResolver)

	dec := decorator.NewDecoratorWithImports(token.NewFileSet(), "main", resolver)

	f, err := dec.Parse(code)
	if err != nil {
		panic(err)
	}

	res := decorator.NewRestorerWithImports("main", packageNameResolver)
	if err := res.Print(f); err != nil {
		panic(err)
	}

	code2 := `package rpc

import (
	"context"
	"testdata/proto"
)

type UserRpc struct {
}

// GetInfo is server rpc method as defined
func (this *UserRpc) GetInfo(ctx context.Context, req *proto.GetInfoReq, reply *proto.GetInfoReply) (err error) {
	// TODO: add some code
	*reply = proto.GetInfoReply{}
	return nil
}`

	df, err := decorator.Parse(code2)
	if err != nil {
		panic(err)
	}

	dst.Print(df)

	code3 := `package services

import (
	"github.com/bmbstack/ripple/sync"

	"github.com/bmbstack/ripple/fixtures/forum/proto"
)

var (
	playerClient     *proto.PlayerClient
	playerClientOnce sync.Once
)

func GetPlayerClient() *proto.PlayerClient {
	playerClientOnce.Do(func() {
		closePlayerClient()
		playerClient = proto.NewPlayerClient(func() {
			playerClientOnce.Reset()
		})
	})
	return playerClient
}

func closePlayerClient() {
	if playerClient != nil {
		playerClient.Discovery.Close()
		playerClient.XClient.Close()
	}
}

func CloseRpcClients() {
	closePlayerClient()
}`

	df2, err := decorator.Parse(code3)
	if err != nil {
		panic(err)
	}

	dst.Print(df2)
}
