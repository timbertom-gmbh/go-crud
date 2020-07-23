package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/types"
	"io"
	"os"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/packages"
)

var (
	modelName  = flag.String("model", "", "Model entity name")
	modelpkg   = flag.String("modelpkg", "", "Model package")
	rpcName    = flag.String("rpc", "", "RPC message object name")
	rpcpkg     = flag.String("rpcpkg", "", "RPC package")
	operations = flag.String("ops", "crudl", "List of characters defining what operations should be availbale")
	out        = flag.String("out", "-", "Package to write to or - for stdout")
)

func main() {
	flag.Parse()
	if *modelName == "" || *modelpkg == "" || *rpcName == "" || *rpcpkg == "" {
		panic("required argument missing")
	}

	patterns := []string{*modelpkg}
	if *modelpkg != *rpcpkg {
		patterns = append(patterns, *rpcpkg)
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode:       packages.LoadSyntax,
		Tests:      false,
		BuildFlags: []string{},
	}, patterns...)
	if err != nil {
		panic(err)
	}

	var modelType *types.Struct
	var rpcType *types.Struct
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			for _, err := range pkg.Errors {
				fmt.Println(err.Kind, err.Msg, err.Pos)
			}
			panic(1)
		}

		tp := pkg.Types
		scope := tp.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			t := obj.Type()

			namedType, namedOk := t.(*types.Named)
			if !namedOk {
				continue
			}

			if pkg.ID == *modelpkg && namedType.Obj().Id() == *modelName {
				structType, structOk := namedType.Underlying().(*types.Struct)
				if !structOk {
					panic("requested thingy found but not a struct")
				}
				modelType = structType
			}
			if pkg.ID == *rpcpkg && namedType.Obj().Id() == *rpcName {
				structType, structOk := namedType.Underlying().(*types.Struct)
				if !structOk {
					panic("requested thingy found but not a struct")
				}
				rpcType = structType
			}
		}
	}

	if modelType == nil || rpcType == nil {
		panic("not found one of the required types")
	}

	buf := bytes.NewBuffer([]byte{})
	write(buf, fmt.Sprintf("import rpc_pkg \"%s\"", *rpcpkg))
	write(buf, fmt.Sprintf("import model_pkg \"%s\"", *modelpkg))
	write(buf, "")

	writeConverters(buf, *modelName, modelType, *rpcName, rpcType)
	formatedBuf, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	os.Stdout.Write(formatedBuf)
}

func writeConverters(out io.Writer, modelName string, modelType *types.Struct, rpcName string, rpcType *types.Struct) {
	write(out, fmt.Sprintf("func convert%sToModel(rpcObj *rpc_pkg.%s) *model_pkg.%s {", rpcName, rpcName, modelName))
	write(out, fmt.Sprintf("\treturn &model_pkg.%s{", modelName))

	for i := 0; i < modelType.NumFields(); i++ {
		f := modelType.Field(i)
		tag := modelType.Tag(i)

		rpcFieldName, skipField := getRPCFieldName(tag, f.Name())
		if skipField {
			continue
		}

		write(out, fmt.Sprintf("\t\t%s: rpcObj.%s,", f.Name(), rpcFieldName))
	}
	write(out, "\t}\n}")
}

func getRPCFieldName(tag string, fallback string) (string, bool) {
	tags, err := structtag.Parse(tag)
	if err != nil {
		return fallback, false
	}

	rpcTag, err := tags.Get("rpc")
	if err != nil {
		return fallback, false
	}

	if rpcTag.Name == "-" {
		return "", true
	}

	return rpcTag.Name, false
}

func write(out io.Writer, str string) {
	out.Write([]byte(str + "\n"))
}
