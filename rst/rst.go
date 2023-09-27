package rst

import (
	"bytes"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/guess"
	"go/token"
	"io"
	"io/ioutil"
	"os"
)

type PkgAlias struct {
	Pkg   string
	Alias string
}

var defaultFileSet = token.NewFileSet()

// ParseSrcFileFromBytes parses the given go src file, in the form of bytes, into *dst.File
func ParseSrcFileFromBytes(src []byte, resolver guess.RestorerResolver) (df *dst.File, err error) {
	dec := decorator.NewDecoratorWithImports(
		defaultFileSet,
		"main",
		goast.WithResolver(resolver))
	return dec.Parse(src)
}

// ParseSrcFile parses the given go src filename, in the form of valid path, into *dst.File
func ParseSrcFile(filename string, resolver guess.RestorerResolver) (df *dst.File, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	return ParseSrcFileFromBytes(src, resolver)
}

// FprintFile writes the *dst.File out to io.Writer
func FprintFile(out io.Writer, df *dst.File, resolver guess.RestorerResolver, alias []PkgAlias) error {
	dec := decorator.NewRestorerWithImports("main", resolver)
	fr := dec.FileRestorer()
	for _, item := range alias {
		fr.Alias[item.Pkg] = item.Alias
	}
	return fr.Fprint(out, df)
}

func PrintToBuf(df *dst.File, resolver guess.RestorerResolver, alias []PkgAlias) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	_ = FprintFile(buf, df, resolver, alias)
	return buf
}
