package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

func With(fields map[string]interface{}) *logrus.Entry {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	_, file, line, _ := runtime.Caller(1)
	fields["caller"] = fmt.Sprintf("%s:%d", file, line)
	return logrus.WithFields(fields)
}
