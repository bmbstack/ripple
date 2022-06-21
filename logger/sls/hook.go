package sls

import (
	"fmt"

	slssdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/sirupsen/logrus"
)

var _ logrus.Hook = (*SLSHook)(nil)

type OptionFunc func(*Option)

type Option struct {
	project  string
	logstore string
	topic    string
	source   string
}

func SetProject(name string) OptionFunc {
	return func(o *Option) {
		o.project = name
	}
}

func SetLogstore(name string) OptionFunc {
	return func(o *Option) {
		o.logstore = name
	}
}

func SetTopic(name string) OptionFunc {
	return func(o *Option) {
		o.topic = name
	}
}

func SetSource(name string) OptionFunc {
	return func(o *Option) {
		o.source = name
	}
}

func NewSLSHook(accessKeyId, accessKeySecret, endpoint, allowLogLevel string, opts ...OptionFunc) *SLSHook {
	opt := &Option{}
	if len(opts) > 0 {
		for _, fun := range opts {
			if fun != nil {
				fun(opt)
			}
		}
	}
	pc := producer.GetDefaultProducerConfig()
	pc.IsJsonType = true
	pc.Endpoint = endpoint
	pc.AccessKeyID = accessKeyId
	pc.AccessKeySecret = accessKeySecret
	pc.AllowLogLevel = allowLogLevel
	p := producer.InitProducer(pc)
	p.Start()
	return &SLSHook{opt, p}
}

type SLSHook struct {
	opt      *Option
	producer *producer.Producer
}

func (hook *SLSHook) Fire(entry *logrus.Entry) error {
	ts := uint32(entry.Time.Unix())

	var contents []*slssdk.LogContent
	for key, value := range entry.Data {
		k := fmt.Sprint(key)
		v := fmt.Sprint(value)
		contents = append(contents, &slssdk.LogContent{
			Key:   &k,
			Value: &v,
		})
	}

	timeKey := "time"
	timeContent := entry.Time.Format("2006-01-02 15:04:05")
	contents = append(contents, &slssdk.LogContent{
		Key:   &timeKey,
		Value: &timeContent,
	})

	levelKey := "level"
	levelContent := entry.Level.String()
	contents = append(contents, &slssdk.LogContent{
		Key:   &levelKey,
		Value: &levelContent,
	})

	msgKey := "message"
	msgContent := entry.Message
	contents = append(contents, &slssdk.LogContent{
		Key:   &msgKey,
		Value: &msgContent,
	})

	log := &slssdk.Log{
		Time:     &ts,
		Contents: contents,
	}

	err := hook.producer.SendLog(hook.opt.project, hook.opt.logstore, hook.opt.topic, hook.opt.source, log)
	if err != nil {
		return err
	}
	return nil
}

func (hook *SLSHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *SLSHook) SafeClose() {
	if hook.producer != nil {
		hook.producer.SafeClose()
	}
}

func (hook *SLSHook) Close(timeoutMs int64) error {
	if hook.producer != nil {
		return hook.producer.Close(timeoutMs)
	}
	return nil
}
