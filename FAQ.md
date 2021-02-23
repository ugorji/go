# FAQ

## Managing Binary Size

This package adds about `900K` to any binary that depends on it.

Prior to 2019-05-16, this package could add about `2MB` to the size of
your binaries.  We have now trimmed that by 60%, and the package
contributes less than `1MB`.  This compares favorably to other packages like
`json-iterator/go` `(1MB)`.

Of that `900Kb`, about `200Kb` is from an *optional* auto-generated file: 
`fast-path.generated.go` that includes static generated implementations 
providing much faster encoding and decoding of slices and maps
containing built-in numeric, boolean, string and interface{} types.

Furthermore, you can bypass including `fast-path.generated.go` in your binary,
by building (or running tests and benchmarks) with the tag: `codec.notfastpath`.

    go install -tags codec.notfastpath
    go build -tags codec.notfastpath
    go test -tags codec.notfastpath

With the tag `codec.notfastpath`, we trim that size to about `700Kb`.

Please be aware of the following:

- At least in our representative microbenchmarks for json (for example),
  passing `codec.notfastpath` tag causes a clear performance loss (about 10%).  
  *YMMV*.
- These values were got from building the test binary that gives > 90% code coverage,
  and running `go tool nm` on it to see how much space these library symbols took.
  We consider that your use will touch much less of the library and you can be sure 
  that it will add less than 1MB to your binary.

## Resolving Module Issues

Prior to v1.1.5, `go-codec` unknowingly introduced some headaches for its
users while introducing module support. We tried to make
`github.com/ugorji/go/codec` a module. At that time, multi-repository
module support was weak, so we reverted and made `github.com/ugorji/go/`
the module.

However, folks previously used go-codec in module mode
before it formally supported modules. Eventually, different established packages
had go.mod files contain various real and pseudo versions of go-codec
which causes `go` to barf with `ambiguous import` error.

To resolve this, from v1.1.5 and up, we use a requirements cycle between
modules `github.com/ugorji/go/codec` and `github.com/ugorji/go/`,
tagging them with parallel consistent tags (`codec/vX.Y.Z and vX.Y.Z`)
to the same commit.

Fixing `ambiguous import` failure is now as simple as running

```
go get -u github.com/ugorji/go/codec@latest
```

## Running on GCCGO

You can use [gccgo](https://gcc.gnu.org/onlinedocs/gccgo/) [1](https://golang.org/doc/install/gccgo)
to build your applications that depend on this library.
This is fully tested.

There's a caveat for contributors to the library. 
If you are contributing to the library, you will need to run `./build.sh -m` at some point.
This has a failure on gccgo. Consequently, contributors must use the `gc` compiler.

