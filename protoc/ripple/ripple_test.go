package ripple

import "testing"

func TestUpperFirstLatter(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"a", "A"},
		{"student", "Student"},
		{"Student", "Student"},
		{"studentLearn", "StudentLearn"},
	}
	for _, c := range cases {
		if got := upperFirstLatter(c.in); got != c.want {
			t.Errorf("upperFirstLatter(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTrimAnnot(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"// @RippleRpc", " @RippleRpc"},
		{"/* @RippleRpc */", " @RippleRpc "},
		{"@NacosGroup DEFAULT_GROUP", "@NacosGroup DEFAULT_GROUP"},
	}
	for _, c := range cases {
		if got := trimAnnot(c.in); got != c.want {
			t.Errorf("trimAnnot(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestHasTag(t *testing.T) {
	if !hasTag("// @RippleRpc", "@RippleRpc") {
		t.Error("hasTag should find @RippleRpc")
	}
	if !hasTag("// @RippleRpc\r\n// @NacosGroup", "@RippleRpc") {
		t.Error("hasTag should find @RippleRpc across newlines")
	}
	if hasTag("// @NacosGroup", "@RippleRpc") {
		t.Error("hasTag should not find @RippleRpc")
	}
}

func TestParseKey(t *testing.T) {
	cases := []struct {
		key, line, want string
	}{
		{"@NacosGroup", "@NacosGroup DEFAULT_GROUP", "DEFAULT_GROUP"},
		{"@NacosCluster", "@NacosCluster ripple", "ripple"},
		{"@NacosGroup", "not-matching", ""},
		{"@NacosGroup", "@NacosGroup", ""}, // 只有一个 token，无值
	}
	for _, c := range cases {
		if got := parseKey(c.key, c.line); got != c.want {
			t.Errorf("parseKey(%q, %q) = %q, want %q", c.key, c.line, got, c.want)
		}
	}
}

// TestNacosSdkPkgPath 验证升级改造后：生成代码引用官方 nacos-group 仓库，
// 而非 vendored 的 bmbstack 路径。这是本次去 vendored 改造的核心约束。
func TestNacosSdkPkgPath(t *testing.T) {
	if got, want := rpcxNacosSdkConsPkgPath, "github.com/nacos-group/nacos-sdk-go/v2/common/constant"; got != want {
		t.Errorf("rpcxNacosSdkConsPkgPath = %q, want %q (生成代码必须引用官方 SDK)", got, want)
	}

	// ripple 自己的 rpcxnacos 插件路径不应变化
	if got, want := rpcxNacosClientPkgPath, "github.com/bmbstack/ripple/nacos/rpcxnacos/client"; got != want {
		t.Errorf("rpcxNacosClientPkgPath = %q, want %q", got, want)
	}
}
