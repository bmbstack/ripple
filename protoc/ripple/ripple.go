package ripple

import (
	"bufio"
	"fmt"
	"github.com/bmbstack/ripple/cmd/ripple/util"
	"github.com/bmbstack/ripple/rst"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver/guess"
	"go/token"
	"golang.org/x/mod/modfile"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	pb "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const (
	ripplePkgPath           = "github.com/bmbstack/ripple"
	rippleHelperPkgPath     = "github.com/bmbstack/ripple/helper"
	rpcxClientPkgPath       = "github.com/smallnest/rpcx/client"
	rpcxNacosClientPkgPath  = "github.com/rpcxio/rpcx-nacos/client"
	rpcxNacosSdkConsPkgPath = "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	rpcxProtocolPkgPath     = "github.com/smallnest/rpcx/protocol"

	Permissions = 0755
)

func init() {
	generator.RegisterPlugin(new(ripple))
}

type ripple struct {
	gen *generator.Generator
}

// Name returns the name of this plugin
func (r *ripple) Name() string {
	return "ripple"
}

// Init initializes the plugin.
func (r *ripple) Init(gen *generator.Generator) {
	r.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (r *ripple) objectNamed(name string) generator.Object {
	r.gen.RecordTypeUse(name)
	return r.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (r *ripple) typeName(str string) string {
	return r.gen.TypeName(r.objectNamed(str))
}

// GenerateImports generates the import declaration for this file.
func (r *ripple) GenerateImports(file *generator.FileDescriptor) {
}

// P forwards to g.gen.r.
func (r *ripple) P(args ...interface{}) { r.gen.P(args...) }

// Generate generates code for the services in the given file.
func (r *ripple) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	_ = r.gen.AddImport(ripplePkgPath)
	_ = r.gen.AddImport(rippleHelperPkgPath)
	_ = r.gen.AddImport(rpcxClientPkgPath)
	_ = r.gen.AddImport(rpcxNacosClientPkgPath)
	_ = r.gen.AddImport(rpcxNacosSdkConsPkgPath)
	_ = r.gen.AddImport(rpcxProtocolPkgPath)
	_ = r.gen.AddImport("context")
	_ = r.gen.AddImport("errors")

	var group string
	var cluster string
	comments := r.gen.Comments(strconv.Itoa(2))
	if hasTag(comments, "@RippleRpc") {
		reader := strings.NewReader(comments)
		br := bufio.NewReader(reader)
		for {
			lineBytes, _, err := br.ReadLine()
			if err != nil || err == io.EOF {
				break
			}
			line := string(lineBytes)
			line = strings.Trim(trimAnnot(line), " ")

			if strings.HasPrefix(line, "@NacosGroup") {
				group = parseKey("@NacosGroup", line)
			} else if strings.HasPrefix(line, "@NacosCluster") {
				cluster = parseKey("@NacosCluster", line)
			}
		}
	}

	// generate all services
	for _, service := range file.FileDescriptorProto.Service {
		r.generateService(file, service, group, cluster)
	}
}

// generateService generates all the code for the named service
func (r *ripple) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, group, cluster string) {
	originServiceName := service.GetName()
	serviceName := upperFirstLatter(originServiceName)
	r.P("// This following code was generated by ripple")
	r.P(fmt.Sprintf("// Gernerated from %s", file.GetName()))

	// ServiceName
	r.P()
	r.P(fmt.Sprintf(`const ServiceNameOf%[1]s = "%[1]sRpc"`, serviceName))

	// interface
	r.P()
	r.P("//================== interface ===================")
	r.P(fmt.Sprintf(`type %sInterface interface {`, serviceName))
	r.P()
	r.P(fmt.Sprintf(`// %sInterface can be used for interface verification.`, serviceName))
	for _, method := range service.Method {
		r.generateInterfaceCode(method)
	}
	r.P(fmt.Sprintf(`}`))

	// server
	r.P()

	r.P("//================== server implement demo ===================")
	r.P("//ripple.Default().RegisterRpc(\"User\", &UserRpcDemo{}, \"\")")
	r.P(fmt.Sprintf(`type %[1]sRpcDemo struct {}`, serviceName))
	r.P()
	for _, method := range service.Method {
		r.generateServerDemoCode(service, method)
	}

	// use ast write *.rpc.go
	r.generateRpcSource(file, service)

	r.P()
	r.P("//================== client stub ===================")
	r.P(fmt.Sprintf(`// newXClientFor%[1]s creates a XClient.
		// You can configure this client with more options such as etcd registry, serialize type, select algorithm and fail mode.
		func newXClientFor%[1]s() (client.XClient, error) {
			config := ripple.GetBaseConfig()
			if helper.IsEmpty(config.Nacos) {
				return nil, errors.New("yaml nacos config is null")
			}
			clientConfig := constant.ClientConfig{
				TimeoutMs:            10 * 1000,
				ListenInterval:       30 * 1000,
				BeatInterval:         5 * 1000,
				NamespaceId:          config.Nacos.NamespaceId,
				CacheDir:             config.Nacos.CacheDir,
				LogDir:               config.Nacos.LogDir,
				UpdateThreadNum:      20,
				NotLoadCacheAtStart:  true,
				UpdateCacheWhenEmpty: true,
			}
		
			serverConfig := []constant.ServerConfig{{
				IpAddr: config.Nacos.Host,
				Port:   config.Nacos.Port,
			}}
		
			d, err := client1.NewNacosDiscovery(ServiceNameOf%[1]s, "%[2]s", "%[3]s", clientConfig, serverConfig)
			if err != nil {
				return nil, err
			}
			
			opt := client.DefaultOption
			opt.SerializeType = protocol.ProtoBuffer

			xclient := client.NewXClient(ServiceNameOf%[1]s, client.Failtry, client.RoundRobin, d, opt)

			return xclient,nil
		}

		// %[1]s is a client wrapped XClient.
		type %[1]sClient struct{
			xclient client.XClient
		}

		// New%[1]sClient wraps a XClient as %[1]sClient.
		// You can pass a shared XClient object created by NewXClientFor%[1]s.
		func New%[1]sClient() *%[1]sClient {
			xc, err := newXClientFor%[1]s()
			if err != nil {
				fmt.Println(fmt.Sprintf("Create rpcx client err: %s", err.Error()))
				return &%[1]sClient{}
			}
			return &%[1]sClient{xclient: xc}
		}
	`, serviceName, group, cluster))
	for _, method := range service.Method {
		r.generateClientCode(service, method)
	}
	r.P("// ======================================================")
}

func (r *ripple) generateInterfaceCode(method *pb.MethodDescriptorProto) {
	methodName := upperFirstLatter(method.GetName())
	inType := r.typeName(method.GetInputType())
	outType := r.typeName(method.GetOutputType())
	r.P(fmt.Sprintf(`// %[1]s is server rpc method as defined
		%[1]s(ctx context.Context, req *%[2]s, reply *%[3]s) (err error)
	`, methodName, inType, outType))
}

func (r *ripple) generateServerDemoCode(service *pb.ServiceDescriptorProto, method *pb.MethodDescriptorProto) {
	methodName := upperFirstLatter(method.GetName())
	serviceName := upperFirstLatter(service.GetName())
	inType := r.typeName(method.GetInputType())
	outType := r.typeName(method.GetOutputType())
	r.P(fmt.Sprintf(`
		func (this *%sRpcDemo) %s(ctx context.Context, req *%s, reply *%s) (err error){
			// TODO: add business logics
			*reply = %s{}
			return nil
		}
	`, serviceName, methodName, inType, outType, outType))
}

func (r *ripple) generateClientCode(service *pb.ServiceDescriptorProto, method *pb.MethodDescriptorProto) {
	methodName := upperFirstLatter(method.GetName())
	serviceName := upperFirstLatter(service.GetName())
	inType := r.typeName(method.GetInputType())
	outType := r.typeName(method.GetOutputType())
	r.P(fmt.Sprintf(`// %s is client rpc method as defined
		func (c *%sClient) %s(ctx context.Context, req *%s)(reply *%s, err error){
			reply = &%s{}
			err = c.xclient.Call(ctx,"%s",req, reply)
			return reply, err
		}
	`, methodName, serviceName, methodName, inType, outType, outType, method.GetName()))
}

// upperFirstLatter make the fisrt charater of given string  upper class
func upperFirstLatter(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return strings.ToUpper(string(s[0]))
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

//=======================================================
//             parse @NacosGroup @NacosCluster
//=======================================================
func trimAnnot(s string) string {
	s = strings.ReplaceAll(s, "//", "")
	s = strings.ReplaceAll(s, "/*", "")
	s = strings.ReplaceAll(s, "*/", "")
	return s
}

func hasTag(str, tag string) (ret bool) {
	str = strings.ReplaceAll(str, "\r", " ")
	str = strings.ReplaceAll(str, "\n", " ")
	arr := strings.Split(str, " ")
	for _, value := range arr {
		if value == tag {
			ret = true
			break
		}
	}
	return ret
}

func parseKey(key, line string) (value string) {
	arr := strings.Split(line, " ")
	if len(arr) >= 2 {
		if strings.EqualFold(key, arr[0]) {
			value = arr[1]
		}
	}
	return value
}

//=======================================================
//                  use AST to write *.rpc.go
//=======================================================
func (r *ripple) generateRpcSource(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto) {
	appPkg, err := getAppPkg()
	if err != nil {
		return
	}

	source := file.GetName()
	currentPath := source[0:strings.Index(source, "proto")]
	pkg := "rpc"
	dir := path.Join(path.Dir(currentPath), "internal", pkg)
	err = os.MkdirAll(dir, Permissions)
	if err != nil {
		return
	}
	target := filepath.Join(dir, fmt.Sprintf("%s.rpc.go", util.StartToLower(service.GetName())))
	upper := util.StartToUpper(service.GetName())
	protoPkg := filepath.Join(appPkg, currentPath, "proto")
	if !util.Exist(target) {
		// create file, write code to file
		_, err := os.Create(target)
		if err != nil {
			return
		}

		code := fmt.Sprintf(`// Code generated by ripple g, You can edit it again.
// source: %s

package %s

import (
	"context"
	"%s"
)

type %sRpc struct {
}
`, source, pkg, protoPkg, upper)

		for _, method := range service.Method {
			methodName := upperFirstLatter(method.GetName())
			serviceName := upperFirstLatter(service.GetName())
			inType := r.typeName(method.GetInputType())
			outType := r.typeName(method.GetOutputType())
			code = code + fmt.Sprintf(`
// %s is server rpc method as defined
func (this *%sRpc) %s(ctx context.Context, req *proto.%s, reply *proto.%s) (err error){
	// TODO: add some code
	*reply = proto.%s{}
	return nil
}
`, methodName, serviceName, methodName, inType, outType, outType)
		}

		_ = ioutil.WriteFile(target, []byte(code), Permissions)
	}

	// convert to ast, append action && setup uri
	df, err := rst.ParseSrcFile(target, guess.New())
	if err != nil {
		return
	}

	for _, method := range service.Method {
		methodName := upperFirstLatter(method.GetName())
		serviceName := upperFirstLatter(service.GetName())
		inType := r.typeName(method.GetInputType())
		outType := r.typeName(method.GetOutputType())
		rpcName := fmt.Sprintf("%sRpc", serviceName)

		// method
		hasMethod := rst.HasFuncDeclWithRecvInFile(df, dst.FuncDecl{
			Name: &dst.Ident{Name: methodName},
		}, rpcName)
		if !hasMethod {
			fd := createRpcMethod(appPkg, currentPath, rpcName, methodName, inType, outType)
			df.Decls = append(df.Decls, fd)
		}
	}

	buf := rst.PrintToBuf(df, guess.New(), []rst.PkgAlias{{Pkg: "context", Alias: ""}})
	_ = ioutil.WriteFile(target, buf.Bytes(), Permissions)
}

func getAppPkg() (string, error) {
	bytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return "", err
	}
	modFile, err := modfile.Parse("go.mod", bytes, nil)
	if err != nil {
		return "", err
	}
	return modFile.Module.Mod.Path, nil
}

func createRpcMethod(appPkg, currentPath, rpcName, methodName, intType, outType string) *dst.FuncDecl {
	return &dst.FuncDecl{
		Name: dst.NewIdent(methodName),
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{Names: []*dst.Ident{dst.NewIdent("this")}, Type: &dst.StarExpr{X: dst.NewIdent(rpcName)}},
			},
		},
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("ctx")},
						Type: &dst.Ident{
							Name: "Context",
							Path: "context",
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("req")},
						Type: &dst.StarExpr{
							X: &dst.Ident{
								Name: intType,
								Path: filepath.Join(appPkg, currentPath, "proto"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("reply")},
						Type: &dst.StarExpr{
							X: &dst.Ident{
								Name: outType,
								Path: filepath.Join(appPkg, currentPath, "proto"),
							},
						},
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("err")},
						Type:  dst.NewIdent("error"),
					},
				},
			},
		},
		Body: &dst.BlockStmt{
			Decs: dst.BlockStmtDecorations{Lbrace: dst.Decorations{"\n"}},
			List: []dst.Stmt{
				&dst.AssignStmt{
					Decs: dst.AssignStmtDecorations{NodeDecs: dst.NodeDecs{Start: []string{"// TODO: add some code"}}},
					Lhs:  []dst.Expr{&dst.StarExpr{X: dst.NewIdent("reply")}},
					Tok:  token.ASSIGN,
					Rhs: []dst.Expr{
						&dst.CompositeLit{
							Type: &dst.Ident{
								Name: outType,
								Path: filepath.Join(appPkg, currentPath, "proto"),
							},
						},
					},
				},
				&dst.ReturnStmt{
					Results: []dst.Expr{dst.NewIdent("nil")},
				},
			},
		},
		Decs: dst.FuncDeclDecorations{NodeDecs: dst.NodeDecs{Before: dst.EmptyLine}},
	}
}
