// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build codec.gen

package codec

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"runtime"
	"slices"
	"strings"
)

// ----
// MARKER 2025

// - [fn:monoBase] for each file of fast-path.generated.go, encode.go, decode.go:
//   - ignore init() and imports (go imports will fix it later - or we just grab all imports and let go-imports fix later)
//   - parse encode.go, decode.go, fast-path.generated.go into an AST
//     - for each function, if a generic function/method, add it into the AST
// - [fn:monoAll] for each handle
//   - [fn:monoHandle] clone prior AST
//   - add <format>.go into it
//   - transform the node names to monomorphize them
//   - transform the method/func calls in each method to monomorphize them
// - output AST into a file <format>.mono.generated.go

const genMonoParserMode = parser.AllErrors | parser.SkipObjectResolution

// This code works very well for things that are monomorphizing things scoped to a specific format.
// As you see, we monomorphize code for each format.
//
// However, for types bound to encWriter or decReader, which are shared across formats,
// there's no place to put them without duplication.
//
// Consequently, since there is only one function ie decByteSlice, we handle it manually.
// This method is in the type helperDecReader.
var genMonoSpecialFieldTypes = []string{"helperDecReader"}

// These functions should take the address of first param when monomorphized
var genMonoSpecialFunc4Addr = []string{"decByteSlice"}

var genMonoRefImportsVia_ = [][2]string{
	{"errors", "New"},
}

type genMono struct {
	files       map[string][]byte
	typParam    map[string]*ast.Field
	imports     map[string]struct{}
	importSpecs []ast.Spec
}

func (x *genMono) init() {
	x.files = make(map[string][]byte)
	x.typParam = make(map[string]*ast.Field)
	x.imports = make(map[string]struct{})
}

func (x *genMono) reset() {
	clear(x.typParam)
	clear(x.imports)
	clear(x.importSpecs)
}

func (m *genMono) do(h Handle) {
	m.reset()

	b, fset := m.base(true)
	r1 := m.handle(h, b, fset)
	m.trFile(r1, h, true)
	// clear(m.typParam)

	b, fset = m.base(false)
	r2 := m.handle(h, b, fset)
	m.trFile(r2, h, false)
	// clear(m.typParam)

	r0 := &ast.File{Name: &ast.Ident{Name: "codec"}}
	r0.Decls = append(r0.Decls, &ast.GenDecl{Tok: token.IMPORT, Specs: m.importSpecs})
	_vd := &ast.GenDecl{Tok: token.VAR}
	for _, v := range genMonoRefImportsVia_ {
		_vs := new(ast.ValueSpec)
		_vs.Names = append(_vs.Names, &ast.Ident{Name: "_"})
		_vs.Values = append(_vs.Values, &ast.SelectorExpr{
			X:   &ast.Ident{Name: v[0]},
			Sel: &ast.Ident{Name: v[1]},
		})
		_vd.Specs = append(_vd.Specs, _vs)
	}
	r0.Decls = append(r0.Decls, _vd)
	r0.Decls = append(r0.Decls, r1.Decls...)
	r0.Decls = append(r0.Decls, r2.Decls...)

	// output r1 to a file
	f, err := os.Create(h.Name() + ".mono.generated.go")
	halt.onerror(err)
	defer f.Close()
	err = format.Node(f, fset, r0)
	halt.onerror(err)

	m.reset()
}

func (x *genMono) file(fname string) (b []byte) {
	b = x.files[fname]
	if b == nil {
		var err error
		b, err = os.ReadFile(fname)
		halt.onerror(err)
		x.files[fname] = b
	}
	return
}

func (x *genMono) base(isbytes bool) (r *ast.File, fset *token.FileSet) {
	fname := "base.io.mono.go"
	if isbytes {
		fname = "base.bytes.mono.go"
	}

	fset = token.NewFileSet()
	r = &ast.File{
		Name: &ast.Ident{Name: "codec"},
	}
	fnames := []string{"encode.go", "decode.go", "fast-path.generated.go"}
	// fnames = []string{"encode.go", "decode.go", "fast-path.not.go"} // MARKER 2025 - remove this
	for _, fname = range fnames {
		fsrc := x.file(fname)
		f, err := parser.ParseFile(fset, fname, fsrc, genMonoParserMode)
		halt.onerror(err)
		x.merge(r, f)
	}
	return
}

func (x *genMono) handle(h Handle, base *ast.File, fset *token.FileSet) (r *ast.File) {
	r = genMonoCopy(base)
	hname := h.Name()
	fname := hname + ".go"
	fsrc := x.file(fname)
	f, err := parser.ParseFile(fset, fname, fsrc, genMonoParserMode)
	halt.onerror(err)
	x.merge(r, f)
	return
}

func (x *genMono) merge(dst, src *ast.File) {
	// we only merge top-level methods and types
	fnIdX := func(node ast.Expr) (ok bool) {
		n2, ok := node.(*ast.IndexExpr)
		if ok {
			n9, ok9 := n2.Index.(*ast.Ident)
			n3, ok := n2.X.(*ast.Ident)
			ok = ok && ok9 && n9.Name == "T"
			if ok {
				_, ok = x.typParam[n3.Name]
			}
		}
		return
	}

	fn := func(node ast.Node) bool {
		var ok bool
		switch n := node.(type) {
		case *ast.FuncDecl:
			// TypeParams is nil for methods, as it is defined at the type node
			// instead, look at the name, and
			// if IndexExpr.Index=T, and IndexExpr.X matches a type name seen already
			//     then ok = true
			if n.Recv == nil || len(n.Recv.List) != 1 {
				return false
			}
			ok = false
			switch nn := n.Recv.List[0].Type.(type) {
			case *ast.IndexExpr:
				ok = fnIdX(nn)
			case *ast.StarExpr:
				ok = fnIdX(nn.X)
			}
			if ok {
				dst.Decls = append(dst.Decls, n)
			}
			return false
		case *ast.GenDecl:
			if n.Tok == token.TYPE {
				for _, v := range n.Specs {
					nn := v.(*ast.TypeSpec)
					ok = genMonoTypeParamsOk(nn.TypeParams)
					if ok {
						// each decl will have only 1 var/type
						x.typParam[nn.Name.Name] = nn.TypeParams.List[0]
						dst.Decls = append(dst.Decls, &ast.GenDecl{Tok: n.Tok, Specs: []ast.Spec{v}})
					}
				}
			} else if n.Tok == token.IMPORT {
				for _, v := range n.Specs {
					nn := v.(*ast.ImportSpec)
					if _, ok := x.imports[nn.Path.Value]; !ok {
						x.importSpecs = append(x.importSpecs, v)
						x.imports[nn.Path.Value] = struct{}{}
					}
				}
			}
			return false
		}
		return true
	}
	ast.Inspect(src, fn)
}

func (x *genMono) trFile(r *ast.File, h Handle, isbytes bool) {
	fn := func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.TypeSpec:
			// type x[T encDriver] struct { ... }
			if !genMonoTypeParamsOk(n.TypeParams) {
				return false
			}
			x.trType(n, h, isbytes)
			return false
		case *ast.FuncDecl:
			if n.Recv == nil || len(n.Recv.List) != 1 {
				return false
			}
			if _, ok := n.Recv.List[0].Type.(*ast.Ident); ok {
				return false
			}
			tp := x.trMethodSign(n, h, isbytes) // receiver, params, results
			// handle the body
			x.trMethodBody(n.Body, tp, h, isbytes)
			return false
		}
		return true
	}
	ast.Inspect(r, fn)

	// set type params to nil, and Pos to NoPos
	fn = func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.FuncType:
			if genMonoTypeParamsOk(n.TypeParams) {
				n.TypeParams = nil
			}
		case *ast.TypeSpec: // for type ...
			if genMonoTypeParamsOk(n.TypeParams) {
				n.TypeParams = nil
			}
		}
		return true
	}
	ast.Inspect(r, fn)
}

func (x *genMono) trType(n *ast.TypeSpec, h Handle, isbytes bool) {
	sfx, _, _, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	tp := n.TypeParams.List[0]
	switch tp.Type.(*ast.Ident).Name {
	case "encDriver", "decDriver":
		n.Name.Name += hnameUp + sfx
	case "encWriter", "decReader":
		n.Name.Name += sfx
	}

	// handle the Struct and Array types
	switch nn := n.Type.(type) {
	case *ast.StructType:
		x.trStruct(nn, tp, h, isbytes)
	case *ast.ArrayType:
		x.trArray(nn, tp, h, isbytes)
	}
}

func (x *genMono) trMethodSign(n *ast.FuncDecl, h Handle, isbytes bool) (tp *ast.Field) {
	// check if recv type is not parameterized
	tp = x.trField(n.Recv.List[0], nil, h, isbytes, true)
	// handle params and results
	x.trMethodSignNonRecv(n.Type.Params, tp, h, isbytes)
	x.trMethodSignNonRecv(n.Type.Results, tp, h, isbytes)
	return
}

func (x *genMono) trMethodSignNonRecv(r *ast.FieldList, tp *ast.Field, h Handle, isbytes bool) {
	if r == nil || len(r.List) == 0 {
		return
	}
	for _, v := range r.List {
		x.trField(v, tp, h, isbytes, false)
	}
}

func (x *genMono) trStruct(r *ast.StructType, tp *ast.Field, h Handle, isbytes bool) {
	// search for fields, and update accordingly
	// type x[T encDriver] struct { w T }
	// var x *A[T]
	// A[T]
	if r == nil || r.Fields == nil || len(r.Fields.List) == 0 {
		return
	}
	for _, v := range r.Fields.List {
		x.trField(v, tp, h, isbytes, false)
	}
}

func (x *genMono) trArray(n *ast.ArrayType, tp *ast.Field, h Handle, isbytes bool) {
	sfx, _, _, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	// type fastpathEs[T encDriver] [56]fastpathE[T]
	// p := tp.Names[0].Name
	switch elt := n.Elt.(type) {
	// case *ast.InterfaceType: // MARKER 2025
	case *ast.IndexExpr:
		if elt.Index.(*ast.Ident).Name == "T" { // generic
			n.Elt = ast.NewIdent(elt.X.(*ast.Ident).Name + hnameUp + sfx)
		}
	}
}

// func (x *genMono) trVar(r *ast.ValueSpec, tp *ast.Field, h Handle, isbytes bool) {
// }

func (x *genMono) trMethodBody(r *ast.BlockStmt, tp *ast.Field, h Handle, isbytes bool) {
	// sfx, writer, reader, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	// find the parent node for an indexExpr, or a T/*T, and set the value back in there

	fn := func(pnode ast.Node) bool {
		var pn *ast.Ident
		fnUp := func() {
			x.updateIdentForT(pn, h, tp, isbytes, false)
		}
		switch n := pnode.(type) {
		// case *ast.SelectorExpr:
		// case *ast.TypeAssertExpr:
		// case *ast.IndexExpr:
		case *ast.StarExpr:
			if genMonoUpdateIndexExprT(&pn, n.X) {
				n.X = pn
				fnUp()
			}
		case *ast.CallExpr:
			for i4, n4 := range n.Args {
				if genMonoUpdateIndexExprT(&pn, n4) {
					n.Args[i4] = pn
					fnUp()
				}
			}
			// special case: if decByteSlice function, take addr of first parameter // MARKER 2025
			if n4, ok4 := n.Fun.(*ast.Ident); ok4 && slices.Contains(genMonoSpecialFunc4Addr, n4.Name) {
				n.Args[0] = &ast.UnaryExpr{Op: token.AND, X: n.Args[0].(*ast.SelectorExpr)}
			}
		case *ast.CompositeLit:
			if genMonoUpdateIndexExprT(&pn, n.Type) {
				n.Type = pn
				fnUp()
			}
		case *ast.ArrayType:
			if genMonoUpdateIndexExprT(&pn, n.Elt) {
				n.Elt = pn
				fnUp()
			}
		case *ast.ValueSpec:
			for i2, n2 := range n.Values {
				if genMonoUpdateIndexExprT(&pn, n2) {
					n.Values[i2] = pn
					fnUp()
				}
			}
			if genMonoUpdateIndexExprT(&pn, n.Type) {
				n.Type = pn
				fnUp()
			}
		case *ast.BinaryExpr:
			// early return here, since the 2 things can apply
			if genMonoUpdateIndexExprT(&pn, n.X) {
				n.X = pn
				fnUp()
			}
			if genMonoUpdateIndexExprT(&pn, n.Y) {
				n.Y = pn
				fnUp()
			}
			return true
		}
		return true
	}
	ast.Inspect(r, fn)
}

func (x *genMono) trField(f *ast.Field, tpt *ast.Field, h Handle, isbytes, isRecv bool) (tp *ast.Field) {
	// sfx, writer, reader, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	var pn *ast.Ident
	switch nn := f.Type.(type) {
	case *ast.IndexExpr:
		if genMonoUpdateIndexExprT(&pn, nn) {
			f.Type = pn
		}
	case *ast.StarExpr:
		if genMonoUpdateIndexExprT(&pn, nn.X) {
			nn.X = pn
		}
	case *ast.FuncType:
		x.trMethodSignNonRecv(nn.Params, tpt, h, isbytes)
		x.trMethodSignNonRecv(nn.Results, tpt, h, isbytes)
		return
	case *ast.ArrayType:
		x.trArray(nn, tpt, h, isbytes)
		return
	case *ast.Ident:
		if isRecv || nn.Name != "T" {
			return
		}
		pn = nn
	}
	if pn == nil {
		return
	}

	tp = x.updateIdentForT(pn, h, tpt, isbytes, true)
	// if isRecv {
	// 	fmt.Printf(">>>> pn: %s\n", pn)
	// }
	return
}

func (x *genMono) updateIdentForT(pn *ast.Ident, h Handle, tp *ast.Field, isbytes, lookupTP bool) (tp2 *ast.Field) {
	sfx, writer, reader, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	// handle special ones e.g. helperDecReader et al
	if slices.Contains(genMonoSpecialFieldTypes, pn.Name) {
		pn.Name += sfx
		return
	}

	// fmt.Printf(">>>> tp: %v\n", tp)
	if pn.Name != "T" && lookupTP {
		tp = x.typParam[pn.Name]
		// fmt.Printf(">>>> tp updated: %v\n", tp)
	}

	paramtyp := tp.Type.(*ast.Ident).Name
	if pn.Name == "T" {
		switch paramtyp {
		case "encDriver", "decDriver":
			pn.Name = h.Name() + genMonoTitleCase(paramtyp) + sfx
		case "encWriter":
			pn.Name = writer
		case "decReader":
			pn.Name = reader
		}
	} else {
		switch paramtyp {
		case "encDriver", "decDriver":
			pn.Name += hnameUp + sfx
		case "encWriter", "decReader":
			pn.Name += sfx
		}
	}
	return tp
}

func genMonoUpdateIndexExprT(pn **ast.Ident, node ast.Node) (pnok bool) {
	*pn = nil
	if n2, ok := node.(*ast.IndexExpr); ok {
		n9, ok9 := n2.Index.(*ast.Ident)
		n3, ok := n2.X.(*ast.Ident)
		if ok && ok9 && n9.Name == "T" {
			*pn, pnok = ast.NewIdent(n3.Name), true
		}
	}
	return
}

func genMonoTitleCase(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func genMonoIsBytesVals(hName string, isbytes bool) (suffix, writer, reader, hNameUp string) {
	hNameUp = genMonoTitleCase(hName)
	if isbytes {
		return "Bytes", "bytesEncAppender", "bytesDecReader", hNameUp
	}
	return "IO", "bufioEncWriter", "ioDecReader", hNameUp
}

func genMonoTypeParamsOk(v *ast.FieldList) bool {
	if v == nil || v.List == nil || len(v.List) != 1 {
		return false
	}
	if len(v.List[0].Names) != 1 || v.List[0].Names[0].Name != "T" {
		return false
	}
	switch v.List[0].Type.(*ast.Ident).Name {
	case "encDriver", "decDriver", "encWriter", "decReader":
		return true
	}
	return false
}

func genMonoCopy(src *ast.File) (dst *ast.File) {
	dst = &ast.File{
		Name: &ast.Ident{Name: "codec"},
	}
	dst.Decls = append(dst.Decls, src.Decls...)
	return
}

func GenMonoAll() {
	hdls := []Handle{
		(*SimpleHandle)(nil),
		(*JsonHandle)(nil),
		(*CborHandle)(nil),
		(*BincHandle)(nil),
		(*MsgpackHandle)(nil),
	}
	// hdls = []Handle{(*MsgpackHandle)(nil)} // MARKER 2025 remove
	var m genMono
	m.init()
	for _, v := range hdls {
		m.do(v)
		runtime.GC()
	}
}

/*
const genMonoBaseStrBytes = `
package codec

type helperDecReaderBytes struct{}

func (helperDecReaderBytes) decByteSlice(r *bytesDecReader, clen, maxInitLen int, bs []byte) (bsOut []byte) {
	if clen <= 0 {
		bsOut = zeroByteSlice
	} else if cap(bs) >= clen {
		bsOut = bs[:clen]
		r.readb(bsOut)
	} else {
		var len2 int
		for len2 < clen {
			len3 := decInferLen(clen-len2, maxInitLen, 1)
			bs3 := bsOut
			bsOut = make([]byte, len2+len3)
			copy(bsOut, bs3)
			r.readb(bsOut[len2:])
			len2 += len3
		}
	}
	return
}
`

const genMonoBaseStrIO = `
package codec

type helperDecReaderIO struct{}

func (helperDecReaderIO) decByteSlice(r *ioDecReader, clen, maxInitLen int, bs []byte) (bsOut []byte) {
	if clen <= 0 {
		bsOut = zeroByteSlice
	} else if cap(bs) >= clen {
		bsOut = bs[:clen]
		r.readb(bsOut)
	} else {
		var len2 int
		for len2 < clen {
			len3 := decInferLen(clen-len2, maxInitLen, 1)
			bs3 := bsOut
			bsOut = make([]byte, len2+len3)
			copy(bsOut, bs3)
			r.readb(bsOut[len2:])
			len2 += len3
		}
	}
	return
}
`

func (x *genMono) transform(r *ast.File, h Handle, isbytes bool) {
	var sfx string
	if isbytes {
		sfx = "Bytes"
	} else {
		sfx = "IO"
	}
	up := func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	hname := h.Name()

	// Type parameters only exist in
	//   - function/method declarations (function, method)
	//     - signature
	//     - body
	//   - type declarations (struct, array)
	//     - signature
	//     - body
	//   - type instantiations (N/A as moved into init.go)
	//   - type constraints (N/A as moved into init.go)
	//
	// Consequently, we only need to work on func and type declarations
	// and things that flow down from there.
	//
	// We only want to work on Field values which could have TypeParams.

	typeParams := make(map[*ast.FieldList]struct{})
	var typParam *ast.Field

	var stack []ast.Node
	pop := func() {
		stack = stack[:len(stack)-1]
	}
	push := func(n ast.Node) {
		stack = append(stack, n)
	}
	// last := func(n depth) ast.Node {
	// 	pos := len(stack) - 1 - n
	// 	if pos >= 0 {
	// 		return stack[pos]
	// 	}
	// 	return nil
	// }
	locate := func(f func(n ast.Node) bool) ast.Node {
		for i := len(stack) - 1; i >= 0; i-- {
			if n := stack[i]; f(n) {
				return n
			}
		}
		return nil
	}

	locateTypeSpec := func() *ast.TypeSpec {
		fp := func(node ast.Node) bool {
			_, ok := node.(*ast.TypeSpec)
			return ok
		}
		v := locate(fp)
		if v != nil {
			return v.(*ast.TypeSpec)
		}
		return nil
	}

	// if type parameter is encDriver, use <name><handle><IO|Bytes>
	// e.g. if type is encoder[T], do encoderJsonBytes --> jsonEncDriverBytes
	fnNameViaType := func(nn *ast.Ident) {
		paramName := typParam.Names[0].Name
		paramType := typParam.Type.(*ast.Ident).Name
		if nn.Name == "" || nn.Name == paramName {
			switch paramType {
			case "encDriver", "decDriver":
				nn.Name = hname + up(paramType) + sfx
			case "encWriter":
				if isbytes {
					nn.Name = "bytesEncAppender"
				} else {
					nn.Name = "bufioEncWriter"
				}
			case "decReader":
				if isbytes {
					nn.Name = "bytesDecReader"
				} else {
					nn.Name = "ioDecReader"
				}
			}
			return
		}

		switch paramType {
		case "encDriver", "decDriver":
			// encoder or decoder has a field of encDriver / decDriver
			nn.Name += up(hname) + sfx
		case "encWriter", "decReader":
			// already en encDriver or decDriver
			nn.Name += sfx
		}
	}

	fn := func(node ast.Node) bool {
		if node == nil {
			pop()
			return false
		}
		push(node)

		typParam = nil
		switch n := node.(type) {
		case *ast.FuncDecl:
			// if !genMonoTypeParamsOk(n.Type.TypeParams) {
			// 	return false
			// }
			// typeParams[n.Type.TypeParams] = struct{}{}
			// always methods. change method's receiver type only

			var rtyp string
			switch nn := n.Recv.List[0].Type.(type) {
			case *ast.IndexExpr:
				if nn2, ok := nn.X.(*ast.Ident); ok {
					rtyp = nn2.Name
				}
			case *ast.StarExpr:
				if nn2, ok := nn.X.(*ast.IndexExpr); ok {
					if nn3, ok := nn2.X.(*ast.Ident); ok {
						rtyp = nn3.Name
					}
				}
			}
			if rtyp != "" {
				typParam = x.typParam[rtyp]
			}
			if typParam == nil {
				return false
			}
			// typParam = n.Type.TypeParams.List[0]
			pnew := ast.NewIdent(rtyp)
			n.Recv.List[0].Type = pnew
			fnNameViaType(pnew)
		case *ast.FuncType:
			// update the params and results, using type parameters
			typeParams[n.TypeParams] = struct{}{}
			// if n.TypeParams == nil {
			// 	return false
			// }
			// field := n.TypeParams.List[0].Type.(*ast.Ident)
			// for _, v0 := range []*ast.FieldList{n.Params, n.Results} {
			// 	if v0 == nil {
			// 		continue
			// 	}
			// 	// fmt.Printf("v0: %v, len(v0.List): %v\n", v0, v0)
			// 	for _, v1 := range v0.List {
			// 		// fmt.Printf("v1: %v, len(v1.Names): %v\n", v1, len(v1.Names))
			// 		for _, v2 := range v1.Names {
			// 			fnNameViaType(v2, field.Name)
			// 		}
			// 	}
			// }
		case *ast.FieldList:
			// skip all type params
			if _, ok := typeParams[n]; ok {
				return false
			}
		case *ast.Field:
			// type x[T encDriver] struct { w T }
			if n.Type == nil {
				return false
			}
			p := locateTypeSpec()
			if p == nil {
				return false
			}

			typParam = p.TypeParams.List[0]
			paramName := typParam.Names[0].Name
			// if any generic field, handle it using p's TypeParams
			switch nt := n.Type.(type) {
			case *ast.Ident:
				if nt.Name == paramName {
					fnNameViaType(nt)
				}
			case *ast.IndexExpr:
				if nt.Index.(*ast.Ident).Name == paramName {
					pnew := ast.NewIdent(nt.X.(*ast.Ident).Name)
					n.Type = pnew
					fnNameViaType(pnew)
				}
			case *ast.StarExpr:
				switch ni := nt.X.(type) {
				case *ast.IndexExpr:
					if ni.Index.(*ast.Ident).Name == paramName {
						pnew := ast.NewIdent(ni.X.(*ast.Ident).Name)
						nt.X = pnew
						fnNameViaType(pnew)
					}
				}
			}
		case *ast.TypeSpec:
			// type x[T encDriver] struct { ... }
			if !genMonoTypeParamsOk(n.TypeParams) {
				return false
			}
			typeParams[n.TypeParams] = struct{}{}
			typParam = n.TypeParams.List[0]
			fnNameViaType(n.Name)
		case *ast.ValueSpec:
			// var X encoder[T]
			p := locateTypeSpec()
			if p == nil {
				return false
			}
			if !genMonoTypeParamsOk(p.TypeParams) {
				return false
			}
			typParam = p.TypeParams.List[0]
			// MARKER 2025 - need to set the type of the variable (not the name)
			// param := p.TypeParams.List[0].Type.(*ast.Ident)
			// fnNameViaType(n.Name, n.TypeParams.List[0].Type.(*ast.Ident).Name)
		case *ast.ArrayType:
			// type fastpathEs[T encDriver] [0]fastpathE[T]
			p := locateTypeSpec()
			if p == nil {
				return false
			}
			typParam = p.TypeParams.List[0]
			paramName := typParam.Names[0].Name
			switch elt := n.Elt.(type) {
			case *ast.IndexExpr:
				if elt.Index.(*ast.Ident).Name == paramName { // generic
					pnew := ast.NewIdent(elt.X.(*ast.Ident).Name)
					n.Elt = pnew
					fnNameViaType(pnew)
				}
			case *ast.InterfaceType: // MARKER 2025
			default:
			}
		}
		return true
	}
	ast.Inspect(r, fn)

	// set type params to nil, and Pos to NoPos
	fn = func(node ast.Node) bool {
		switch n := node.(type) {
		// case *ast.FuncDecl:
		// 	n.Type.TypeParams = nil
		case *ast.FuncType:
			if genMonoTypeParamsOk(n.TypeParams) {
				n.TypeParams = nil
			}
		case *ast.TypeSpec: // for type ...
			if genMonoTypeParamsOk(n.TypeParams) {
				n.TypeParams = nil
			}
		}
		return true
	}
	ast.Inspect(r, fn)

}

func (x *genMono) trMethodBody(r *ast.BlockStmt, tp *ast.Field, h Handle, isbytes bool) {
	// sfx, writer, reader, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	// find the parent node for an indexExpr, and set the value back in there
	// var ns genMonoNodeStack
	fn = func(pnode ast.Node) bool {
		// if node == nil {
		// 	ns.pop()
		// 	return false
		// }
		var pn *ast.Ident
		switch n := pnode.(type) {
		// case *ast.SelectorExpr:
		case *ast.IndexExpr:
			if nn2, ok := nn.X.(*ast.Ident); ok {
				pn = ast.NewIdent(nn2.Name)
				nlast := ns.peek()
				switch nn3 := nlast.(type) {
				case *ast.StarExpr:
					f.Type = pn
				}
			}
		}
		if pn != nil {
			x.updateIdentForT(pn, h, tp, false)
		}
			// ns.push(node) // marker 2025 - should this be here?
		return true
	}
	ast.Inspect(r, fn)
}

func (x *genMono) trField(f *ast.Field, tpt *ast.Field, h Handle, isbytes, isRecv bool) (tp *ast.Field) {
	// sfx, writer, reader, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	var pn *ast.Ident
	fnIdX := func(node2 ast.Node) bool {
		var pnok bool
		pn = nil
		if n2, ok := node2.(*ast.IndexExpr); ok {
			n9, ok9 := n2.Index.(*ast.Index)
			n3, ok := n2.X.(*ast.Ident)
			if ok && ok9 {
				pn, pnok = ast.NewIdent(n3.Name), true
			}
		}
		return pnok
	}
	switch nn := f.Type.(type) {
	case *ast.IndexExpr:
		if fnIdX(nn) {
			f.Type = pn
		}
		if nn2, ok := nn.X.(*ast.Ident); ok {
			pn = ast.NewIdent(nn2.Name)
			f.Type = pn
		}
	case *ast.StarExpr:
		if nn2, ok := nn.X.(*ast.IndexExpr); ok {
			if nn3, ok := nn2.X.(*ast.Ident); ok {
				pn = ast.NewIdent(nn3.Name)
				nn.X = pn
			}
		}
	case *ast.FuncType:
		x.trMethodSignNonRecv(nn.Params, tpt, h, isbytes)
		x.trMethodSignNonRecv(nn.Results, tpt, h, isbytes)
		return
	case *ast.ArrayType:
		x.trArray(nn, tpt, h, isbytes)
		return
	case *ast.Ident:
		if isRecv || nn.Name != "T" {
			return
		}
		pn = nn
	}
	if pn == nil {
		return
	}

	x.updateIdentForT(pn, h, tpt, isbytes, true)
	// // handle special ones e.g. helperDecReader et al
	// if slices.Contains(genMonoSpecialFieldTypes, pn.Name) {
	// 	pn.Name += sfx
	// 	return
	// }

	// if pn.Name == "T" {
	// 	tp = tpt
	// } else {
	// 	tp = x.typParam[pn.Name]
	// }

	// paramtyp := tp.Type.(*ast.Ident).Name
	// if pn.Name == "T" {
	// 	switch paramtyp {
	// 	case "encDriver", "decDriver":
	// 		pn.Name = h.Name() + genMonoTitleCase(paramtyp) + sfx
	// 	case "encWriter":
	// 		pn.Name = writer
	// 	case "decReader":
	// 		pn.Name = reader
	// 	}
	// } else {
	// 	switch paramtyp {
	// 	case "encDriver", "decDriver":
	// 		pn.Name += hnameUp + sfx
	// 	case "encWriter", "decReader":
	// 		pn.Name += sfx
	// 	}
	// }
	return
}

const genMonoBaseStrBytes = `
package codec
`

const genMonoBaseStrIO = `
package codec
`

func (x *genMono) base(isbytes bool) (r *ast.File, fset *token.FileSet) {
	fsrc := genMonoBaseStrIO
	fname := "base.io.mono.go"
	if isbytes {
		fsrc = genMonoBaseStrBytes
		fname = "base.bytes.mono.go"
	}

	fset = token.NewFileSet()
	r, err := parser.ParseFile(fset, fname, fsrc, genMonoParserMode)
	halt.onerror(err)
	r = &ast.File{
		Name: &ast.Ident{Name: "codec"},
	}
	fnames := []string{"encode.go", "decode.go", "fast-path.generated.go"}
	fnames = []string{"encode.go", "decode.go", "fast-path.not.go"} // MARKER 2025 - remove this
	for _, fname = range fnames {
		fsrc := x.file(fname)
		f, err := parser.ParseFile(fset, fname, fsrc, genMonoParserMode)
		halt.onerror(err)
		x.merge(r, f)
	}
	return
    }

type genMonoNodeStack struct {
	n []ast.Node
}

func (x *genMonoNodeStack) pop() {
	x.n = x.n[:len(x.n)-1]
}
func (x *genMonoNodeStack) push(n ast.Node) {
	x.n = append(x.n, n)
}
func (x *genMonoNodeStack) peek() (n ast.Node) {
	if pos := len(x.n) - 1; pos >= 0 {
		n = x.n[pos]
	}
	return
}

// ----

			} else if n.Tok == token.IMPORT {
				for _, v := range n.Specs {
					nn := v.(*ast.ImportSpec)
					if _, ok := x.imports[nn.Path.Value]; !ok {
						// dst.Decls = append(dst.Decls, n)
						// dd := &ast.GenDecl{Tok: n.Tok, Specs: []ast.Spec{v}}
						// dst.Decls = slices.Insert[[]ast.Decl, ast.Decl](dst.Decls, len(x.imports), dd)

						// ispec := &ast.ImportSpec{
						// 	Path: &ast.BasicLit{
						// 		Kind: token.STRING,
						// 		Value: nn.Path.Value
						x.importSpecs = append(x.importSpecs, v)
						x.imports[nn.Path.Value] = struct{}{}
					}
				}
			}


// ----

*/
