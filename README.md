[![Sourcegraph](https://sourcegraph.com/github.com/ugorji/go/-/badge.svg?v=4)](https://sourcegraph.com/github.com/ugorji/go/-/tree/codec?badge)
[![Build Status](https://travis-ci.org/ugorji/go.svg?branch=master)](https://travis-ci.org/ugorji/go)
[![codecov](https://codecov.io/gh/ugorji/go/branch/master/graph/badge.svg?v=4)](https://codecov.io/gh/ugorji/go)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ugorji/go/codec)](https://pkg.go.dev/github.com/ugorji/go/codec)
[![rcard](https://goreportcard.com/badge/github.com/ugorji/go/codec?v=4)](https://goreportcard.com/report/github.com/ugorji/go/codec)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/ugorji/go/master/LICENSE)

# go-codec

This repository contains the `go-codec` library, the `codecgen` tool and
benchmarks for comparing against other libraries.

This is a High Performance, Feature-Rich Idiomatic Go 1.4+ codec/encoding library
for binary and text formats: binc, msgpack, cbor, json and simple.

It fully supports the legacy `GOPATH` and the new `go modules` modes.

## Code Organization and Module Support

This repository consists of 4 modules:

- `github.com/ugorji/go` (requires `github.com/ugorji/go/codec`)
- `github.com/ugorji/go/codec` (requires `github.com/ugorji/go`) [README](codec/README.md)
- `github.com/ugorji/go/codec/codecgen` (requires `github.com/ugorji/go/codec`) [README](codec/codecgen/README.md)
- `github.com/ugorji/go/codec/bench` (requires `github.com/ugorji/go/codec`) [README](codec/bench/README.md)

For encoding and decoding, the `github.com/ugorji/go/codec` module is sufficient.

To install:

```
go get github.com/ugorji/go/codec
```

The other modules exist for specific uses, and all require `github.com/ugorji/go/codec`

- `github.com/ugorji/go` is here for [historical compatibility reasons, as modules was initially introduced only at repo root](https://github.com/ugorji/go/issues/299)
- `github.com/ugorji/go/codec/codecgen` generates high performance static encoders/decoders for given types
- `github.com/ugorji/go/codec/bench` benchmarks codec against other popular go libraries
