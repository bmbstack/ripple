package client

import (
	"strings"
	"testing"

	model2 "github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/smallnest/rpcx/client"
)

func TestConvertInstances(t *testing.T) {
	instances := []model2.Instance{
		{Ip: "10.0.0.1", Port: 8080, Metadata: map[string]string{"network": "tcp", "weight": "1"}},
		{Ip: "10.0.0.2", Port: 8081, Metadata: map[string]string{"network": "quic", "weight": "2"}},
	}
	pairs := convertInstances(instances, nil)
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0].Key != "tcp@10.0.0.1:8080" {
		t.Errorf("pairs[0].Key = %q, want %q", pairs[0].Key, "tcp@10.0.0.1:8080")
	}
	if pairs[1].Key != "quic@10.0.0.2:8081" {
		t.Errorf("pairs[1].Key = %q, want %q", pairs[1].Key, "quic@10.0.0.2:8081")
	}
}

func TestConvertInstancesWithFilter(t *testing.T) {
	instances := []model2.Instance{
		{Ip: "10.0.0.1", Port: 8080, Metadata: map[string]string{"network": "tcp"}},
		{Ip: "10.0.0.2", Port: 8081, Metadata: map[string]string{"network": "quic"}},
	}
	// 过滤掉 network 为 quic 的实例
	filter := func(pair *client.KVPair) bool {
		return !strings.HasPrefix(pair.Key, "quic")
	}
	pairs := convertInstances(instances, filter)
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair after filter, got %d", len(pairs))
	}
	if pairs[0].Key != "tcp@10.0.0.1:8080" {
		t.Errorf("pairs[0].Key = %q, want tcp@10.0.0.1:8080", pairs[0].Key)
	}
}

func TestConvertInstancesEmpty(t *testing.T) {
	if pairs := convertInstances(nil, nil); len(pairs) != 0 {
		t.Errorf("expected 0 pairs for nil input, got %d", len(pairs))
	}
}

func TestConvertInstancesNilMetadata(t *testing.T) {
	instances := []model2.Instance{
		{Ip: "10.0.0.1", Port: 8080},
	}
	pairs := convertInstances(instances, nil)
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	// Metadata 为 nil 时 network 为空，key 形如 "@ip:port"
	if pairs[0].Key != "@10.0.0.1:8080" {
		t.Errorf("pairs[0].Key = %q, want %q", pairs[0].Key, "@10.0.0.1:8080")
	}
}
