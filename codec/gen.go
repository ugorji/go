// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build codec.gen

package codec

import (
	"bytes"
	"encoding/base32"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"

	// "ugorji.net/zz"
	"unicode"
	"unicode/utf8"
)

// ---------------------------------------------------
// codecgen supports the full cycle of reflection-based codec:
//    - RawExt
//    - Raw
//    - Extensions
//    - (Binary|Text|JSON)(Unm|M)arshal
//    - generic by-kind
//
// This means that, for dynamic things, we MUST use reflection to at least get the reflect.Type.
// In those areas, we try to only do reflection or interface-conversion when NECESSARY:
//    - Extensions, only if Extensions are configured.
//
// However, note following codecgen caveats:
//   - Canonical option.
//     If Canonical=true, codecgen'ed code may delegate encoding maps to reflection-based code.
//     This is due to the runtime work needed to marshal a map in canonical mode.
//     However, if map key is a pre-defined/builtin numeric or string type, codecgen
//     will try to write it out itself
//   - CheckCircularRef option.
//     When encoding a struct, a circular reference can lead to a stack overflow.
//     If CheckCircularRef=true, codecgen'ed code will delegate encoding structs to reflection-based code.
//   - MissingFielder implementation.
//     If a type implements MissingFielder, a Selfer is not generated (with a warning message).
//     Statically reproducing the runtime work needed to extract the missing fields and marshal them
//     along with the struct fields, while handling the Canonical=true special case, was onerous to implement.
//
// During encode/decode, Selfer takes precedence.
// A type implementing Selfer will know how to encode/decode itself statically.
//
// The following field types are supported:
//     array: [n]T
//     slice: []T
//     map: map[K]V
//     primitive: [u]int[n], float(32|64), bool, string
//     struct
//
// ---------------------------------------------------
// Note that a Selfer cannot call (e|d).(En|De)code on itself,
// as this will cause a circular reference, as (En|De)code will call Selfer methods.
// Any type that implements Selfer must implement completely and not fallback to (En|De)code.
//
// In addition, code in this file manages the generation of fast-path implementations of
// encode/decode of slices/maps of primitive keys/values.
//
// Users MUST re-generate their implementations whenever the code shape changes.
// The generated code will panic if it was generated with a version older than the supporting library.
// ---------------------------------------------------
//
// codec framework is very feature rich.
// When encoding or decoding into an interface, it depends on the runtime type of the interface.
// The type of the interface may be a named type, an extension, etc.
// Consequently, we fallback to runtime codec for encoding/decoding interfaces.
// In addition, we fallback for any value which cannot be guaranteed at runtime.
// This allows us support ANY value, including any named types, specifically those which
// do not implement our interfaces (e.g. Selfer).
//
// This explains some slowness compared to other code generation codecs (e.g. msgp).
// This reduction in speed is only seen when your refers to interfaces,
// e.g. type T struct { A interface{}; B []interface{}; C map[string]interface{} }
//
// codecgen will panic if the file was generated with an old version of the library in use.
//
// Note:
//   It was a conscious decision to have gen.go always explicitly call EncodeNil or TryDecodeAsNil.
//   This way, there isn't a function call overhead just to see that we should not enter a block of code.
//
// Note:
//   codecgen-generated code depends on the variables defined by fast-path.generated.go.
//   consequently, you cannot run with tags "codecgen codec.notfastpath".
//
// Note:
//   genInternalXXX functions are used for generating fast-path and other internally generated
//   files, and not for use in codecgen.

// Size of a struct or value is not portable across machines, especially across 32-bit vs 64-bit
// operating systems. This is due to types like int, uintptr, pointers, (and derived types like slice), etc
// which use the natural word size on those machines, which may be 4 bytes (on 32-bit) or 8 bytes (on 64-bit).
//
// Within decInferLen calls, we may generate an explicit size of the entry.
// We do this because decInferLen values are expected to be approximate,
// and serve as a good hint on the size of the elements or key+value entry.
//
// Since development is done on 64-bit machines, the sizes will be roughly correctly
// on 64-bit OS, and slightly larger than expected on 32-bit OS.
// This is ok.
//
// For reference, look for 'Size' in fast-path.go.tmpl, gen-dec-(array|map).go.tmpl and gen.go (this file).

const (
	genTopLevelVarName = "x"

	// genFastpathCanonical configures whether we support Canonical in fast path. Low savings.
	//
	// MARKER: This MUST ALWAYS BE TRUE. fast-path.go.tmpl doesn't handle it being false.
	genFastpathCanonical = true

	// genFastpathTrimTypes configures whether we trim uncommon fastpath types.
	genFastpathTrimTypes = true
)

type genStringDecAsBytes string
type genStringDecZC string

var genStringDecAsBytesTyp = reflect.TypeOf(genStringDecAsBytes(""))
var genStringDecZCTyp = reflect.TypeOf(genStringDecZC(""))
var genFormats = []string{"Json", "Cbor", "Msgpack", "Binc", "Simple"}

var (
	errGenAllTypesSamePkg        = errors.New("All types must be in the same package")
	errGenExpectArrayOrMap       = errors.New("unexpected type - expecting array/map/slice")
	errGenUnexpectedTypeFastpath = errors.New("fast-path: unexpected type - requires map or slice")

	// don't use base64, only 63 characters allowed in valid go identifiers
	// ie ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_
	//
	// don't use numbers, as a valid go identifer must start with a letter.
	genTypenameEnc = base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef")
	genQNameRegex  = regexp.MustCompile(`[A-Za-z_.]+`)
)

type genBuf struct {
	buf []byte
}

func (x *genBuf) sIf(b bool, s, t string) *genBuf {
	if b {
		x.buf = append(x.buf, s...)
	} else {
		x.buf = append(x.buf, t...)
	}
	return x
}
func (x *genBuf) s(s string) *genBuf              { x.buf = append(x.buf, s...); return x }
func (x *genBuf) b(s []byte) *genBuf              { x.buf = append(x.buf, s...); return x }
func (x *genBuf) v() string                       { return string(x.buf) }
func (x *genBuf) f(s string, args ...interface{}) { x.s(fmt.Sprintf(s, args...)) }
func (x *genBuf) reset() {
	if x.buf != nil {
		x.buf = x.buf[:0]
	}
}

// genRunner holds some state used during a Gen run.
type genRunner struct {
	w io.Writer // output
	c uint64    // counter used for generating varsfx
	f uint64    // counter used for saying false

	t  []reflect.Type   // list of types to run selfer on
	tc reflect.Type     // currently running selfer on this type
	te map[uintptr]bool // types for which the encoder has been created
	td map[uintptr]bool // types for which the decoder has been created
	tz map[uintptr]bool // types for which GenIsZero has been created

	cp string // codec import path

	im  map[string]reflect.Type // imports to add
	imn map[string]string       // package names of imports to add
	imc uint64                  // counter for import numbers

	is map[reflect.Type]struct{} // types seen during import search
	bp string                    // base PkgPath, for which we are generating for

	cpfx string // codec package prefix

	ty map[reflect.Type]struct{} // types for which GenIsZero *should* be created
	tm map[reflect.Type]struct{} // types for which enc/dec must be generated
	ts []reflect.Type            // types for which enc/dec must be generated

	xs string // top level variable/constant suffix
	hn string // fn helper type name

	ti *TypeInfos
	// rr *rand.Rand // random generator for file-specific types

	jsonOnlyWhen, toArrayWhen, omitEmptyWhen *bool

	nx bool // no extensions
}

type genIfClause struct {
	hasIf bool
}

func (g *genIfClause) end(x *genRunner) {
	if g.hasIf {
		x.line("}")
	}
}

func (g *genIfClause) c(last bool) (v string) {
	if last {
		if g.hasIf {
			v = " } else { "
		}
	} else if g.hasIf {
		v = " } else if "
	} else {
		v = "if "
		g.hasIf = true
	}
	return
}

func (x *genRunner) checkForSelfer(t reflect.Type, varname string) bool {
	// return varname != genTopLevelVarName && t != x.tc
	// the only time we checkForSelfer is if we are not at the TOP of the generated code.
	return varname != genTopLevelVarName
}

func (x *genRunner) arr2str(t reflect.Type, s string) string {
	if t.Kind() == reflect.Array {
		return s
	}
	return ""
}

func (x *genRunner) genRefPkgs(t reflect.Type) {
	if _, ok := x.is[t]; ok {
		return
	}
	x.is[t] = struct{}{}
	tpkg, tname := genImportPath(t), t.Name()
	if tpkg != "" && tpkg != x.bp && tpkg != x.cp && tname != "" && tname[0] >= 'A' && tname[0] <= 'Z' {
		if _, ok := x.im[tpkg]; !ok {
			x.im[tpkg] = t
			if idx := strings.LastIndex(tpkg, "/"); idx < 0 {
				x.imn[tpkg] = tpkg
			} else {
				x.imc++
				x.imn[tpkg] = "pkg" + strconv.FormatUint(x.imc, 10) + "_" + genGoIdentifier(tpkg[idx+1:], false)
			}
		}
	}
	switch t.Kind() {
	case reflect.Array, reflect.Slice, reflect.Ptr, reflect.Chan:
		x.genRefPkgs(t.Elem())
	case reflect.Map:
		x.genRefPkgs(t.Elem())
		x.genRefPkgs(t.Key())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if fname := t.Field(i).Name; fname != "" && fname[0] >= 'A' && fname[0] <= 'Z' {
				x.genRefPkgs(t.Field(i).Type)
			}
		}
	}
}

// sayFalse will either say "false" or use a function call that returns false.
func (x *genRunner) sayFalse() string {
	x.f++
	if x.f%2 == 0 {
		return x.hn + "False()"
	}
	return "false"
}

// sayFalse will either say "true" or use a function call that returns true.
func (x *genRunner) sayTrue() string {
	x.f++
	if x.f%2 == 0 {
		return x.hn + "True()"
	}
	return "true"
}

func (x *genRunner) varsfx() string {
	x.c++
	return strconv.FormatUint(x.c, 10)
}

func (x *genRunner) varsfxreset() {
	x.c = 0
}

func (x *genRunner) out(s string) {
	_, err := io.WriteString(x.w, s)
	genCheckErr(err)
}

func (x *genRunner) outf(s string, params ...interface{}) {
	_, err := fmt.Fprintf(x.w, s, params...)
	genCheckErr(err)
}

func (x *genRunner) line(s string) {
	x.out(s)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		x.out("\n")
	}
}

func (x *genRunner) lineIf(s string) {
	if s != "" {
		x.line(s)
	}
}

func (x *genRunner) linef(s string, params ...interface{}) {
	x.outf(s, params...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		x.out("\n")
	}
}

func (x *genRunner) genTypeName(t reflect.Type) (n string) {
	// if the type has a PkgPath, which doesn't match the current package,
	// then include it.
	// We cannot depend on t.String() because it includes current package,
	// or t.PkgPath because it includes full import path,
	//
	var ptrPfx string
	for t.Kind() == reflect.Ptr {
		ptrPfx += "*"
		t = t.Elem()
	}
	if tn := t.Name(); tn != "" {
		return ptrPfx + x.genTypeNamePrim(t)
	}
	switch t.Kind() {
	case reflect.Map:
		return ptrPfx + "map[" + x.genTypeName(t.Key()) + "]" + x.genTypeName(t.Elem())
	case reflect.Slice:
		return ptrPfx + "[]" + x.genTypeName(t.Elem())
	case reflect.Array:
		return ptrPfx + "[" + strconv.FormatInt(int64(t.Len()), 10) + "]" + x.genTypeName(t.Elem())
	case reflect.Chan:
		return ptrPfx + t.ChanDir().String() + " " + x.genTypeName(t.Elem())
	default:
		if t == intfTyp {
			return ptrPfx + "interface{}"
		} else {
			return ptrPfx + x.genTypeNamePrim(t)
		}
	}
}

func (x *genRunner) genTypeNamePrim(t reflect.Type) (n string) {
	if t.Name() == "" {
		return t.String()
	} else if genImportPath(t) == "" || genImportPath(t) == genImportPath(x.tc) {
		return t.Name()
	} else {
		return x.imn[genImportPath(t)] + "." + t.Name()
		// return t.String() // best way to get the package name inclusive
	}
}

func (x *genRunner) genZeroValueR(t reflect.Type) string {
	// if t is a named type, w
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func,
		reflect.Slice, reflect.Map, reflect.Invalid:
		return "nil"
	case reflect.Bool:
		return "false"
	case reflect.String:
		return `""`
	case reflect.Struct, reflect.Array:
		return x.genTypeName(t) + "{}"
	default: // all numbers
		return "0"
	}
}

func (x *genRunner) genMethodNameT(t reflect.Type) (s string) {
	return genMethodNameT(t, x.tc)
}

// used for chan, array, slice, map
func (x *genRunner) xtraSM(varname string, t reflect.Type, ti *typeInfo, encode, isptr bool) {
	var ptrPfx, addrPfx string
	if isptr {
		ptrPfx = "*"
	} else {
		addrPfx = "&"
	}
	if encode {
		x.linef("h.enc%s((%s%s)(%s), e)", x.genMethodNameT(t), ptrPfx, x.genTypeName(t), varname)
	} else {
		x.linef("h.dec%s((*%s)(%s%s), d)", x.genMethodNameT(t), x.genTypeName(t), addrPfx, varname)
	}
	x.registerXtraT(t, ti)
}

func (x *genRunner) registerXtraT(t reflect.Type, ti *typeInfo) {
	// recursively register the types
	tk := t.Kind()
	if tk == reflect.Ptr {
		x.registerXtraT(t.Elem(), nil)
		return
	}
	if _, ok := x.tm[t]; ok {
		return
	}

	switch tk {
	case reflect.Chan, reflect.Slice, reflect.Array, reflect.Map:
	default:
		return
	}
	// only register the type if it will not default to a fast-path
	if ti == nil {
		ti = x.ti.get(rt2id(t), t)
	}
	_, rtidu := genFastpathUnderlying(t, ti.rtid, ti)
	if _, ok := fastpathAvIndex(rtidu); ok {
		return
	}
	x.tm[t] = struct{}{}
	x.ts = append(x.ts, t)
	// check if this refers to any xtra types eg. a slice of array: add the array
	x.registerXtraT(t.Elem(), nil)
	if tk == reflect.Map {
		x.registerXtraT(t.Key(), nil)
	}
}

func (x *genRunner) encZero(t reflect.Type) {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x.line("r.EncodeInt(0)")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x.line("r.EncodeUint(0)")
	case reflect.Float32:
		x.line("r.EncodeFloat32(0)")
	case reflect.Float64:
		x.line("r.EncodeFloat64(0)")
	case reflect.Complex64:
		x.line("z.EncEncodeComplex64(0)")
	case reflect.Complex128:
		x.line("z.EncEncodeComplex128(0)")
	case reflect.Bool:
		x.line("r.EncodeBool(false)")
	case reflect.String:
		x.linef(`r.EncodeString("")`)
	default:
		x.line("r.EncodeNil()")
	}
}

func genOmitEmptyLinePreChecks(varname string, t reflect.Type, si *structFieldInfo, omitline *genBuf, oneLevel bool) (t2 reflect.StructField) {
	// xdebug2f("calling genOmitEmptyLinePreChecks on: %v", t)
	t2typ := t
	varname3 := varname
	// go through the loop, record the t2 field explicitly,
	// and gather the omit line if embedded in pointers.
	fullpath := si.path.fullpath()
	for i, path := range fullpath {
		for t2typ.Kind() == reflect.Ptr {
			t2typ = t2typ.Elem()
		}
		t2 = t2typ.Field(int(path.index))
		t2typ = t2.Type
		varname3 = varname3 + "." + t2.Name
		// do not include actual field in the omit line.
		// that is done subsequently (right after - below).
		if i+1 < len(fullpath) && t2typ.Kind() == reflect.Ptr {
			omitline.s(varname3).s(" != nil && ")
		}
		if oneLevel {
			break
		}
	}
	return
}

// --------

type fastpathGenV struct {
	// fastpathGenV is either a primitive (Primitive != "") or a map (MapKey != "") or a slice
	MapKey      string
	Elem        string
	Primitive   string
	Size        int
	NoCanonical bool
}

func (x *genRunner) newFastpathGenV(t reflect.Type) (v fastpathGenV) {
	v.NoCanonical = !genFastpathCanonical
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		te := t.Elem()
		v.Elem = x.genTypeName(te)
		v.Size = int(te.Size())
	case reflect.Map:
		te := t.Elem()
		tk := t.Key()
		v.Elem = x.genTypeName(te)
		v.MapKey = x.genTypeName(tk)
		v.Size = int(te.Size() + tk.Size())
	default:
		halt.onerror(errGenUnexpectedTypeFastpath)
	}
	return
}

func (x *fastpathGenV) MethodNamePfx(prefix string, prim bool) string {
	var name []byte
	if prefix != "" {
		name = append(name, prefix...)
	}
	if prim {
		name = append(name, genTitleCaseName(x.Primitive)...)
	} else {
		if x.MapKey == "" {
			name = append(name, "Slice"...)
		} else {
			name = append(name, "Map"...)
			name = append(name, genTitleCaseName(x.MapKey)...)
		}
		name = append(name, genTitleCaseName(x.Elem)...)
	}
	return string(name)
}

// genImportPath returns import path of a non-predeclared named typed, or an empty string otherwise.
//
// This handles the misbehaviour that occurs when 1.5-style vendoring is enabled,
// where PkgPath returns the full path, including the vendoring pre-fix that should have been stripped.
// We strip it here.
func genImportPath(t reflect.Type) (s string) {
	s = t.PkgPath()
	// HACK: always handle vendoring. It should be typically on in go 1.6, 1.7
	s = genStripVendor(s)
	return
}

// A go identifier is (letter|_)[letter|number|_]*
func genGoIdentifier(s string, checkFirstChar bool) string {
	b := make([]byte, 0, len(s))
	t := make([]byte, 4)
	var n int
	for i, r := range s {
		if checkFirstChar && i == 0 && !unicode.IsLetter(r) {
			b = append(b, '_')
		}
		// r must be unicode_letter, unicode_digit or _
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			n = utf8.EncodeRune(t, r)
			b = append(b, t[:n]...)
		} else {
			b = append(b, '_')
		}
	}
	return string(b)
}

func genNonPtr(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func genFastpathUnderlying(t reflect.Type, rtid uintptr, ti *typeInfo) (tu reflect.Type, rtidu uintptr) {
	tu = t
	rtidu = rtid
	if ti.flagHasPkgPath {
		tu = ti.fastpathUnderlying
		rtidu = rt2id(tu)
	}
	return
}

func genTitleCaseName(s string) string {
	switch s {
	case "interface{}", "interface {}":
		return "Intf"
	case "[]byte", "[]uint8", "bytes":
		return "Bytes"
	default:
		return strings.ToUpper(s[0:1]) + s[1:]
	}
}

func genMethodNameT(t reflect.Type, tRef reflect.Type) (n string) {
	var ptrPfx string
	for t.Kind() == reflect.Ptr {
		ptrPfx += "Ptrto"
		t = t.Elem()
	}
	tstr := t.String()
	if tn := t.Name(); tn != "" {
		if tRef != nil && genImportPath(t) == genImportPath(tRef) {
			return ptrPfx + tn
		} else {
			if genQNameRegex.MatchString(tstr) {
				return ptrPfx + strings.Replace(tstr, ".", "_", 1000)
			} else {
				return ptrPfx + genCustomTypeName(tstr)
			}
		}
	}
	switch t.Kind() {
	case reflect.Map:
		return ptrPfx + "Map" + genMethodNameT(t.Key(), tRef) + genMethodNameT(t.Elem(), tRef)
	case reflect.Slice:
		return ptrPfx + "Slice" + genMethodNameT(t.Elem(), tRef)
	case reflect.Array:
		return ptrPfx + "Array" + strconv.FormatInt(int64(t.Len()), 10) + genMethodNameT(t.Elem(), tRef)
	case reflect.Chan:
		var cx string
		switch t.ChanDir() {
		case reflect.SendDir:
			cx = "ChanSend"
		case reflect.RecvDir:
			cx = "ChanRecv"
		default:
			cx = "Chan"
		}
		return ptrPfx + cx + genMethodNameT(t.Elem(), tRef)
	default:
		if t == intfTyp {
			return ptrPfx + "Interface"
		} else {
			if tRef != nil && genImportPath(t) == genImportPath(tRef) {
				if t.Name() != "" {
					return ptrPfx + t.Name()
				} else {
					return ptrPfx + genCustomTypeName(tstr)
				}
			} else {
				// best way to get the package name inclusive
				if t.Name() != "" && genQNameRegex.MatchString(tstr) {
					return ptrPfx + strings.Replace(tstr, ".", "_", 1000)
				} else {
					return ptrPfx + genCustomTypeName(tstr)
				}
			}
		}
	}
}

// genCustomNameForType base32 encodes the t.String() value in such a way
// that it can be used within a function name.
func genCustomTypeName(tstr string) string {
	len2 := genTypenameEnc.EncodedLen(len(tstr))
	bufx := make([]byte, len2)
	genTypenameEnc.Encode(bufx, []byte(tstr))
	for i := len2 - 1; i >= 0; i-- {
		if bufx[i] == '=' {
			len2--
		} else {
			break
		}
	}
	return string(bufx[:len2])
}

func genIsImmutable(t reflect.Type) (v bool) {
	return scalarBitset.isset(byte(t.Kind()))
}

type genInternal struct {
	Values  []fastpathGenV
	Formats []string
}

func (x genInternal) FastpathLen() (l int) {
	for _, v := range x.Values {
		// if v.Primitive == "" && !(v.MapKey == "" && v.Elem == "uint8") {
		if v.Primitive == "" {
			l++
		}
	}
	return
}

func genInternalZeroValue(s string) string {
	switch s {
	case "interface{}", "interface {}":
		return "nil"
	case "[]byte", "[]uint8", "bytes":
		return "nil"
	case "bool":
		return "false"
	case "string":
		return `""`
	default:
		return "0"
	}
}

var genInternalNonZeroValueIdx [6]uint64
var genInternalNonZeroValueStrs = [...][6]string{
	{`"string-is-an-interface-1"`, "true", `"some-string-1"`, `[]byte("some-string-1")`, "11.1", "111"},
	{`"string-is-an-interface-2"`, "false", `"some-string-2"`, `[]byte("some-string-2")`, "22.2", "77"},
	{`"string-is-an-interface-3"`, "true", `"some-string-3"`, `[]byte("some-string-3")`, "33.3e3", "127"},
}

// Note: last numbers must be in range: 0-127 (as they may be put into a int8, uint8, etc)

func genInternalNonZeroValue(s string) string {
	var i int
	switch s {
	case "interface{}", "interface {}":
		i = 0
	case "bool":
		i = 1
	case "string":
		i = 2
	case "bytes", "[]byte", "[]uint8":
		i = 3
	case "float32", "float64", "float", "double", "complex", "complex64", "complex128":
		i = 4
	default:
		i = 5
	}
	genInternalNonZeroValueIdx[i]++
	idx := genInternalNonZeroValueIdx[i]
	slen := uint64(len(genInternalNonZeroValueStrs))
	return genInternalNonZeroValueStrs[idx%slen][i] // return string, to remove ambiguity
}

// Note: used for fastpath only
func genInternalEncCommandAsString(s string, vname string) string {
	switch s {
	case "uint64":
		return "e.e.EncodeUint(" + vname + ")"
	case "uint", "uint8", "uint16", "uint32":
		return "e.e.EncodeUint(uint64(" + vname + "))"
	case "int64":
		return "e.e.EncodeInt(" + vname + ")"
	case "int", "int8", "int16", "int32":
		return "e.e.EncodeInt(int64(" + vname + "))"
	case "[]byte", "[]uint8", "bytes":
		return "e.e.EncodeStringBytesRaw(" + vname + ")"
	case "string":
		return "e.e.EncodeString(" + vname + ")"
	case "float32":
		return "e.e.EncodeFloat32(" + vname + ")"
	case "float64":
		return "e.e.EncodeFloat64(" + vname + ")"
	case "bool":
		return "e.e.EncodeBool(" + vname + ")"
	// case "symbol":
	// 	return "e.e.EncodeSymbol(" + vname + ")"
	default:
		return "e.encode(" + vname + ")"
	}
}

// Note: used for fastpath only
func genInternalDecCommandAsString(s string, mapkey bool) string {
	switch s {
	case "uint":
		return "uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))"
	case "uint8":
		return "uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))"
	case "uint16":
		return "uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))"
	case "uint32":
		return "uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))"
	case "uint64":
		return "d.d.DecodeUint64()"
	case "uintptr":
		return "uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))"
	case "int":
		return "int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))"
	case "int8":
		return "int8(chkOvf.IntV(d.d.DecodeInt64(), 8))"
	case "int16":
		return "int16(chkOvf.IntV(d.d.DecodeInt64(), 16))"
	case "int32":
		return "int32(chkOvf.IntV(d.d.DecodeInt64(), 32))"
	case "int64":
		return "d.d.DecodeInt64()"

	case "string":
		// if mapkey {
		// 	return "d.stringZC(d.d.DecodeStringAsBytes())"
		// }
		// return "string(d.d.DecodeStringAsBytes())"
		return "d.stringZC(d.d.DecodeStringAsBytes())"
	case "[]byte", "[]uint8", "bytes":
		return "d.d.DecodeBytes(zeroByteSlice)"
	case "float32":
		return "float32(d.d.DecodeFloat32())"
	case "float64":
		return "d.d.DecodeFloat64()"
	case "complex64":
		return "complex(d.d.DecodeFloat32(), 0)"
	case "complex128":
		return "complex(d.d.DecodeFloat64(), 0)"
	case "bool":
		return "d.d.DecodeBool()"
	default:
		halt.onerror(errors.New("gen internal: unknown type for decode: " + s))
	}
	return ""
}

// func genInternalSortType(s string, elem bool) string {
// 	for _, v := range [...]string{
// 		"int",
// 		"uint",
// 		"float",
// 		"bool",
// 		"string",
// 		"bytes", "[]uint8", "[]byte",
// 	} {
// 		if v == "[]byte" || v == "[]uint8" {
// 			v = "bytes"
// 		}
// 		if strings.HasPrefix(s, v) {
// 			if v == "int" || v == "uint" || v == "float" {
// 				v += "64"
// 			}
// 			if elem {
// 				return v
// 			}
// 			return v + "Slice"
// 		}
// 	}
// 	halt.onerror(errors.New("sorttype: unexpected type: " + s))
// }

func genInternalSortType(s string, elem bool) string {
	if elem {
		return s
	}
	return s + "Slice"
}

// MARKER: keep in sync with codecgen/gen.go
func genStripVendor(s string) string {
	// HACK: Misbehaviour occurs in go 1.5. May have to re-visit this later.
	// if s contains /vendor/ OR startsWith vendor/, then return everything after it.
	const vendorStart = "vendor/"
	const vendorInline = "/vendor/"
	if i := strings.LastIndex(s, vendorInline); i >= 0 {
		s = s[i+len(vendorInline):]
	} else if strings.HasPrefix(s, vendorStart) {
		s = s[len(vendorStart):]
	}
	return s
}

// var genInternalMu sync.Mutex
var genInternalV = genInternal{}
var genInternalTmplFuncs template.FuncMap
var genInternalOnce sync.Once

func genInternalInit() {
	wordSizeBytes := int(intBitsize) / 8

	typesizes := map[string]int{
		"interface{}": 2 * wordSizeBytes,
		"string":      2 * wordSizeBytes,
		"[]byte":      3 * wordSizeBytes,
		"uint":        1 * wordSizeBytes,
		"uint8":       1,
		"uint16":      2,
		"uint32":      4,
		"uint64":      8,
		"uintptr":     1 * wordSizeBytes,
		"int":         1 * wordSizeBytes,
		"int8":        1,
		"int16":       2,
		"int32":       4,
		"int64":       8,
		"float32":     4,
		"float64":     8,
		"complex64":   8,
		"complex128":  16,
		"bool":        1,
	}

	// keep as slice, so it is in specific iteration order.
	// Initial order was uint64, string, interface{}, int, int64, ...

	var types = [...]string{
		"interface{}",
		"string",
		"[]byte",
		"float32",
		"float64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"bool",
	}

	var primitivetypes, slicetypes, mapkeytypes, mapvaltypes []string

	primitivetypes = types[:]

	slicetypes = types[:]
	mapkeytypes = types[:]
	mapvaltypes = types[:]

	if genFastpathTrimTypes {
		// Note: we only create fast-paths for commonly used types.
		// Consequently, things like int8, uint16, uint, etc are commented out.
		slicetypes = []string{
			"interface{}",
			"string",
			"[]byte",
			"float32",
			"float64",
			"uint8", // keep fast-path, so it doesn't have to go through reflection
			"uint64",
			"int",
			"int32", // rune
			"int64",
			"bool",
		}
		mapkeytypes = []string{
			"string",
			"uint8",  // byte
			"uint64", // used for keys
			"int",    // default number key
			"int32",  // rune
		}
		mapvaltypes = []string{
			"interface{}",
			"string",
			"[]byte",
			"uint8",  // byte
			"uint64", // used for keys, etc
			"int",    // default number
			"int32",  // rune (mostly used for unicode)
			"float64",
			"bool",
		}
	}

	var gt = genInternal{Formats: genFormats}

	// For each slice or map type, there must be a (symmetrical) Encode and Decode fast-path function

	for _, s := range primitivetypes {
		gt.Values = append(gt.Values,
			fastpathGenV{Primitive: s, Size: typesizes[s], NoCanonical: !genFastpathCanonical})
	}
	for _, s := range slicetypes {
		gt.Values = append(gt.Values,
			fastpathGenV{Elem: s, Size: typesizes[s], NoCanonical: !genFastpathCanonical})
	}
	for _, s := range mapkeytypes {
		for _, ms := range mapvaltypes {
			gt.Values = append(gt.Values,
				fastpathGenV{MapKey: s, Elem: ms, Size: typesizes[s] + typesizes[ms], NoCanonical: !genFastpathCanonical})
		}
	}

	funcs := make(template.FuncMap)
	// funcs["haspfx"] = strings.HasPrefix
	funcs["encmd"] = genInternalEncCommandAsString
	funcs["decmd"] = genInternalDecCommandAsString
	funcs["zerocmd"] = genInternalZeroValue
	funcs["nonzerocmd"] = genInternalNonZeroValue
	funcs["hasprefix"] = strings.HasPrefix
	funcs["sorttype"] = genInternalSortType

	genInternalV = gt
	genInternalTmplFuncs = funcs
}

// genInternalGoFile is used to generate source files from templates.
func genInternalGoFile(r io.Reader, w io.Writer) (err error) {
	genInternalOnce.Do(genInternalInit)

	gt := genInternalV

	t := template.New("").Funcs(genInternalTmplFuncs)

	tmplstr, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	if t, err = t.Parse(string(tmplstr)); err != nil {
		return
	}

	var out bytes.Buffer
	err = t.Execute(&out, gt)
	if err != nil {
		return
	}

	bout, err := format.Source(out.Bytes())
	if err != nil {
		w.Write(out.Bytes()) // write out if error, so we can still see.
		// w.Write(bout) // write out if error, as much as possible, so we can still see.
		return
	}
	w.Write(bout)
	return
}

func genTypeForShortName(s string) string {
	switch s {
	case "time":
		return "time.Time"
	case "bytes":
		return "[]byte"
	}
	return s
}

func genArgs(args ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); {
		m[args[i].(string)] = args[i+1]
		i += 2
	}
	return m
}

func genEndsWith(s0 string, sn ...string) bool {
	for _, s := range sn {
		if strings.HasSuffix(s0, s) {
			return true
		}
	}
	return false
}

func genCheckErr(err error) {
	halt.onerror(err)
}

func genRunTmpl2Go(fnameIn, fnameOut string) {
	// println("____ " + fnameIn + " --> " + fnameOut + " ______")
	fin, err := os.Open(fnameIn)
	genCheckErr(err)
	defer fin.Close()
	fout, err := os.Create(fnameOut)
	genCheckErr(err)
	defer fout.Close()
	err = genInternalGoFile(fin, fout)
	genCheckErr(err)
}

// --- some methods here for other types, which are only used in codecgen

// depth returns number of valid nodes in the hierachy
func (path *structFieldInfoPathNode) root() *structFieldInfoPathNode {
TOP:
	if path.parent != nil {
		path = path.parent
		goto TOP
	}
	return path
}

func (path *structFieldInfoPathNode) fullpath() (p []*structFieldInfoPathNode) {
	// this method is mostly called by a command-line tool - it's not optimized, and that's ok.
	// it shouldn't be used in typical runtime use - as it does unnecessary allocation.
	d := path.depth()
	p = make([]*structFieldInfoPathNode, d)
	for d--; d >= 0; d-- {
		p[d] = path
		path = path.parent
	}
	return
}
