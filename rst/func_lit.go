package rst

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

// AddStmtToFuncLitBody add statement to anonymous function body
func AddStmtToFuncLitBody(df *dst.File, scope Scope, stmt dst.Stmt, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncLit:
			if !scope.IsInScope() {
				return true
			}
			nn := node.(*dst.FuncLit)
			stmtList := nn.Body.List
			pos = normalizePos(pos, len(stmtList))

			nn.Body.List = append(
				stmtList[:pos],
				append([]dst.Stmt{dst.Clone(stmt).(dst.Stmt)}, stmtList[pos:]...)...)
			modified = true
		default:
			scope.TryEnterScope(node)
		}
		return true
	}

	post := func(c *dstutil.Cursor) bool {
		scope.TryLeaveScope(c.Node())
		return true
	}

	dstutil.Apply(df, pre, post)
	return
}

// AddFieldToFuncLitParams add statement to anonymous function params
func AddFieldToFuncLitParams(df *dst.File, scope Scope, field *dst.Field, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncLit:
			if !scope.IsInScope() {
				return true
			}
			nn := node.(*dst.FuncLit)

			fieldList := nn.Type.Params.List
			pos = normalizePos(pos, len(fieldList))
			nn.Type.Params.List = append(
				fieldList[:pos],
				append([]*dst.Field{dst.Clone(field).(*dst.Field)}, fieldList[pos:]...)...)
			modified = true
		default:
			scope.TryEnterScope(node)
		}
		return true
	}

	post := func(c *dstutil.Cursor) bool {
		scope.TryLeaveScope(c.Node())
		return true
	}

	dstutil.Apply(df, pre, post)
	return
}
