package rst

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

// HasArgInCallExpr checks if the arguments of the function call has given arg
func HasArgInCallExpr(df *dst.File, scope Scope, funcName string, arg dst.Expr) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.CallExpr:
			if !scope.IsInScope() {
				return true
			}

			var found bool
			nn := node.(*dst.CallExpr)
			if ie, ok := nn.Fun.(*dst.Ident); ok && ie.Name == funcName {
				found = true
			}

			if se, ok := nn.Fun.(*dst.SelectorExpr); ok && se.Sel.Name == funcName {
				found = true
			}

			if found {
				for _, cArg := range nn.Args {
					if nodesEqual(arg, cArg) {
						ret = true
					}
				}
				return false
			}
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

// DeleteArgFromCallExpr deletes any arg, in the function call's argument list,
// that is semantically equal to the given arg.
func DeleteArgFromCallExpr(df *dst.File, scope Scope, funcName string, arg dst.Expr) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.CallExpr:
			if !scope.IsInScope() {
				return true
			}

			var found bool
			nn := node.(*dst.CallExpr)
			if ie, ok := nn.Fun.(*dst.Ident); ok && ie.Name == funcName {
				found = true
			}

			if se, ok := nn.Fun.(*dst.SelectorExpr); ok && se.Sel.Name == funcName {
				found = true
			}

			if found {
				var newArgs []dst.Expr
				for _, cArg := range nn.Args {
					if !nodesEqual(arg, cArg) {
						newArgs = append(newArgs, cArg)
					} else {
						modified = true
					}
				}
				nn.Args = newArgs
				return false
			}
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

// AddArgToCallExpr adds given arg, to the function call's argument list, in the given position
func AddArgToCallExpr(df *dst.File, scope Scope, funcName string, arg dst.Expr, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.CallExpr:
			if !scope.IsInScope() {
				return true
			}

			nn := node.(*dst.CallExpr)
			var ce *dst.CallExpr

			if ie, ok := nn.Fun.(*dst.Ident); ok && ie.Name == funcName {
				ce = nn
			}

			if se, ok := nn.Fun.(*dst.SelectorExpr); ok && se.Sel.Name == funcName {
				ce = nn
			}

			if ce != nil {
				args := ce.Args
				pos = normalizePos(pos, len(args))
				ce.Args = append(
					args[:pos],
					append([]dst.Expr{dst.Clone(arg).(dst.Expr)}, args[pos:]...)...)
				modified = true
				return false
			}
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

func SetMethodOnReceiver(df *dst.File, scope Scope, receiver, oldMethod, newMethod string) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.CallExpr:
			if !scope.IsInScope() {
				return true
			}
			nn := node.(*dst.CallExpr)

			switch nn.Fun.(type) {
			case *dst.Ident:
				si := nn.Fun.(*dst.Ident)
				if si.Path == receiver && si.Name == oldMethod {
					si.Name = newMethod
					modified = true
				}
			case *dst.SelectorExpr:
				se := nn.Fun.(*dst.SelectorExpr)

				// TODO: deal with other se.X type
				switch se.X.(type) {
				case *dst.Ident:
					xi := se.X.(*dst.Ident)
					if xi.Name != receiver {
						return true
					}
				}
				if se.Sel.Name == oldMethod {
					se.Sel.Name = newMethod
					modified = true
				}
			}
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
