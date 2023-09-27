package rst

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

// AddEmptyLineIntoFuncBody add an empty line, inside the body of function,
func AddEmptyLineIntoFuncBody(df *dst.File, funcName string) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				nn.Body.Decs = dst.BlockStmtDecorations{Lbrace: dst.Decorations{"\n"}}
				modified = true
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

// HasStmtInsideFuncBodyWithRecv checks if the body of function has given statement
func HasStmtInsideFuncBodyWithRecv(df *dst.File, recvName, funcName string, stmt dst.Stmt) (ret bool) {
	var inside bool

	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				if nn.Recv.NumFields() > 0 {
					recvType := nn.Recv.List[0].Type
					id := &dst.Ident{Name: recvName, Path: ""}
					star := &dst.StarExpr{
						X: id,
					}
					if nodesEqual(recvType, id) || nodesEqual(recvType, star) {
						inside = true
					}
				}
			}
		case dst.Stmt:
			if inside && nodesEqual(node, stmt) {
				ret = true
			}
		}
		return true
	}

	post := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName && inside {
				if nn.Recv.NumFields() > 0 {
					recvType := nn.Recv.List[0].Type
					id := &dst.Ident{Name: recvName, Path: ""}
					star := &dst.StarExpr{
						X: id,
					}
					if nodesEqual(recvType, id) || nodesEqual(recvType, star) {
						inside = false
					}
				}
			}
		}

		return true
	}

	dstutil.Apply(df, pre, post)
	return
}

// HasStmtInsideFuncBody checks if the body of function has given statement
func HasStmtInsideFuncBody(df *dst.File, funcName string, stmt dst.Stmt) (ret bool) {
	var inside bool

	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				inside = true
			}
		case dst.Stmt:
			if inside && nodesEqual(node, stmt) {
				ret = true
			}
		}
		return true
	}

	post := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName && inside {
				inside = false
			}
		}

		return true
	}

	dstutil.Apply(df, pre, post)
	return
}

// DeleteStmtFromFuncBody deletes any statement, inside the body of function,
// that is semantically equal to the given statement.
func DeleteStmtFromFuncBody(df *dst.File, funcName string, stmt dst.Stmt) (modified bool) {
	var inside bool

	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				inside = true
			}
		case dst.Stmt:
			if inside && nodesEqual(node, stmt) {
				c.Delete()
				modified = true
			}
		}
		return true
	}

	post := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName && inside {
				inside = false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, post)
	return
}

// DeleteAllStmtFromFuncBody deletes any statement, inside the body of function,
// that is semantically equal to the given statement.
func DeleteAllStmtFromFuncBody(df *dst.File, funcName string) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				nn.Body.List = []dst.Stmt{}
				modified = true
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

// DeleteAllStmtFromFuncBodyWithRecv deletes any statement, inside the body of function,
// that is semantically equal to the given statement.
func DeleteAllStmtFromFuncBodyWithRecv(df *dst.File, recvName, funcName string) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				if nn.Recv.NumFields() > 0 {
					recvType := nn.Recv.List[0].Type
					id := &dst.Ident{Name: recvName, Path: ""}
					star := &dst.StarExpr{
						X: id,
					}
					if nodesEqual(recvType, id) || nodesEqual(recvType, star) {
						nn.Body.List = []dst.Stmt{}
						modified = true
						return false
					}
				}
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

// DeleteSelectorExprFromFuncBody deletes any SelectorExpr equal to the given one, inside the body of function.
func DeleteSelectorExprFromFuncBody(df *dst.File, funcName string, selectorExpr dst.Expr) (modified bool) {
	var inside bool
	var found bool

	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				inside = true
			}
		case *dst.SelectorExpr:
			nn := node.(*dst.SelectorExpr)
			if inside && nodesEqual(nn, selectorExpr) {
				found = true
			}
		}

		return true
	}

	post := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName && inside {
				inside = false
			}
		case *dst.ExprStmt:
			if found {
				found, modified = false, true
				c.Delete()
			}
		}
		return true
	}

	dstutil.Apply(df, pre, post)
	return
}

// AddStmtToFuncBody adds given statement, to the body of function, in the given position
func AddStmtToFuncBody(df *dst.File, funcName string, stmt dst.Stmt, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				stmtList := nn.Body.List
				pos = normalizePos(pos, len(stmtList))

				nn.Body.List = append(
					stmtList[:pos],
					append([]dst.Stmt{dst.Clone(stmt).(dst.Stmt)}, stmtList[pos:]...)...)
				modified = true
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

// AddStmtToFuncBodyStart adds given statement, to the start of function body
func AddStmtToFuncBodyStart(df *dst.File, funcName string, stmt dst.Stmt) (modified bool) {
	return AddStmtToFuncBody(df, funcName, stmt, 0)
}

// AddStmtToFuncBodyEnd adds given statement, to the end of function body
func AddStmtToFuncBodyEnd(df *dst.File, funcName string, stmt dst.Stmt) (modified bool) {
	return AddStmtToFuncBody(df, funcName, stmt, -1)
}

// AddStmtToFuncBodyWithRecv adds given statement, to the body of function, in the given position
func AddStmtToFuncBodyWithRecv(df *dst.File, recvName string, funcName string, stmt dst.Stmt, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				if nn.Recv.NumFields() > 0 {
					recvType := nn.Recv.List[0].Type
					id := &dst.Ident{Name: recvName, Path: ""}
					star := &dst.StarExpr{
						X: id,
					}
					if nodesEqual(recvType, id) || nodesEqual(recvType, star) {
						stmtList := nn.Body.List
						pos = normalizePos(pos, len(stmtList))

						nn.Body.List = append(
							stmtList[:pos],
							append([]dst.Stmt{dst.Clone(stmt).(dst.Stmt)}, stmtList[pos:]...)...)
						modified = true
						return false
					}
				}
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

// AddStmtToFuncBodyStartWithRecv adds given statement, to the start of function body
func AddStmtToFuncBodyStartWithRecv(df *dst.File, recvName, funcName string, stmt dst.Stmt) (modified bool) {
	return AddStmtToFuncBodyWithRecv(df, recvName, funcName, stmt, 0)
}

// AddStmtToFuncBodyEndWithRecv adds given statement, to the end of function body
func AddStmtToFuncBodyEndWithRecv(df *dst.File, recvName, funcName string, stmt dst.Stmt) (modified bool) {
	return AddStmtToFuncBodyWithRecv(df, recvName, funcName, stmt, -1)
}

const (
	relativeDirectionBefore = iota
	relativeDirectionAfter
)

func addStmtToFuncBodyRelativeTo(df *dst.File, funcName string, stmt, refStmt dst.Stmt, relDirection int) (modified bool) {
	var inside bool
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				inside = true
			}
		case dst.Stmt:
			ss := node.(dst.Stmt)
			if inside && nodesEqual(ss, refStmt) {
				switch relDirection {
				case relativeDirectionBefore:
					c.InsertBefore(dst.Clone(stmt))
					modified = true
				case relativeDirectionAfter:
					c.InsertAfter(dst.Clone(stmt))
					modified = true
				}
			}
		}
		return true
	}

	post := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName && inside {
				inside = false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, post)
	return
}

// AddStmtToFuncBodyBefore adds given statement, to the function body, before the position of refStmt.
// if refStmt not found, nothing will happen
func AddStmtToFuncBodyBefore(df *dst.File, funcName string, stmt, refStmt dst.Stmt) (modified bool) {
	return addStmtToFuncBodyRelativeTo(df, funcName, stmt, refStmt, relativeDirectionBefore)
}

// AddStmtToFuncBodyAfter adds given statement, to the function body, after the position of refStmt,
// if refStmt not found, nothing will happen
func AddStmtToFuncBodyAfter(df *dst.File, funcName string, stmt, refStmt dst.Stmt) (modified bool) {
	return addStmtToFuncBodyRelativeTo(df, funcName, stmt, refStmt, relativeDirectionAfter)
}
