package naming_grpc

import (
	naming_proxy2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_proxy"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestRedoSubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProxy := naming_proxy2.NewMockINamingProxy(ctrl)
	evListener := NewConnectionEventListener(mockProxy)

	cases := []struct {
		serviceName string
		groupName   string
		clusters    string
	}{
		{"service-a", "group-a", ""},
		{"service-b", "group-b", "cluster-b"},
	}

	for _, v := range cases {
		fullServiceName := util2.GetGroupName(v.serviceName, v.groupName)
		evListener.CacheSubscriberForRedo(fullServiceName, v.clusters)
		mockProxy.EXPECT().Subscribe(v.serviceName, v.groupName, v.clusters)
		evListener.redoSubscribe()
		evListener.RemoveSubscriberForRedo(fullServiceName, v.clusters)
	}
}
