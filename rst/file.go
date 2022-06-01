package rst

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
	"go/token"
)

func HasStructDeclInFile(df *dst.File, structName string) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.GenDecl:
			nn := node.(*dst.GenDecl)
			if nn.Tok == token.TYPE {
				if isInStructTypeArray(structName, nn.Specs) {
					ret = true
					return false
				}

			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

func HasFuncDeclInFile(df *dst.File, decl dst.FuncDecl) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == decl.Name.Name {
				ret = true
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

func HasFuncDeclWithRecvInFile(df *dst.File, decl dst.FuncDecl, recvName string) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == decl.Name.Name {
				if nn.Recv.NumFields() > 0 {
					recvType := nn.Recv.List[0].Type
					id := &dst.Ident{Name: recvName, Path: ""}
					star := &dst.StarExpr{
						X: id,
					}
					if nodesEqual(recvType, id) || nodesEqual(recvType, star) {
						ret = true
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

func isInStructTypeArray(structName string, list []dst.Spec) (ret bool) {
	for _, item := range list {
		ts, ok := item.(*dst.TypeSpec)
		if ok {
			if structName == ts.Name.Name {
				ret = true
				break
			}

		}
	}
	return ret
}
