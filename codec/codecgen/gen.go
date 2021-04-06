// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// codecgen generates static implementations of the encoder and decoder functions
// for a given type, bypassing reflection, and giving some performance benefits in terms of
// wall and cpu time, and memory usage.
//
// Benchmarks (as of Dec 2018) show that codecgen gives about
//
//   - for binary formats (cbor, etc): 25% on encoding and 30% on decoding to/from []byte
//   - for text formats (json, etc): 15% on encoding and 25% on decoding to/from []byte
//
// Note that (as of Dec 2018) codecgen completely ignores
//
// - MissingFielder interface
//   (if you types implements it, codecgen ignores that)
// - decode option PreferArrayOverSlice
//   (we cannot dynamically create non-static arrays without reflection)
//
// In explicit package terms: codecgen generates codec.Selfer implementations for a set of types.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// MARKER: keep in sync with ../gen.go (genVersion) and ../go.mod (module version)
const (
	codecgenModuleVersion = `1.2.6` // default version - overridden if available via go.mod
	minimumCodecVersion   = `1.2.6`
	genVersion            = 25
)

const genCodecPkg = "codec1978" // MARKER: keep in sync with ../gen.go

const genFrunMainTmpl = `// +build ignore

// Code generated - temporary main package for codecgen - DO NOT EDIT.

package main
{{ if .Types }}import "{{ .ImportPath }}"{{ end }}
func main() {
	{{ $.PackageName }}.CodecGenTempWrite{{ .RandString }}()
}
`

// const genFrunPkgTmpl = `// +build codecgen
const genFrunPkgTmpl = `

// Code generated - temporary package for codecgen - DO NOT EDIT.

package {{ $.PackageName }}

import (
	{{ if not .CodecPkgFiles }}{{ .CodecPkgName }} "{{ .CodecImportPath }}"{{ end }}
	"os"
	"reflect"
	"bytes"
	"strings"
	"go/format"
	"fmt"
)

func codecGenBoolPtr(b bool) *bool {
	return &b
}

func CodecGenTempWrite{{ .RandString }}() {
	os.Remove("{{ .OutFile }}")
	fout, err := os.Create("{{ .OutFile }}")
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	var bJO, bS2A, bOE *bool
	_, _, _ = bJO, bS2A, bOE
	{{ if .JsonOnly }}bJO = codecGenBoolPtr({{ boolvar .JsonOnly }}){{end}}
	{{ if .StructToArrayAlways }}bS2A = codecGenBoolPtr({{ boolvar .StructToArrayAlways }}){{end}}
	{{ if .OmitEmptyAlways }}bOE = codecGenBoolPtr({{ boolvar .OmitEmptyAlways }}){{end}}

	var typs []reflect.Type
	var typ reflect.Type
	var numfields int
{{ range $index, $element := .Types }}
	var t{{ $index }} {{ . }}
typ = reflect.TypeOf(t{{ $index }})
	typs = append(typs, typ)
	if typ.Kind() == reflect.Struct { numfields += typ.NumField() } else { numfields += 1 }
{{ end }}

	// println("initializing {{ .OutFile }}, buf size: {{ .AllFilesSize }}*16",
	// 	{{ .AllFilesSize }}*16, "num fields: ", numfields)
	var out = bytes.NewBuffer(make([]byte, 0, numfields*1024)) // {{ .AllFilesSize }}*16
	var warnings = {{ if not .CodecPkgFiles }}{{ .CodecPkgName }}.{{ end }}Gen(out,
		"{{ .BuildTag }}", "{{ .PackageName }}", "{{ .RandString }}", {{ .NoExtensions }},
		bJO, bS2A, bOE,
		{{ if not .CodecPkgFiles }}{{ .CodecPkgName }}.{{ end }}NewTypeInfos(strings.Split("{{ .StructTags }}", ",")),
		 typs...)

	for _, warning := range warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}

	bout, err := format.Source(out.Bytes())
	// println("... lengths: before formatting: ", len(out.Bytes()), ", after formatting", len(bout))
	if err != nil {
		fout.Write(out.Bytes())
		panic(err)
	}
	fout.Write(bout)
}

`

var genFuncs template.FuncMap

func init() {
	genFuncs = make(template.FuncMap)
	genFuncs["boolvar"] = func(x *bool) bool { return *x }
}

type regexFlagValue struct {
	v *regexp.Regexp
}

func (v *regexFlagValue) Set(s string) (err error) {
	// fmt.Printf("calling regexFlagValue.Set with %s\n", s)
	v.v, err = regexp.Compile(s)
	return
}
func (v *regexFlagValue) Get() interface{} { return v.v }
func (v *regexFlagValue) String() string   { return fmt.Sprintf("%v", v.v) }

// boolFlagValue can be set to true or false, or unset as nil (default)
type boolFlagValue struct {
	v *bool
}

func (v *boolFlagValue) Set(s string) (err error) {
	// fmt.Printf("calling boolFlagValue.Set with %s\n", s)
	b, err := strconv.ParseBool(s)
	if err == nil {
		v.v = &b
	}
	return
}
func (v *boolFlagValue) Get() interface{} { return v.v }
func (v *boolFlagValue) String() string   { return fmt.Sprintf("%#v", v.v) }

type mainCfg struct {
	CodecPkgName    string
	CodecImportPath string
	ImportPath      string
	OutFile         string
	PackageName     string
	RandString      string
	BuildTag        string
	StructTags      string
	Types           []string
	AllFilesSize    int64
	CodecPkgFiles   bool
	NoExtensions    bool

	JsonOnly            *bool
	StructToArrayAlways *bool
	OmitEmptyAlways     *bool

	uid          int64
	goRunTags    string
	regexName    *regexp.Regexp
	notRegexName *regexp.Regexp

	keepTempFile bool // !deleteTempFile
}

// Generate is given a list of *.go files to parse, and an output file (fout).
//
// It finds all types T in the files, and it creates 2 tmp files (frun).
//   - main package file passed to 'go run'
//   - package level file which calls *genRunner.Selfer to write Selfer impls for each T.
// We use a package level file so that it can reference unexported types in the package being worked on.
// Tool then executes: "go run __frun__" which creates fout.
// fout contains Codec(En|De)codeSelf implementations for every type T.
//
func mainGen(tv *mainCfg, infiles ...string) (err error) {
	// For each file, grab AST, find each type, and write a call to it.
	if len(infiles) == 0 {
		return
	}
	if tv.CodecImportPath == "" {
		return errors.New("codec package path cannot be blank")
	}
	if tv.OutFile == "" {
		return errors.New("outfile cannot be blank")
	}
	if tv.regexName == nil {
		tv.regexName = regexp.MustCompile(".*")
	}
	if tv.notRegexName == nil {
		tv.notRegexName = regexp.MustCompile("^$")
	}
	if tv.uid < 0 {
		tv.uid = -tv.uid
	} else if tv.uid == 0 {
		rr := rand.New(rand.NewSource(time.Now().UnixNano()))
		tv.uid = 101 + rr.Int63n(9777)
	}

	// We have to parse dir for package, before opening the temp file for writing (else ImportDir fails).
	// Also, ImportDir(...) must take an absolute path.
	lastdir := filepath.Dir(tv.OutFile)
	absdir, err := filepath.Abs(lastdir)
	if err != nil {
		return
	}
	importPath, err := pkgPath(absdir)
	if err != nil {
		return
	}

	tv.CodecPkgName = genCodecPkg
	tv.RandString = strconv.FormatInt(tv.uid, 10)

	tv.ImportPath = importPath
	if tv.ImportPath == tv.CodecImportPath {
		tv.CodecPkgFiles = true
		tv.CodecPkgName = "codec"
	} else {
		// HACK: always handle vendoring. It should be typically on in go 1.6, 1.7
		tv.ImportPath = genStripVendor(tv.ImportPath)
	}
	astfiles := make([]*ast.File, len(infiles))
	var fi os.FileInfo
	for i, infile := range infiles {
		if filepath.Dir(infile) != lastdir {
			err = errors.New("all input files must all be in same directory as output file")
			return
		}
		if fi, err = os.Stat(infile); err != nil {
			return
		}
		tv.AllFilesSize += fi.Size()

		fset := token.NewFileSet()
		astfiles[i], err = parser.ParseFile(fset, infile, nil, 0)
		if err != nil {
			return
		}
		if i == 0 {
			tv.PackageName = astfiles[i].Name.Name
			if tv.PackageName == "main" {
				// codecgen cannot be run on types in the 'main' package.
				// A temporary 'main' package must be created, and should reference the fully built
				// package containing the types.
				// Also, the temporary main package will conflict with the main package which already has a main method.
				err = errors.New("codecgen cannot be run on types in the 'main' package")
				return
			}
		}
	}

	// keep track of types with selfer methods
	// selferMethods := []string{"CodecEncodeSelf", "CodecDecodeSelf"}
	selferEncTyps := make(map[string]bool)
	selferDecTyps := make(map[string]bool)
	for _, f := range astfiles {
		for _, d := range f.Decls {
			// if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv != nil && fd.Recv.NumFields() == 1 {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv != nil && len(fd.Recv.List) == 1 {
				recvType := fd.Recv.List[0].Type
				if ptr, ok := recvType.(*ast.StarExpr); ok {
					recvType = ptr.X
				}
				if id, ok := recvType.(*ast.Ident); ok {
					switch fd.Name.Name {
					case "CodecEncodeSelf":
						selferEncTyps[id.Name] = true
					case "CodecDecodeSelf":
						selferDecTyps[id.Name] = true
					}
				}
			}
		}
	}

	// now find types
	for _, f := range astfiles {
		for _, d := range f.Decls {
			if gd, ok := d.(*ast.GenDecl); ok {
				for _, dd := range gd.Specs {
					if td, ok := dd.(*ast.TypeSpec); ok {
						// if len(td.Name.Name) == 0 || td.Name.Name[0] > 'Z' || td.Name.Name[0] < 'A' {
						if len(td.Name.Name) == 0 {
							continue
						}

						// only generate for:
						//   struct: StructType
						//   primitives (numbers, bool, string): Ident
						//   map: MapType
						//   slice, array: ArrayType
						//   chan: ChanType
						// do not generate:
						//   FuncType, InterfaceType, StarExpr (ptr), etc
						//
						// We generate for all these types (not just structs), because they may be a field
						// in another struct which doesn't have codecgen run on it, and it will be nice
						// to take advantage of the fact that the type is a Selfer.
						switch td.Type.(type) {
						case *ast.StructType, *ast.Ident, *ast.MapType, *ast.ArrayType, *ast.ChanType:
							// only add to tv.Types iff
							//   - it matches per the -r parameter
							//   - it doesn't match per the -nr parameter
							//   - it doesn't have any of the Selfer methods in the file
							if tv.regexName.FindStringIndex(td.Name.Name) != nil &&
								tv.notRegexName.FindStringIndex(td.Name.Name) == nil &&
								!selferEncTyps[td.Name.Name] &&
								!selferDecTyps[td.Name.Name] {
								tv.Types = append(tv.Types, td.Name.Name)
							}
						}
					}
				}
			}
		}
	}

	if len(tv.Types) == 0 {
		return
	}

	// we cannot use ioutil.TempFile, because we cannot guarantee the file suffix (.go).
	// Also, we cannot create file in temp directory,
	// because go run will not work (as it needs to see the types here).
	// Consequently, create the temp file in the current directory, and remove when done.

	// frun, err = ioutil.TempFile("", "codecgen-")
	// frunName := filepath.Join(os.TempDir(), "codecgen-"+strconv.FormatInt(time.Now().UnixNano(), 10)+".go")

	frunMainName := filepath.Join(lastdir, "codecgen-main-"+tv.RandString+".generated.go")
	frunPkgName := filepath.Join(lastdir, "codecgen-pkg-"+tv.RandString+".generated.go")

	// var frunMain, frunPkg *os.File
	if _, err = gen1(frunMainName, genFrunMainTmpl, &tv); err != nil {
		return
	}
	if _, err = gen1(frunPkgName, genFrunPkgTmpl, &tv); err != nil {
		return
	}

	// remove outfile, so "go run ..." will not think that types in outfile already exist.
	os.Remove(tv.OutFile)

	// execute go run
	// MARKER: it must be run with codec.safe tag, so that we don't use wrong optimizations that only make sense at runtime.
	// e.g. compositeUnderlyingType will cause generated code to have bad conversions for defined composite types.
	cmd := exec.Command("go", "run", "-tags", "codecgen.exec codec.safe "+tv.goRunTags, frunMainName) //, frunPkg.Name())
	cmd.Dir = lastdir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("error running 'go run %s': %v, console: %s", frunMainName, err, buf.Bytes())
		return
	}
	os.Stdout.Write(buf.Bytes())

	// only delete these files if codecgen ran successfully.
	// if unsuccessful, these files are here for diagnosis.
	if !tv.keepTempFile {
		os.Remove(frunMainName)
		os.Remove(frunPkgName)
	}

	return
}

func gen1(frunName, tmplStr string, tv interface{}) (frun *os.File, err error) {
	os.Remove(frunName)
	if frun, err = os.Create(frunName); err != nil {
		return
	}
	defer frun.Close()

	t := template.New("").Funcs(genFuncs)
	if t, err = t.Parse(tmplStr); err != nil {
		return
	}
	bw := bufio.NewWriter(frun)
	if err = t.Execute(bw, tv); err != nil {
		bw.Flush()
		return
	}
	if err = bw.Flush(); err != nil {
		return
	}
	return
}

// MARKER: keep in sync with ../gen.go
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

func main() {
	var unusedBool bool
	var printVersion bool
	var g mainCfg
	var r1, r2 regexFlagValue
	var b1, b2, b3 boolFlagValue

	flag.BoolVar(&printVersion, "version", false, "show version information")
	flag.StringVar(&g.OutFile, "o", "", "`output file` that contains generated type")
	flag.StringVar(&g.CodecImportPath, "c", genCodecPath, "`codec import path` useful when building against a mirror or fork")
	flag.StringVar(&g.BuildTag, "t", "", "`build tag` to put into generated file")

	flag.Var(&r1, "r", "`regex` for type names to match - defaults to '.*' (all)")
	flag.Var(&r2, "nr", "`regex` for type names to exclude - defaults to the '^$' (none)")

	flag.Var(&b1, "j", "if set to `true or false`, generated file supports 'json only' or omits json specific optimizations")
	flag.Var(&b2, "ta", "if set to `true or false`, support 'to_array' always or never, regardless of Handle or struct tags")
	flag.Var(&b3, "oe", "if set to `true or false`, 'omit empty' always or never, regardless of struct tags")

	flag.StringVar(&g.goRunTags, "rt", "", "`runtime tags` for go run to select which files to use during codecgen execution")
	flag.StringVar(&g.StructTags, "st", "codec,json", "`struct tag keys` to introspect")
	flag.BoolVar(&g.keepTempFile, "x", false, "set to `true or false` to 'keep' temp file created during codecgen execution - for debugging sessions")
	flag.BoolVar(&unusedBool, "u", false, "set to `true or false` to 'allow unsafe' use. **Deprecated and Ignored: kept for backwards compatibility**")
	flag.Int64Var(&g.uid, "d", 0, "`random integer identifier` for use in generated code, to prevent excessive churn")
	flag.BoolVar(&g.NoExtensions, "nx", false, "set to `true or false` to elide/ignore checking for extensions (if you do not use extensions)")

	flag.Parse()

	if printVersion {
		var modVersion string = codecgenModuleVersion
		if bi, ok := debug.ReadBuildInfo(); ok {
			if modVersion = bi.Main.Version; len(modVersion) > 0 && modVersion[0] == 'v' {
				modVersion = modVersion[1:]
			}
		}
		fmt.Printf("codecgen v%s (internal version %d) works with %s library v%s +\n",
			modVersion, genVersion, genCodecPath, minimumCodecVersion)
		return
	}

	g.JsonOnly = b1.v
	g.StructToArrayAlways = b2.v
	g.OmitEmptyAlways = b3.v

	g.regexName = r1.v
	g.notRegexName = r2.v

	// fmt.Printf("jsonOnly: %v, StructToArrayAlways: %v, OmitEmptyAlways: %v\n", g.JsonOnly, g.StructToArrayAlways, g.OmitEmptyAlways)
	err := mainGen(&g, flag.Args()...)
	// err := mainGen(*o, *t, *c, *d, *rt, *st,
	// 	regexp.MustCompile(*r), regexp.MustCompile(*nr), !*x, *nx, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "codecgen error: %v\n", err)
		os.Exit(1)
	}
}
