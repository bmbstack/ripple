package scripts

import (
	"fmt"
	"github.com/bmbstack/ripple/cmd/ripple/logger"
	"github.com/bmbstack/ripple/cmd/ripple/util"
	"github.com/bmbstack/ripple/rst"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator/resolver/guess"
	"go/token"
	"golang.org/x/mod/modfile"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func Generate(currentPath, component, name string) {
	appPkg, err := getAppPkg()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("parse go.mod err: %v", err))
		logger.Logger.Error("Please execute the command `ripple g` under the same level directory of go.mod")
		return
	}
	logger.Logger.Notice(fmt.Sprintf("the project package: %s", appPkg))

	if strings.HasSuffix(currentPath, ".proto") {
		logger.Logger.Debug(fmt.Sprintf("[%s]generate code *.pb.go, *.rpc.go, ref file: %s ...", "singleFile", currentPath))
		goPath := getGoPath()
		generateOnePb(goPath, currentPath)
		return
	} else if strings.HasSuffix(currentPath, ".dto.go") {
		logger.Logger.Debug(fmt.Sprintf("[%s]generate code *.controller.go, *.service.go, ref file: %s ...", "singleFile", currentPath))
		item := fileInfo{Name: getFileNameByPath(currentPath), Path: currentPath}
		if !strings.Contains(currentPath, "internal") {
			logger.Logger.Error("please put your dto file into the internal directory")
			return
		}
		currentPath = currentPath[0:strings.Index(currentPath, "internal")]
		generateOneController(currentPath, appPkg, item.Name, item.Path)
		generateOneService(currentPath, item.Name, item.Path)
		return
	}

	if strings.EqualFold(component, "all") {
		logger.Logger.Debug(fmt.Sprintf("[%s]generate code *.pb.go, *.rpc.go, *.controller.go, *.service.go ...", "all"))
		generatePb(currentPath)
		generateController(currentPath)
		generateService(currentPath)
	} else if strings.EqualFold(component, "proto") {
		logger.Logger.Debug(fmt.Sprintf("[%s]generate code *.pb.go, *.rpc.go ...", component))
		generatePb(currentPath)
	} else if strings.EqualFold(component, "controller") {
		logger.Logger.Debug(fmt.Sprintf("[%s]generate code *.controller.go ...", component))
		generateController(currentPath)
	} else if strings.EqualFold(component, "service") {
		logger.Logger.Debug(fmt.Sprintf("[%s]generate code *.service.go ...", component))
		if strings.EqualFold(name, "") {
			generateService(currentPath)
		} else {
			logger.Logger.Info(fmt.Sprintf("auto generate %s.service.go, ref file: %s", name, "no ref file, auto generate by component name"))
			generateOneServiceByName(currentPath, name, "no ref file, auto generate by component name")
		}
	} else if strings.EqualFold(component, "ecode") {
		logger.Logger.Debug(fmt.Sprintf("[%s]generate ecode.go ...", component))
		generateEcode(currentPath)
	} else if strings.EqualFold(component, "rpc.client") {
		logger.Logger.Debug(fmt.Sprintf("[%s]generate code rpc.client.go ...", component))
		if util.IsEmpty(name) {
			logger.Logger.Error("Please input the source *.pb.go file")
			return
		}
		if !strings.HasSuffix(name, ".pb.go") {
			logger.Logger.Error("Please input name:  *.pb.go file")
			return
		}
		generateRpcClient(currentPath, name)
	} else {
		logger.Logger.Error("Please input component name: '', proto, controller, service")
	}

	logger.Logger.Notice("auto generate file finish")
}

func generatePb(currentPath string) {
	goPathArray := strings.Split(os.Getenv("GOPATH"), ":")
	goPath := goPathArray[0]

	list := collect(currentPath, ".proto")
	for _, item := range list {
		generateOnePb(goPath, item.Path)
	}
}

func generateOnePb(goPath, itemPath string) {
	logger.Logger.Info(fmt.Sprintf("auto generate *.pb.go, ref file: %s", itemPath))
	cmd := fmt.Sprintf("protoc -I.:%s/src --gofast_out=plugins=ripple:. %s", goPath, itemPath)
	logger.Logger.Debug(fmt.Sprintf("Run command: %s", cmd))
	out, err := RunCommand("bash", "-c", cmd)
	logger.Logger.Info(fmt.Sprintf("protoc gen out: %s", string(out)))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("protoc gen err: %s", err.Error()))
	} else {
		logger.Logger.Notice(fmt.Sprintf("ref file: %s, generate *.pb.go, *.rpc.go success", itemPath))
	}
}

func generateController(currentPath string) {
	appPkg, err := getAppPkg()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("parse go.mod err: %v", err))
		logger.Logger.Error("Please execute the command `ripple g` under the same level directory of go.mod")
		return
	}

	list := collect(currentPath, ".dto.go")
	for _, item := range list {
		generateOneController(currentPath, appPkg, item.Name, item.Path)
	}
}

func generateOneController(currentPath, appPkg, itemName, itemPath string) {
	nameArr := strings.Split(itemName, ".")
	if !strings.HasPrefix(itemName, ".") && len(nameArr) > 0 {
		module := nameArr[0]
		source := itemPath

		logger.Logger.Info(fmt.Sprintf("auto parse %s annotation according to the ast", source))

		pkgNames := map[string]string{"github.com/labstack/echo/v4": "echo"}
		resolver := guess.WithMap(pkgNames)
		dfs, _ := rst.ParseSrcFile(source, resolver)

		routers := parseRouterFromAnnot(dfs, source)
		if util.IsEmpty(routers) {
			return
		}

		// create ecode.go
		generateEcode(currentPath)

		// create *.controller.go
		logger.Logger.Info(fmt.Sprintf("auto generate %s.controller.go, ref file: %s", module, source))
		pkg := "v1"
		dir := path.Join(currentPath, "internal", "controllers", pkg)
		err := os.MkdirAll(dir, Permissions)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("The project path could not be created: %s", err))
			return
		}
		target := filepath.Join(dir, fmt.Sprintf("%s.controller.go", module))
		logger.Logger.Info(fmt.Sprintf("check: %s, file exist: %t", target, util.Exist(target)))

		upper := util.StartToUpper(module)
		if !util.Exist(target) {
			// create file, write code to file
			_, err := os.Create(target)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("%s not be created: %v", target, err))
				return
			}
			logger.Logger.Info(fmt.Sprintf("file: %s, create success", target))

			code := fmt.Sprintf(`// Code generated by ripple g, You can edit it again.
// source: %s

package %s

import (
	"github.com/labstack/echo/v4"
)

type %[3]sController struct {
	Group *echo.Group
	BaseController
}

func (this %[3]sController) Setup() {}
		`, source, pkg, upper)
			err = ioutil.WriteFile(target, []byte(code), Permissions)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("%s not be write: %v", target, err))
				return
			}
		}

		// convert to ast, append action && setup uri
		df, err := rst.ParseSrcFile(target, resolver)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("file: %s, file syntax error, convert to ast error: %v", target, err))
			return
		}
		logger.Logger.Info(fmt.Sprintf("file: %s, convert to ast success", target))

		for _, item := range routers {
			// Action
			hasAction := rst.HasFuncDeclWithRecvInFile(df, dst.FuncDecl{
				Name: &dst.Ident{Name: item.Action},
			}, fmt.Sprintf("%sController", upper))
			hasResp := rst.HasStructDeclInFile(dfs, item.RespIdent)
			if !hasResp {
				logger.Logger.Error(fmt.Sprintf("resp struct is error, you must have the corresponding Resp struct: %s, ref file: %s", item.RespIdent, source))
			}
			if !hasAction && hasResp {
				fd := createActionFunc(appPkg, currentPath, upper, item)
				df.Decls = append(df.Decls, fd)
			}

			// Setup
			routeStmt := &dst.ExprStmt{
				X: &dst.CallExpr{
					Fun: &dst.SelectorExpr{
						X: &dst.SelectorExpr{
							X:   dst.NewIdent("this"),
							Sel: dst.NewIdent("Group"),
						},
						Sel: dst.NewIdent(item.Method),
					},
					Args: []dst.Expr{
						&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(item.Uri)},
						&dst.SelectorExpr{X: dst.NewIdent("this"), Sel: dst.NewIdent(item.Action)},
					},
				},
			}

			arg := &dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(item.Uri)}
			sign := &dst.SelectorExpr{X: dst.NewIdent("this"), Sel: dst.NewIdent(item.Action)}

			var hasSign bool
			var existsMethod string
			methods := []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodConnect,
				http.MethodOptions,
				http.MethodTrace,
			}
			for _, m := range methods {
				if rst.HasArgInCallExpr(df, rst.Scope{FuncName: "Setup"}, m, sign) {
					existsMethod = m
					hasSign = true
					break
				}
			}
			if hasSign {
				rst.ReplaceFuncNameAndArgWithIndexInCallExpr(df, rst.Scope{FuncName: "Setup"}, existsMethod, sign, item.Method, 0, arg)
			} else {
				rst.AddStmtToFuncBodyEndWithRecv(df, fmt.Sprintf("%sController", upper), "Setup", routeStmt)
			}
		}

		buf := rst.PrintToBuf(df, resolver, []rst.PkgAlias{})
		err = ioutil.WriteFile(target, buf.Bytes(), Permissions)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("%s not be write: %v", target, err))
		} else {
			logger.Logger.Notice(fmt.Sprintf("file: %s, generate success", target))
		}
	}
}

func generateService(currentPath string) {
	list := collect(currentPath, ".dto.go")
	for _, item := range list {
		generateOneService(currentPath, item.Name, item.Path)
	}
}

func generateOneService(currentPath, itemName, itemPath string) {
	nameArr := strings.Split(itemName, ".")
	if !strings.HasPrefix(itemName, ".") && len(nameArr) > 0 {
		module := nameArr[0]
		source := itemPath

		pkgNames := map[string]string{"github.com/labstack/echo/v4": "echo"}
		resolver := guess.WithMap(pkgNames)
		dfs, _ := rst.ParseSrcFile(source, resolver)

		routers := parseRouterFromAnnot(dfs, source)
		if util.IsEmpty(routers) {
			return
		}

		logger.Logger.Info(fmt.Sprintf("auto generate %s.service.go, ref file: %s", module, source))
		generateOneServiceByName(currentPath, module, source)
	}
}

func generateOneServiceByName(currentPath, module, source string) {
	pkg := "services"
	dir := path.Join(currentPath, "internal", pkg)
	err := os.MkdirAll(dir, Permissions)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("The project path could not be created: %s", err))
	}
	target := filepath.Join(dir, fmt.Sprintf("%s.service.go", module))
	logger.Logger.Info(fmt.Sprintf("check: %s, file exist: %t", target, util.Exist(target)))

	lower := util.StartToLower(module)
	upper := util.StartToUpper(module)
	if !util.Exist(target) {
		// create file, write code to file
		_, err := os.Create(target)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("%s not be created: %v", target, err))
			return
		}
		logger.Logger.Info(fmt.Sprintf("file: %s, create success", target))

		code := fmt.Sprintf(`// Code generated by ripple g, You can edit it again.
// source: %s

package %s

import (
	"github.com/labstack/echo/v4"
	"sync"
)

var (
	%[3]sService     *%[4]sService
	%[3]sServiceOnce sync.Once
)

// Get%[4]sService 单例
func Get%[4]sService(ctx echo.Context) *%[4]sService {
	%[3]sServiceOnce.Do(func() {
		%[3]sService = &%[4]sService{}
	})
	%[3]sService.ctx = ctx
	return %[3]sService
}

type %[4]sService struct {
	ctx echo.Context
}
		`, source, pkg, lower, upper)

		err = ioutil.WriteFile(target, []byte(code), Permissions)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("%s not be write: %v", target, err))
		} else {
			logger.Logger.Notice(fmt.Sprintf("file: %s, generate success", target))
		}
	}
}

func generateEcode(currentPath string) {
	logger.Logger.Info("auto generate internal/ecode/ecode.go")

	pkg := "ecode"
	dir := path.Join(currentPath, "internal", pkg)
	err := os.MkdirAll(dir, Permissions)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("The project path could not be created: %s", err))
	}
	target := filepath.Join(dir, "ecode.go")
	logger.Logger.Info(fmt.Sprintf("check: %s, file exist: %t", target, util.Exist(target)))

	if !util.Exist(target) {
		// create file, write code to file
		_, err := os.Create(target)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("%s not be created: %v", target, err))
			return
		}
		logger.Logger.Info(fmt.Sprintf("file: %s, create success", target))

		code := fmt.Sprintf(`// Code generated by ripple g, You can edit it again.

package %s

import (
	"github.com/labstack/echo/v4"
	"net/http"
)


//==============================================
//                  code
//==============================================

var (
	Success     = value(200, "成功")
	ServerError = value(500, "服务器错误")
	ParamError  = value(501, "参数错误")
)

//==============================================
//                  OK
//==============================================

func OK(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": Success.Code,
		"msg":  Success.Msg,
		"data": data,
	})
}

//==============================================
//                  Error
//==============================================

func Error(ctx echo.Context, err error) error {
	ec, ok := err.(Ecode)
	if !ok {
		ec = ServerError
		ec.Msg = err.Error()
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": ec.Code,
		"msg":  ec.Msg,
		"data": []map[string]interface{}{},
	})
}

func ErrorWithHttpCode(ctx echo.Context, httpCode int, err error) error {
	ec, ok := err.(Ecode)
	if !ok {
		ec = ServerError
		ec.Msg = err.Error()
	}
	return ctx.JSON(httpCode, map[string]interface{}{
		"code": ec.Code,
		"msg":  ec.Msg,
		"data": []map[string]interface{}{},
	})
}

//==============================================
//                  Ecode
//==============================================

type Ecode struct {
	Code int64
	Msg  string
}

func value(code int64, msg string) Ecode {
	return Ecode{Code: code, Msg: msg}
}

// Status implement error interface
func (this Ecode) Error() string {
	return this.Msg
}
		`, pkg)

		err = ioutil.WriteFile(target, []byte(code), Permissions)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("%s not be write: %v", target, err))
		} else {
			logger.Logger.Notice(fmt.Sprintf("file: %s, generate success", target))
		}
	}
}

func generateRpcClient(currentPath, name string) {
	logger.Logger.Info("auto generate internal/service/rpc.client.go")

	pkg := "services"
	dir := path.Join(currentPath, "internal", pkg)
	err := os.MkdirAll(dir, Permissions)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("The project path could not be created: %s", err))
	}
	target := filepath.Join(dir, "rpc.client.go")
	logger.Logger.Info(fmt.Sprintf("check: %s, file exist: %t", target, util.Exist(target)))

	if !util.Exist(name) {
		logger.Logger.Error(fmt.Sprintf("check: %s, file exist: %t", name, util.Exist(name)))
		return
	}

	dfs, err := rst.ParseSrcFile(name, guess.New())
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("file: %s, file syntax error, convert to ast error: %v", name, err))
		return
	}

	funcList := rst.GetFuncDeclListWithPrefixSuffixInFile(dfs, "New", "Client")
	if util.IsEmpty(funcList) {
		logger.Logger.Notice("*.pb.go doesn't have NewClient func")
		return
	}

	appPkg, err := getAppPkg()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("parse go.mod err: %v", err))
		logger.Logger.Error("Please execute the command `ripple g` under the same level directory of go.mod")
		return
	}
	nameArr := strings.Split(name, "/")
	pbArr := nameArr[0 : len(nameArr)-1]
	proto := pbArr[len(pbArr)-1:][0]
	pbPkg := filepath.Join(appPkg, strings.Join(pbArr, "/"))

	if !util.Exist(target) {
		// create file, write code to file
		_, err := os.Create(target)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("%s not be created: %v", target, err))
			return
		}
		logger.Logger.Info(fmt.Sprintf("file: %s, create success", target))

		var stmtVar string
		var codeFunc string
		var closeFuncCall string
		for key, item := range funcList {
			itemValue := strings.ReplaceAll(item, "New", "")
			stmtVar += fmt.Sprintf(`    %s     *%s.%s
`, util.StartToLower(itemValue), proto, itemValue)
			if key < (len(funcList) - 1) {
				stmtVar += fmt.Sprintf(`    %sOnce sync.Once

`, util.StartToLower(itemValue))
			} else {
				stmtVar += fmt.Sprintf(`    %sOnce sync.Once`, util.StartToLower(itemValue))
			}

			codeFunc += fmt.Sprintf(`
func Get%[1]s() *proto.%[1]s {
	%[2]sOnce.Do(func() {
		close%[1]s()
		%[2]s = proto.New%[1]s(func() {
			%[2]sOnce.Reset()
		})
	})
	return %[2]s
}
`, itemValue, util.StartToLower(itemValue))

			codeFunc += fmt.Sprintf(`
func close%[1]s() {
	if %[2]s != nil {
		%[2]s.Discovery.Close()
		%[2]s.XClient.Close()
	}
}
`, itemValue, util.StartToLower(itemValue))

			if key < (len(funcList) - 1) {
				closeFuncCall += fmt.Sprintf(`	close%s()
`, itemValue)
			} else {
				closeFuncCall += fmt.Sprintf(`	close%s()`, itemValue)
			}
		}

		codeFunc += fmt.Sprintf(`
func CloseRpcClients() {
%s
}
`, closeFuncCall)

		codeVar := fmt.Sprintf(`
var (
%s
)`, stmtVar)
		code := fmt.Sprintf(`// Code generated by ripple g, DO NOT EDIT.

package %s

import (
	"github.com/bmbstack/ripple/sync"

	"%s"
)
%s
%s`, pkg, pbPkg, codeVar, codeFunc)

		err = ioutil.WriteFile(target, []byte(code), Permissions)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("%s not be write: %v", target, err))
		} else {
			logger.Logger.Notice(fmt.Sprintf("file: %s, generate success", target))
		}
		return
	}

	// convert to ast, append action && setup uri
	df, err := rst.ParseSrcFile(target, guess.New())
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("file: %s, file syntax error, convert to ast error: %v", target, err))
		return
	}
	logger.Logger.Info(fmt.Sprintf("file: %s, convert to ast success", target))

	for _, item := range funcList {
		itemValue := strings.ReplaceAll(item, "New", "")

		spec1 := &dst.ValueSpec{
			Decs:  dst.ValueSpecDecorations{NodeDecs: dst.NodeDecs{Before: dst.EmptyLine}},
			Names: []*dst.Ident{dst.NewIdent(util.StartToLower(itemValue))},
			Type:  &dst.StarExpr{X: &dst.Ident{Name: itemValue, Path: pbPkg}},
		}
		hasSpec1 := rst.HasVarInBlockGenDecl(df, 0, spec1)
		if !hasSpec1 {
			rst.AddVarIntoBlockGenDecl(df, 0, spec1)
		}

		spec2 := &dst.ValueSpec{
			Names: []*dst.Ident{dst.NewIdent(fmt.Sprintf("%sOnce", util.StartToLower(itemValue)))},
			Type:  &dst.Ident{Name: "Once", Path: "github.com/bmbstack/ripple/sync"},
		}
		hasSpec2 := rst.HasVarInBlockGenDecl(df, 0, spec2)
		if !hasSpec2 {
			rst.AddVarIntoBlockGenDecl(df, 0, spec2)
		}

		hasFunc1 := rst.HasFuncDeclInFile(df, dst.FuncDecl{Name: &dst.Ident{Name: fmt.Sprintf("Get%s", itemValue)}})
		if !hasFunc1 {
			fd := createGetRpcClientFunc(itemValue, util.StartToLower(itemValue), pbPkg)
			df.Decls = append(df.Decls, fd)
		}

		hasFunc2 := rst.HasFuncDeclInFile(df, dst.FuncDecl{Name: &dst.Ident{Name: fmt.Sprintf("close%s", itemValue)}})
		if !hasFunc2 {
			fd := createCloseOneRpcClientFunc(itemValue, util.StartToLower(itemValue), pbPkg)
			df.Decls = append(df.Decls, fd)
		}
	}

	rst.DeleteFuncFromFile(df, "CloseRpcClients")
	fd := createCloseAllRpcClientFunc()
	df.Decls = append(df.Decls, fd)

	closeFuncList := rst.GetFuncDeclListWithPrefixSuffixInFile(df, "close", "Client")
	for _, item := range closeFuncList {
		closeOneCallStmt := &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Name: fmt.Sprintf("%s", item)}}}
		rst.AddStmtToFuncBodyEnd(df, "CloseRpcClients", closeOneCallStmt)
	}

	buf := rst.PrintToBuf(df, guess.New(), []rst.PkgAlias{})
	err = ioutil.WriteFile(target, buf.Bytes(), Permissions)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s not be write: %v", target, err))
	} else {
		logger.Logger.Notice(fmt.Sprintf("file: %s, generate success", target))
	}
}

//=================================================================
//                        common
//=================================================================

type fileInfo struct {
	Name string
	Path string
}

type routerInfo struct {
	ReqIdent  string
	RespIdent string

	Action string
	Uri    string
	Method string
}

// collect files with suffix in current path
func collect(currentPath string, suffix string) (list []fileInfo) {
	err := filepath.Walk(currentPath, func(file string, info fs.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return err
		}
		if strings.HasSuffix(file, suffix) {
			list = append(list, fileInfo{
				Name: info.Name(),
				Path: file,
			})
		}
		return nil
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("filepath.Walk err: %v", err))
	}
	return list
}

func parseKey(key, line string) (uri string) {
	arr := strings.Split(line, " ")
	if len(arr) >= 2 {
		if strings.EqualFold(key, arr[0]) {
			uri = arr[1]
		}
	}
	return uri
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

func getGoPath() string {
	goPathArray := strings.Split(os.Getenv("GOPATH"), ":")
	if len(goPathArray) > 0 {
		return goPathArray[0]
	}
	return ""
}

func getFileNameByPath(path string) string {
	arr := strings.Split(path, "/")
	if len(arr) >= 1 {
		return arr[len(arr)-1]
	}
	return ""
}

func parseRouterFromAnnot(dfs *dst.File, source string) []routerInfo {
	ret := rst.GetStructDecsInStructComment(dfs, "@RippleApi")
	var list []routerInfo
	for _, item := range ret {
		var uri string
		var method string
		for _, line := range item.Decs {
			line := strings.Trim(rst.TrimAnnot(line), " ")
			if strings.HasPrefix(line, "@Uri") {
				uri = parseKey("@Uri", line)
			} else if strings.HasPrefix(line, "@Method") {
				method = parseKey("@Method", line)
			}
		}
		if !strings.HasPrefix(item.Name, "Req") {
			logger.Logger.Error(fmt.Sprintf("req struct is error, please has prefix `Req`, ref file: %s", source))
			continue
		}
		name := strings.ReplaceAll(item.Name, "Req", "")
		list = append(list, routerInfo{
			ReqIdent:  fmt.Sprintf("Req%s", name),
			RespIdent: fmt.Sprintf("Resp%s", name),
			Action:    fmt.Sprintf("Action%s", name),
			Uri:       uri,
			Method:    method,
		})
	}
	return list
}

func createActionFunc(appPkg, currentPath, moduleUpper string, router routerInfo) *dst.FuncDecl {
	return &dst.FuncDecl{
		Name: dst.NewIdent(router.Action),
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{Names: []*dst.Ident{dst.NewIdent("this")}, Type: dst.NewIdent(fmt.Sprintf("%sController", moduleUpper))},
			},
		},
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("ctx")},
						Type: &dst.SelectorExpr{
							X:   &dst.Ident{Name: "echo"},
							Sel: &dst.Ident{Name: "Context"},
						},
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: dst.NewIdent("error")},
				},
			},
		},
		Body: &dst.BlockStmt{
			Decs: dst.BlockStmtDecorations{Lbrace: dst.Decorations{"\n"}},
			List: []dst.Stmt{
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("params")},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{&dst.UnaryExpr{
						Op: token.AND,
						X: &dst.CompositeLit{
							Type: &dst.Ident{
								Name: router.ReqIdent,
								Path: filepath.Join(appPkg, currentPath, "internal", "dto"),
							},
						},
					}},
				},
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("err")},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{
						&dst.CallExpr{
							Fun:  &dst.SelectorExpr{X: dst.NewIdent("ctx"), Sel: dst.NewIdent("Bind")},
							Args: []dst.Expr{dst.NewIdent("params")},
						},
					},
				},
				&dst.IfStmt{
					Decs: dst.IfStmtDecorations{
						NodeDecs: dst.NodeDecs{
							After: dst.EmptyLine,
						},
					},
					Cond: &dst.BinaryExpr{X: dst.NewIdent("err"), Op: token.NEQ, Y: dst.NewIdent("nil")},
					Body: &dst.BlockStmt{
						List: []dst.Stmt{
							&dst.ReturnStmt{
								Results: []dst.Expr{
									&dst.CallExpr{
										Fun: &dst.Ident{
											Name: "Error",
											Path: filepath.Join(appPkg, currentPath, "internal", "ecode"),
										},
										Args: []dst.Expr{
											dst.NewIdent("ctx"),
											&dst.Ident{
												Name: "ParamError",
												Path: filepath.Join(appPkg, currentPath, "internal", "ecode"),
											},
										},
									},
								},
							},
						},
					},
				},
				&dst.AssignStmt{
					Decs: dst.AssignStmtDecorations{NodeDecs: dst.NodeDecs{Start: []string{"// TODO: add some code"}}},
					Lhs:  []dst.Expr{dst.NewIdent("result")},
					Tok:  token.DEFINE,
					Rhs: []dst.Expr{&dst.UnaryExpr{
						Op: token.AND,
						X: &dst.CompositeLit{
							Type: &dst.Ident{
								Name: router.RespIdent,
								Path: filepath.Join(appPkg, currentPath, "internal", "dto"),
							},
						},
					}},
				},
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.CallExpr{
							Fun: &dst.Ident{
								Name: "OK",
								Path: filepath.Join(appPkg, currentPath, "internal", "ecode"),
							},
							Args: []dst.Expr{
								dst.NewIdent("ctx"),
								dst.NewIdent("result"),
							},
						},
					},
				},
			},
		},
		Decs: dst.FuncDeclDecorations{NodeDecs: dst.NodeDecs{Before: dst.EmptyLine}},
	}
}

func createGetRpcClientFunc(upper, lower, pbPkg string) *dst.FuncDecl {
	return &dst.FuncDecl{
		Name: dst.NewIdent(fmt.Sprintf("Get%s", upper)),
		Type: &dst.FuncType{
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: &dst.StarExpr{X: &dst.Ident{Name: upper, Path: pbPkg}}},
				},
			},
		},
		Body: &dst.BlockStmt{
			Decs: dst.BlockStmtDecorations{Lbrace: dst.Decorations{"\n"}},
			List: []dst.Stmt{
				&dst.ExprStmt{X: &dst.CallExpr{
					Fun: &dst.SelectorExpr{
						X:   dst.NewIdent(fmt.Sprintf("%sOnce", lower)),
						Sel: dst.NewIdent("Do"),
					},
					Args: []dst.Expr{
						&dst.FuncLit{
							Type: &dst.FuncType{Func: true, Params: &dst.FieldList{Opening: true, Closing: true}},
							Body: &dst.BlockStmt{
								Decs: dst.BlockStmtDecorations{Lbrace: dst.Decorations{"\n"}},
								List: []dst.Stmt{
									&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Name: fmt.Sprintf("close%s", upper)}}},
									&dst.AssignStmt{
										Lhs: []dst.Expr{dst.NewIdent(lower)},
										Tok: token.ASSIGN,
										Rhs: []dst.Expr{&dst.CallExpr{
											Fun: &dst.Ident{Name: fmt.Sprintf("New%s", upper), Path: pbPkg},
											Args: []dst.Expr{
												&dst.FuncLit{
													Type: &dst.FuncType{Func: true},
													Body: &dst.BlockStmt{
														Decs: dst.BlockStmtDecorations{Lbrace: dst.Decorations{"\n"}},
														List: []dst.Stmt{
															&dst.ExprStmt{
																X: &dst.CallExpr{Fun: &dst.SelectorExpr{X: &dst.Ident{Name: fmt.Sprintf("%sOnce", lower)}, Sel: &dst.Ident{Name: "Reset"}}},
															},
														},
													},
												},
											},
										}},
									},
								},
							},
						},
					},
				}},
				&dst.ReturnStmt{
					Results: []dst.Expr{dst.NewIdent(lower)},
				},
			},
		},
		Decs: dst.FuncDeclDecorations{NodeDecs: dst.NodeDecs{Before: dst.EmptyLine}},
	}
}

func createCloseOneRpcClientFunc(upper, lower, pbPkg string) *dst.FuncDecl {
	return &dst.FuncDecl{
		Name: dst.NewIdent(fmt.Sprintf("close%s", upper)),
		Type: &dst.FuncType{},
		Body: &dst.BlockStmt{
			Decs: dst.BlockStmtDecorations{Lbrace: dst.Decorations{"\n"}},
			List: []dst.Stmt{
				&dst.IfStmt{
					Cond: &dst.BinaryExpr{X: &dst.Ident{Name: lower}, Op: token.NEQ, Y: &dst.Ident{Name: "nil"}},
					Body: &dst.BlockStmt{
						List: []dst.Stmt{
							&dst.ExprStmt{
								X: &dst.CallExpr{
									Fun: &dst.SelectorExpr{
										X:   &dst.SelectorExpr{X: &dst.Ident{Name: lower}, Sel: &dst.Ident{Name: "Discovery"}},
										Sel: &dst.Ident{Name: "Close"},
									},
								},
							},
							&dst.ExprStmt{
								X: &dst.CallExpr{
									Fun: &dst.SelectorExpr{
										X:   &dst.SelectorExpr{X: &dst.Ident{Name: lower}, Sel: &dst.Ident{Name: "XClient"}},
										Sel: &dst.Ident{Name: "Close"},
									},
								},
							},
						},
					},
				},
			},
		},
		Decs: dst.FuncDeclDecorations{NodeDecs: dst.NodeDecs{Before: dst.EmptyLine}},
	}
}

func createCloseAllRpcClientFunc() *dst.FuncDecl {
	return &dst.FuncDecl{
		Name: dst.NewIdent("CloseRpcClients"),
		Type: &dst.FuncType{},
		Body: &dst.BlockStmt{
			Decs: dst.BlockStmtDecorations{Lbrace: dst.Decorations{"\n"}},
			List: []dst.Stmt{

			},
		},
		Decs: dst.FuncDeclDecorations{NodeDecs: dst.NodeDecs{Before: dst.EmptyLine}},
	}
}
