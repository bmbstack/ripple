package rst

import (
	"github.com/dave/dst"
)

type Scope struct {
	FuncName string

	currentScopeNode dst.Node
}

var EmptyScope = Scope{}

func (s Scope) isEmptyScope() bool {
	return s == EmptyScope
}

func (s Scope) IsInScope() bool {
	if s.isEmptyScope() {
		return true
	}
	return s.currentScopeNode != nil
}

func (s *Scope) TryEnterScope(node dst.Node) (ok bool) {
	switch node.(type) {
	case *dst.FuncDecl:
		nn := node.(*dst.FuncDecl)
		if nn.Name.Name == s.FuncName {
			s.currentScopeNode = nn
			ok = true
		}
	}
	return
}

func (s *Scope) TryLeaveScope(node dst.Node) (ok bool) {
	if s.IsInScope() {
		if s.currentScopeNode == node {
			ok = true
			s.currentScopeNode = nil
		}
	}
	return
}
