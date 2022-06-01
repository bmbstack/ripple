package rst

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

// HasFieldInFuncDeclParams checks if the declaration params of the function, contains the given field
func HasFieldInFuncDeclParams(df *dst.File, funcName string, field *dst.Field) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				funcType := nn.Type
				for _, ff := range funcType.Params.List {
					if nodesEqual(ff, field) {
						ret = true
					}
				}
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

// DeleteFieldFromFuncDeclParams deletes any field, in the declaration params of the function,
// that is semantically equal to given field
func DeleteFieldFromFuncDeclParams(df *dst.File, funcName string, field *dst.Field) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				funcType := nn.Type
				var newList []*dst.Field
				for _, ff := range funcType.Params.List {
					if !nodesEqual(ff, field) {
						newList = append(newList, ff)
					} else {
						modified = true
					}
				}
				funcType.Params.List = newList
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

// AddFieldToFuncDeclParams adds given field, to the declaration params of the function, in the given position
func AddFieldToFuncDeclParams(df *dst.File, funcName string, field *dst.Field, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			nn := node.(*dst.FuncDecl)
			if nn.Name.Name == funcName {
				funcType := nn.Type
				fieldList := funcType.Params.List
				pos = normalizePos(pos, len(fieldList))
				funcType.Params.List = append(
					fieldList[:pos],
					append([]*dst.Field{dst.Clone(field).(*dst.Field)}, fieldList[pos:]...)...)
				modified = true
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}
