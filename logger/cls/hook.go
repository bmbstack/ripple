package cls

import (
	"fmt"
	"github.com/sirupsen/logrus"
	clssdk "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

var _ logrus.Hook = (*CLSHook)(nil)

type OptionFunc func(*Option)

type Option struct {
	topic string
}

func SetTopic(name string) OptionFunc {
	return func(o *Option) {
		o.topic = name
	}
}

func NewCLSHook(accessKeyId, accessKeySecret, endpoint, allowLogLevel string, opts ...OptionFunc) *CLSHook {
	opt := &Option{}
	if len(opts) > 0 {
		for _, fun := range opts {
			if fun != nil {
				fun(opt)
			}
		}
	}
	pc := clssdk.GetDefaultAsyncProducerClientConfig()
	pc.AccessKeyID = accessKeyId
	pc.AccessKeySecret = accessKeySecret
	pc.Endpoint = endpoint

	p, err := clssdk.NewAsyncProducerClient(pc)
	if err != nil {
		fmt.Println("clssdk.NewAsyncProducerClient err: ", err)
		return &CLSHook{}
	}
	p.Start()
	return &CLSHook{opt, p}
}

type CLSHook struct {
	opt      *Option
	producer *clssdk.AsyncProducerClient
}

func (hook *CLSHook) Fire(entry *logrus.Entry) error {
	var out = map[string]string{
		"time":    entry.Time.Format("2006-01-02 15:04:05"),
		"level":   entry.Level.String(),
		"message": entry.Message,
	}

	for key, value := range entry.Data {
		k := fmt.Sprint(key)
		v := fmt.Sprint(value)
		out[k] = v
	}

	return hook.producer.SendLog(
		hook.opt.topic,
		clssdk.NewCLSLog(
			entry.Time.Unix(),
			out,
		), nil,
	)
}

func (hook *CLSHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *CLSHook) Close(timeoutMs int64) error {
	if hook.producer != nil {
		return hook.producer.Close(timeoutMs)
	}
	return nil
}
