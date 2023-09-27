package ripple

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
	"reflect"
	"strings"
)

// fieldTagKey is the field tag key for ripple
const fieldTagKey = "controller"

// Controller is the interface for a Controller to be applied to an echo Group
type Controller interface {
	// Path is the namespace ripple will create the Group at, eg /posts
	Path() string
}

// AddControllers applies the Controller to the echo via a new Group using the
// Controller's ripple tags as a manifest to properly associate methods/path and
// handler.
func AddController(echoMux *echo.Echo, c Controller) {
	ctlValue, ctlType, err := reflectCtrl(c)
	if err != nil {
		panic(err)
	}

	grp := echoMux.Group(c.Path())

	i := 0
	n := ctlType.NumField()
	for ; i < n; i++ {
		res, err := newResource(ctlType.Field(i), ctlValue)
		if err != nil {
			panic(err)
		}
		if res == nil {
			continue // if there is no route
		}

		res.Set(grp, c.Path())
	}
}

func reflectCtrl(c Controller) (reflect.Value, reflect.Type, error) {
	ctlValue := reflect.ValueOf(c)
	ctlType := ctlValue.Type()

	if ctlType.Kind() == reflect.Ptr {
		ctlValue = ctlValue.Elem()
		ctlType = ctlValue.Type()
	}

	var err error
	if ctlType.Kind() != reflect.Struct {
		err = errNotStruct
	}

	fmt.Println(fmt.Sprintf("%s: %v, Path: %s", color.Bold("[RegisterController]"), color.Bold(color.Blue(ctlType)), color.Bold(color.Yellow(c.Path()))))
	return ctlValue, ctlType, err
}

var errNotStruct = errors.New("invalid controller type: requires a struct type")

//=============================Field===============================
// echoType represents one of the 2 types that can be mounted onto an Echo Group
// either an Handler or a Middlerware
type echoType int

const (
	_ echoType = iota

	middleware
	handler
)

// structFielder is the basic interface we need for a struct field
type structFielder interface {
	Tag() string
	Name() string
	Type() reflect.Type
}

// fieldInfo is the basic meta data parsed from a struct field. This does not
// include the actual field value or the <name>Func method it represents.
type fieldInfo struct {
	Path   string
	Method string
	Name   string
	Type   reflect.Type

	// EchoType represents the type of field in relation to echo either a handler
	// or middleware
	EchoType echoType
}

func newFieldInfo(f structFielder) (*fieldInfo, error) {
	tagInf, err := parseTag(f.Tag())
	if err != nil {
		return nil, err
	}
	if tagInf == nil {
		return nil, nil
	}

	return &fieldInfo{
		Method: tagInf.meth,
		Path:   strings.TrimRight(tagInf.path, "/"),
		Name:   f.Name(),
		Type:   f.Type(),

		EchoType: tagInf.EchoType,
	}, nil
}

// MethodName returns the associated method name for ripple field.
// eg. Index -> ActionIndex
func (f fieldInfo) MethodName() string {
	return fmt.Sprintf("Action%s", f.Name)
}

// methodMap maps all echo methods that match the func(string, echo.Handler)
// signature used to add method routes
var methodMap = map[string]string{
	"GET":     "GET",
	"POST":    "POST",
	"PUT":     "PUT",
	"PATCH":   "PATCH",
	"DELETE":  "DELETE",
	"HEAD":    "HEAD",
	"OPTIONS": "OPTIONS",
	"CONNECT": "CONNECT",
	"TRACE":   "TRACE",
}

// tagInfo represents the decoded tag string
type tagInfo struct {
	meth, path string

	EchoType echoType
}

func parseTag(tag string) (*tagInfo, error) {
	if tag == "" {
		return nil, nil
	}
	if tag == ",middleware" {
		return &tagInfo{EchoType: middleware}, nil
	}

	arr := strings.Split(tag, " ")
	if len(arr) != 2 {
		return nil, errTagFormat
	}

	meth := arr[0]
	path := arr[1]

	_, ok := methodMap[meth]
	if !ok {
		return nil, errHttpMethod(meth)
	}

	return &tagInfo{meth, path, handler}, nil
}

var errTagFormat = errors.New("invalid tag format")

type errHttpMethod string

func (e errHttpMethod) Error() string {
	return fmt.Sprintf("invalid HTTP method: %s", string(e))
}

//===========================Resource=================================
// resource represents the handler/middleware to be mounted onto an Echo Group
type resource struct {
	*fieldInfo

	Func reflect.Value
}

func newResource(f reflect.StructField, v reflect.Value) (*resource, error) {
	fieldinf, err := newFieldInfo(structField{f})
	if err != nil {
		return nil, err
	}
	if fieldinf == nil {
		return nil, nil
	}

	fn, err := getResourceFunc(fieldinf, v)
	if err != nil {
		return nil, err
	}
	if !fn.Type().ConvertibleTo(fieldinf.Type) {
		return nil, errTypeMismatch
	}

	return &resource{
		fieldInfo: fieldinf,

		Func: fn,
	}, nil
}

func (r resource) isMiddleware() bool {
	return r.EchoType == middleware
}

func (r resource) callName() string {
	if r.isMiddleware() {
		return "Use"
	}

	return methodMap[r.Method]
}

func (r resource) callArgs() []reflect.Value {
	if r.isMiddleware() {
		return []reflect.Value{r.Func}
	}

	return []reflect.Value{
		reflect.ValueOf(r.Path),
		r.Func,
	}
}

// Set sets the resources on the given group
func (r resource) Set(grp *echo.Group, path string) {
	fmt.Println(fmt.Sprintf("%s method: %v, path: %s, action: %v",
		color.Bold(color.Yellow("[router]")),
		color.Green(r.callName()),
		color.Red(path+r.fieldInfo.Path),
		color.Yellow("Action"+r.fieldInfo.Name)))
	reflect.ValueOf(grp).MethodByName(r.callName()).Call(r.callArgs())
}

// structField is a wrapper that implements structFielder
type structField struct {
	field reflect.StructField
}

func (f structField) Tag() string {
	return f.field.Tag.Get(fieldTagKey)
}

func (f structField) Name() string {
	return f.field.Name
}

func (f structField) Type() reflect.Type {
	return f.field.Type
}

var errTypeMismatch = errors.New("field and method types do not match")

// getResourceFunc returns the associated <name>Func method for a defined ripple
// field or the actual field value if the <name>Func association is not found.
func getResourceFunc(
	fieldinf *fieldInfo, v reflect.Value) (reflect.Value, error) {

	var fn reflect.Value

	// first search methods
	fn = v.MethodByName(fieldinf.MethodName())
	if fn.IsValid() {
		return fn, nil
	}

	// then search fields
	fn = v.FieldByName(fieldinf.Name)
	if fn.IsValid() && !reflect.ValueOf(fn.Interface()).IsNil() {
		return fn, nil
	}

	return fn, errActionNotFound(fieldinf.Name)
}

type errActionNotFound string

func (e errActionNotFound) Error() string {
	return fmt.Sprintf("action not found: %s", string(e))
}

//===========================Namespace=================================
// Namespace provides an embeddable type that will allow a struct to implement
// Controller.
type Namespace string

var _ Controller = Namespace("")

// Path returns a string implementing Controller
func (n Namespace) Path() string {
	if n == "" {
		return "/"
	}

	return string(n)
}
