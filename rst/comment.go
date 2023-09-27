package rst

import (
	"github.com/bmbstack/ripple/cmd/ripple/util"
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
	"go/token"
	"strings"
)

type StructDecs struct {
	Name string
	Decs dst.Decorations
}

func GetStructDecsInStructComment(df *dst.File, sign string) (ret []StructDecs) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.GenDecl:
			nn := node.(*dst.GenDecl)
			if nn.Tok == token.TYPE {
				if len(nn.Specs) > 0 {
					ts := nn.Specs[0].(*dst.TypeSpec)
					_, ok := ts.Type.(*dst.StructType)
					if ok {
						// is struct
						if hasTagInNodeDecs(nn.Decs.Start, sign) {
							ret = append(ret, StructDecs{
								Name: ts.Name.Name,
								Decs: nn.Decs.Start,
							})
						}
					}
				}
			}
		default:
		}
		return true
	}
	dstutil.Apply(df, pre, nil)
	return ret
}

func HasSignInFuncComment(df *dst.File, funcName, sign string) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			nn := node.(*dst.FuncDecl)
			if nn.Name.Name == funcName {
				ret = hasTagInNodeDecs(nn.Decs.Start, sign)
				return false
			}
		default:
		}
		return true
	}
	dstutil.Apply(df, pre, nil)
	return ret
}

func GetAnnotLineInFuncDeclComment(df *dst.File, funcName, sign, annot string) (ret string) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			nn := node.(*dst.FuncDecl)
			if nn.Name.Name == funcName {
				if hasTagInNodeDecs(nn.Decs.Start, sign) {
					ret = getTagLineInNodeDecs(nn.Decs.Start, annot)
					if util.IsNotEmpty(ret) {
						return false
					}
				}

			}
		default:
		}
		return true
	}
	dstutil.Apply(df, pre, nil)
	return ret
}

func TrimAnnot(s string) string {
	s = strings.ReplaceAll(s, "//", "")
	s = strings.ReplaceAll(s, "/*", "")
	s = strings.ReplaceAll(s, "*/", "")
	return s
}

func getTagLineInNodeDecs(decorations []string, tag string) (ret string) {
	for _, item := range decorations {
		if strings.Contains(item, tag) {
			ret = item
			break
		}
	}
	return ret
}

func hasTagInNodeDecs(decorations []string, tag string) (ret bool) {
	for _, item := range decorations {
		arr := strings.Split(item, " ")
		for _, value := range arr {
			if value == tag {
				ret = true
				break
			}
		}
	}
	return ret
}
