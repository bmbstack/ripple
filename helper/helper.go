package helper

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
)

const (
	DateShortLayout  = "2006-01-02"
	DateFullLayout   = "2006-01-02 15:04:05"
	TimeLocationName = "Asia/Shanghai"
)

// GetTypeName returns a string representing the name of the object typ.
// if the name is defined then it is used, otherwise, the name is derived from the
// Stringer interface.
//
// the stringer returns something like *somepkg.MyStruct, so skip
// the *somepkg and return MyStruct
func GetTypeName(typ reflect.Type) string {
	if typ.Name() != "" {
		return typ.Name()
	}
	split := strings.Split(typ.String(), ".")
	return split[len(split)-1]
}

func IsEmpty(i interface{}) bool {
	return isEmpty(reflect.ValueOf(i))
}

func IsNotEmpty(i interface{}) bool {
	return !isEmpty(reflect.ValueOf(i))
}

func isEmpty(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool() == false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return v.Float() == 0

	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0

	case reflect.Ptr, reflect.Interface:
		return isEmpty(v.Elem())

	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isEmpty(v.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Slice, reflect.String, reflect.Map:
		return v.Len() == 0

	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			if !isEmpty(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return v.IsNil()
	}
}

func Substring(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func CurrentMethodName() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(1, pc)
	f := runtime.FuncForPC(pc[1])
	arr := strings.Split(f.Name(), ".")
	return arr[len(arr)-1]
}

// OnStop 退出信号拦截
// 在POSIX.1-1990标准中定义的信号列表
//
//信号	    值	    	动作			说明
//SIGHUP	1			Term	终端控制进程结束(终端连接断开)
//SIGINT	2			Term	用户发送INTR字符(Ctrl+C)触发
//SIGQUIT	3			Core	用户发送QUIT字符(Ctrl+/)触发
//SIGILL	4			Core	非法指令(程序错误、试图执行数据段、栈溢出等)
//SIGABRT	6			Core	调用abort函数触发
//SIGFPE	8			Core	算术运行错误(浮点运算错误、除数为零等)
//SIGKILL	9			Term	无条件结束程序(不能被捕获、阻塞或忽略)
//SIGSEGV	11			Core	无效内存引用(试图访问不属于自己的内存空间、对只读内存空间进行写操作)
//SIGPIPE	13			Term	消息管道损坏(FIFO/Socket通信时，管道未打开而进行写操作)
//SIGALRM	14			Term	时钟定时信号
//SIGTERM	15			Term	结束程序(可以被捕获、阻塞或忽略)
//SIGUSR1	30,10,16	Term	用户保留
//SIGUSR2	31,12,17	Term	用户保留
//SIGCHLD	20,17,18	Ign		子进程结束(由父进程接收)
//SIGCONT	19,18,25	Cont	继续执行已经停止的进程(不能被阻塞)
//SIGSTOP	17,19,23	Stop	停止进程(不能被捕获、阻塞或忽略)
//SIGTSTP	18,20,24	Stop	停止进程(可以被捕获、阻塞或忽略)
//SIGTTIN	21,21,26	Stop	后台程序从终端中读取数据时触发
//SIGTTOU	22,22,27	Stop	后台程序向终端中写数据时触发
func OnStop(clean func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range signalChan {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program exit signal...", s)
				clean()
				os.Exit(0)
			default:
				fmt.Println("Other signal", s)
			}
		}
	}()
}
