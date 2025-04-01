// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build codec.build

package codec

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strings"
)

// ----

// - [fn:monoBase] for each file of fastpath.generated.go, encode.go, decode.go:
//   - ignore init() and imports (go imports will fix it later - or we just grab all imports and let go-imports fix later)
//   - parse encode.go, decode.go, fastpath.generated.go into an AST
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

// Consequently, since there is only one function ie decByteSlice, we handle it manually.
// This method is in the type helperDecReader.

// Key constraints
//   - Every generic types MUST be defined in each file before the methods/functions
//     that reference them in the signature (receiver, parameters or results signature)
//   - No generic top level functions. Only generic methods
//     (so that things can be scoped to the driver)

var genMonoSpecialFieldTypes = []string{"helperDecReader"}

// These functions should take the address of first param when monomorphized
var genMonoSpecialFunc4Addr = []string{"decByteSlice"}

var genMonoImportsToSkip = []string{`"errors"`, `"fmt"`, `"net/rpc"`}

var genMonoRefImportsVia_ = [][2]string{
	// {"errors", "New"},
}

var genMonoCallsToSkip = []string{"callMake"}

type genMonoFieldState uint

const (
	genMonoFieldRecv genMonoFieldState = iota << 1
	genMonoFieldParamsResult
	genMonoFieldStruct
)

type genMonoImports struct {
	set   map[string]struct{}
	specs []*ast.ImportSpec
}

type genMono struct {
	files    map[string][]byte
	typParam map[string]*ast.Field
}

func (x *genMono) init() {
	x.files = make(map[string][]byte)
	x.typParam = make(map[string]*ast.Field)
}

func (x *genMono) reset() {
	clear(x.typParam)
}

func (m *genMono) do(h Handle) {
	m.reset()
	hdlFname := h.Name() + ".go"
	m.do4(h, "", []string{"encode.go", "decode.go", "fastpath.not.go", hdlFname}, ``, "fastpath")
	// any type defined in above file, should not be re-defined below
	m.do4(h, ".notfastpath", []string{"fastpath.not.go"}, ` && (notfastpath || codec.notfastpath)`, "--")
	m.do4(h, ".fastpath", []string{"fastpath.generated.go"}, ` && !notfastpath && !codec.notfastpath`, "--")
}

func (m *genMono) do4(h Handle, fnameInfx string, fnames []string, buildTagsSfx string, skipPfx string) {
	// keep m.typParams across whole call, as all others use it

	const fnameSfx = ".mono.generated.go"

	var imports = genMonoImports{set: make(map[string]struct{})}

	r1, fset := m.base(fnames, &imports, skipPfx)
	m.trFile(r1, h, true)

	r2, fset := m.base(fnames, &imports, skipPfx)
	m.trFile(r2, h, false)

	fname := h.Name() + fnameInfx + fnameSfx

	r0 := genMonoOutInit(imports.specs, fname)
	r0.Decls = append(r0.Decls, r1.Decls...)
	r0.Decls = append(r0.Decls, r2.Decls...)

	// output r1 to a file
	f, err := os.Create(fname)
	halt.onerror(err)
	defer f.Close()

	var s genMonoStrBuilder
	s.s(`//go:build !codec.notmono `).s(buildTagsSfx)
	s.s(`

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

`)
	_, err = f.Write(s.v)
	halt.onerror(err)
	err = format.Node(f, fset, r0)
	halt.onerror(err)

	// m.reset()
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

func (x *genMono) base(fnames []string, imports *genMonoImports, skipPfx string) (r *ast.File, fset *token.FileSet) {
	fset = token.NewFileSet()
	r = &ast.File{
		Name: &ast.Ident{Name: "codec"},
	}
	for _, fname := range fnames {
		fsrc := x.file(fname)
		f, err := parser.ParseFile(fset, fname, fsrc, genMonoParserMode)
		halt.onerror(err)
		x.merge(r, f, imports, skipPfx)
	}
	return
}

func (x *genMono) merge(dst, src *ast.File, imports *genMonoImports, skipPfx string) {
	// we only merge top-level methods and types
	fnIdX := func(n *ast.FuncDecl, n2 *ast.IndexExpr) (ok bool) {
		n9, ok9 := n2.Index.(*ast.Ident)
		n3, ok := n2.X.(*ast.Ident) // n3 = type name
		ok = ok && ok9 && n9.Name == "T"
		if ok {
			_, ok = x.typParam[n3.Name]
		}
		if ok {
			ok = !strings.HasPrefix(n.Name.Name, skipPfx)
		}
		if ok {
			ok = !strings.HasPrefix(n3.Name, skipPfx)
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
				ok = fnIdX(n, nn)
			case *ast.StarExpr:
				switch nn2 := nn.X.(type) {
				case *ast.IndexExpr:
					ok = fnIdX(n, nn2)
				}
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
					// always keep the typ param (since it is used, even if type is skipped
					if ok {
						// each decl will have only 1 var/type
						x.typParam[nn.Name.Name] = nn.TypeParams.List[0]
					}
					if ok {
						ok = !strings.HasPrefix(nn.Name.Name, skipPfx)
					}
					if ok {
						dst.Decls = append(dst.Decls, &ast.GenDecl{Tok: n.Tok, Specs: []ast.Spec{v}})
					}
				}
			} else if n.Tok == token.IMPORT {
				for _, v := range n.Specs {
					nn := v.(*ast.ImportSpec)
					if slices.Contains(genMonoImportsToSkip, nn.Path.Value) {
						continue
					}
					if _, ok = imports.set[nn.Path.Value]; !ok {
						imports.specs = append(imports.specs, nn)
						imports.set[nn.Path.Value] = struct{}{}
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
	tp = x.trField(n.Recv.List[0], nil, h, isbytes, genMonoFieldRecv)
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
		x.trField(v, tp, h, isbytes, genMonoFieldParamsResult)
	}
}

func (x *genMono) trStruct(r *ast.StructType, tp *ast.Field, h Handle, isbytes bool) {
	// search for fields, and update accordingly
	//   type x[T encDriver] struct { w T }
	//   var x *A[T]
	//   A[T]
	if r == nil || r.Fields == nil || len(r.Fields.List) == 0 {
		return
	}
	for _, v := range r.Fields.List {
		x.trField(v, tp, h, isbytes, genMonoFieldStruct)
	}
}

func (x *genMono) trArray(n *ast.ArrayType, tp *ast.Field, h Handle, isbytes bool) {
	sfx, _, _, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	// type fastpathEs[T encDriver] [56]fastpathE[T]
	// p := tp.Names[0].Name
	switch elt := n.Elt.(type) {
	// case *ast.InterfaceType:
	case *ast.IndexExpr:
		if elt.Index.(*ast.Ident).Name == "T" { // generic
			n.Elt = ast.NewIdent(elt.X.(*ast.Ident).Name + hnameUp + sfx)
		}
	}
}

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

func (x *genMono) trField(f *ast.Field, tpt *ast.Field, h Handle, isbytes bool, state genMonoFieldState) (tp *ast.Field) {
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
		if state == genMonoFieldRecv || nn.Name != "T" {
			return
		}
		pn = nn // "T"
		if state == genMonoFieldParamsResult {
			f.Type = &ast.StarExpr{X: pn}
		}
	}
	if pn == nil {
		return
	}

	tp = x.updateIdentForT(pn, h, tpt, isbytes, true)
	return
}

func (x *genMono) updateIdentForT(pn *ast.Ident, h Handle, tp *ast.Field,
	isbytes bool, lookupTP bool) (tp2 *ast.Field) {
	sfx, writer, reader, hnameUp := genMonoIsBytesVals(h.Name(), isbytes)
	// handle special ones e.g. helperDecReader et al
	if slices.Contains(genMonoSpecialFieldTypes, pn.Name) {
		pn.Name += sfx
		return
	}

	if pn.Name != "T" && lookupTP {
		tp = x.typParam[pn.Name]
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
	pn := v.List[0]
	if len(pn.Names) != 1 {
		return false
	}
	pnName := pn.Names[0].Name
	if pnName != "T" {
		return false
	}
	switch pn.Type.(*ast.Ident).Name {
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

type genMonoStrBuilder struct {
	v []byte
}

func (x *genMonoStrBuilder) s(v string) *genMonoStrBuilder {
	x.v = append(x.v, v...)
	return x
}

func genMonoOutInit(importSpecs []*ast.ImportSpec, fname string) (f *ast.File) {
	// ParseFile seems to skip the //go:build stanza
	// it should be written directly into the file
	var s genMonoStrBuilder
	s.s(`
package codec

import (
`)
	for _, v := range importSpecs {
		s.s("\t").s(v.Path.Value).s("\n")
	}
	s.s(")\n")
	for _, v := range genMonoRefImportsVia_ {
		s.s("var _ = ").s(v[0]).s(".").s(v[1]).s("\n")
	}
	f, err := parser.ParseFile(token.NewFileSet(), fname, s.v, genMonoParserMode)
	halt.onerror(err)
	return
}

func genMonoAll() {
	hdls := []Handle{
		(*SimpleHandle)(nil),
		(*JsonHandle)(nil),
		(*CborHandle)(nil),
		(*BincHandle)(nil),
		(*MsgpackHandle)(nil),
	}
	// hdls = []Handle{(*SimpleHandle)(nil)} // MARKER 2025 uncomment
	var m genMono
	m.init()
	for _, v := range hdls {
		m.do(v)
	}
}

// func genMonoOutInit(importSpecs []ast.Spec) (r0 *ast.File) {
// 	r0 = &ast.File{Name: &ast.Ident{Name: "codec"}}
// 	r0.Decls = append(r0.Decls, &ast.GenDecl{Tok: token.IMPORT, Specs: importSpecs})
// 	if _vd := genMonoRefImportsVia_Decl(); _vd != nil {
// 		r0.Decls = append(r0.Decls, _vd)
// 	}
// 	return
// }
//
// func genMonoRefImportsVia_Decl() (_vd *ast.GenDecl) {
// 	if len(genMonoRefImportsVia_) == 0 {
// 		return
// 	}
// 	_vd = &ast.GenDecl{Tok: token.VAR}
// 	for _, v := range genMonoRefImportsVia_ {
// 		_vs := new(ast.ValueSpec)
// 		_vs.Names = append(_vs.Names, &ast.Ident{Name: "_"})
// 		_vs.Values = append(_vs.Values, &ast.SelectorExpr{
// 			X:   &ast.Ident{Name: v[0]},
// 			Sel: &ast.Ident{Name: v[1]},
// 		})
// 		_vd.Specs = append(_vd.Specs, _vs)
// 	}
// 	return
// }
