package rst

import (
	"github.com/dave/dst"
)

func exprListsEqual(la, lb []dst.Expr) bool {
	if len(la) != len(lb) {
		return false
	}

	for i := 0; i < len(la); i++ {
		if !nodesEqual(la[i], lb[i]) {
			return false
		}
	}
	return true
}

func fieldListsEqual(la, lb *dst.FieldList) bool {
	if la.NumFields() != lb.NumFields() {
		return false
	}

	for i := 0; i < la.NumFields(); i++ {
		if !nodesEqual(la.List[i], lb.List[i]) {
			return false
		}
	}
	return true
}

func nodeListsEqual(la, lb []dst.Node) bool {
	if len(la) != len(lb) {
		return false
	}

	for i := 0; i < len(la); i++ {
		if !nodesEqual(la[i], lb[i]) {
			return false
		}
	}
	return true
}

func stmtListsEqual(la, lb []dst.Stmt) bool {
	if len(la) != len(lb) {
		return false
	}

	for i := 0; i < len(la); i++ {
		if !nodesEqual(la[i], lb[i]) {
			return false
		}
	}
	return true
}

func nodesEqual(a, b dst.Node) (ret bool) {
	switch a.(type) {
	case *dst.Field:
		na := a.(*dst.Field)
		nb, ok := b.(*dst.Field)

		if len(na.Names) != len(nb.Names) {
			return false
		}

		if len(na.Names) == 0 {
			return ok && nodesEqual(na.Type, nb.Type)
		}
		return ok && nodesEqual(na.Names[0], nb.Names[0]) && nodesEqual(na.Type, nb.Type)

	case *dst.FuncType:
		na := a.(*dst.FuncType)
		nb, ok := b.(*dst.FuncType)

		if (len(na.Params.List) != len(nb.Params.List)) || (len(na.Results.List) != len(nb.Results.List)) {
			return false
		}
		return ok && fieldListsEqual(na.Params, nb.Params) && fieldListsEqual(na.Results, nb.Results)
	case *dst.StarExpr:
		na := a.(*dst.StarExpr)
		nb, ok := b.(*dst.StarExpr)
		return ok && nodesEqual(na.X, nb.X)
	case *dst.SelectorExpr:
		na := a.(*dst.SelectorExpr)
		nb, ok := b.(*dst.SelectorExpr)
		return ok && nodesEqual(na.X, nb.X) && nodesEqual(na.Sel, nb.Sel)
	case *dst.Ident:
		na := a.(*dst.Ident)
		nb, ok := b.(*dst.Ident)
		return ok && (na.Name == nb.Name && na.Path == nb.Path)
	case *dst.BasicLit:
		na := a.(*dst.BasicLit)
		nb, ok := b.(*dst.BasicLit)
		return ok && na.Kind == nb.Kind && na.Value == nb.Value
	case *dst.CallExpr:
		na := a.(*dst.CallExpr)
		nb, ok := b.(*dst.CallExpr)
		return ok && nodesEqual(na.Fun, nb.Fun) && exprListsEqual(na.Args, nb.Args)
	case *dst.UnaryExpr:
		na := a.(*dst.UnaryExpr)
		nb, ok := b.(*dst.UnaryExpr)
		return ok && na.Op == nb.Op && nodesEqual(na.X, nb.X)
	case *dst.BinaryExpr:
		na := a.(*dst.BinaryExpr)
		nb, ok := b.(*dst.BinaryExpr)
		return ok && nodesEqual(na.X, nb.X) && na.Op == nb.Op && nodesEqual(na.Y, nb.Y)
	case *dst.CompositeLit:
		na := a.(*dst.CompositeLit)
		nb, ok := b.(*dst.CompositeLit)
		return ok && nodesEqual(na.Type, nb.Type) && exprListsEqual(na.Elts, nb.Elts)
	case *dst.AssignStmt:
		na := a.(*dst.AssignStmt)
		nb, ok := b.(*dst.AssignStmt)
		return ok && exprListsEqual(na.Lhs, nb.Lhs) && exprListsEqual(na.Rhs, nb.Rhs) && na.Tok == nb.Tok
	case *dst.IfStmt:
		na := a.(*dst.IfStmt)
		nb, ok := b.(*dst.IfStmt)
		return ok && nodesEqual(na.Cond, nb.Cond) && nodesEqual(na.Body, nb.Body)
	case *dst.BlockStmt:
		na := a.(*dst.BlockStmt)
		nb, ok := b.(*dst.BlockStmt)
		return ok && stmtListsEqual(na.List, nb.List)
	case *dst.DeferStmt:
		na := a.(*dst.DeferStmt)
		nb, ok := b.(*dst.DeferStmt)
		return ok && nodesEqual(na.Call, nb.Call)
	case *dst.ExprStmt:
		na := a.(*dst.ExprStmt)
		nb, ok := b.(*dst.ExprStmt)
		return ok && nodesEqual(na.X, nb.X)
	case *dst.InterfaceType:
		na := a.(*dst.InterfaceType)
		nb, ok := b.(*dst.InterfaceType)
		return ok && fieldListsEqual(na.Methods, nb.Methods)
	default:
		return false
	}
}

func normalizePos(pos, total int) int {
	if pos == -1 || pos > total {
		return total
	}

	if pos < 0 {
		return 0
	}

	return pos
}
