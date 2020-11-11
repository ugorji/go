# go-codec-bench

This is a comparison of different binary and text encodings.

We compare the codecs provided by github.com/ugorji/go/codec package,
against other libraries:

[github.com/ugorji/go/codec](http://github.com/ugorji/go) provides:

  - msgpack: [http://github.com/msgpack/msgpack] 
  - binc:    [http://github.com/ugorji/binc]
  - cbor:    [http://cbor.io] [http://tools.ietf.org/html/rfc7049]
  - simple: 
  - json:    [http://json.org] [http://tools.ietf.org/html/rfc7159] 

Other codecs compared include:

  - [github.com/vmihailenco/msgpack/v4](https://pkg.go.dev/github.com/vmihailenco/msgpack/v4)
  - [github.com/globalsign/mgo/bson](https://pkg.go.dev/github.com/globalsign/mgo/bson)
  - [go.mongodb.org/mongo-driver/bson](https://pkg.go.dev/go.mongodb.org/mongo-driver/bson)
  - [github.com/davecgh/go-xdr/xdr2](https://pkg.go.dev/github.com/davecgh/go-xdr/xdr2)
  - [github.com/Sereal/Sereal/Go/sereal](https://pkg.go.dev/github.com/Sereal/Sereal/Go/sereal)
  - [bitbucket.org/bodhisnarkva/cbor/go](https://pkg.go.dev/bitbucket.org/bodhisnarkva/cbor/go)
  - [github.com/tinylib/msgp](https://pkg.go.dev/github.com/tinylib/msgp)
  - [github.com/pquerna/ffjson/ffjson](https://pkg.go.dev/github.com/pquerna/ffjson/ffjson)
  - [github.com/json-iterator/go](https://pkg.go.dev/github.com/json-iterator/go)
  - [github.com/mailru/easyjson](https://pkg.go.dev/github.com/mailru/easyjson)
  - [github.com/fxamacker/cbor/v2](https://pkg.go.dev/github.com/fxamacker/cbor/v2)
  
# Data

The data being serialized is a `TestStruc` randomly generated values.
See [values_test.go](values_test.go) for the
definition of the TestStruc.

# Run Benchmarks

See [bench.sh](bench.sh)
for how to download the external libraries which we benchmark against,
generate the files for the types when needed, 
and run the suite of tests.

The 3 suite of benchmarks are

  - CodecSuite
  - XSuite
  - CodecXSuite

```
# Note that `bench.sh` may be in the codec sub-directory, and should be run from there.

# download the code and all its dependencies
./bench.sh -d

# code-generate files needed for benchmarks against ffjson, easyjson, msgp, etc
./bench.sh -c

# run the full suite of tests (not including external formats)
./bench.sh -s

# Below, see how to just run some specific suite of tests, knowing the right tags and flags ...
# See bench.sh for different iterations

# Run suite of tests in default mode (selectively using unsafe in specific areas)
go test -tags "alltests x" -bench "CodecXSuite" -benchmem 
# Run suite of tests in safe mode (no usage of unsafe)
go test -tags "alltests x safe" -bench "CodecXSuite" -benchmem 
# Run suite of tests in codecgen mode, including all tests which are generated (msgp, ffjson, etc)
go test -tags "alltests x generated" -bench "CodecXGenSuite" -benchmem 

```

# Issues

The following issues are seen currently (11/20/2014):

- _code.google.com/p/cbor/go_ fails on encoding and decoding the test struct
- _github.com/davecgh/go-xdr/xdr2_ fails on encoding and decoding the test struct
- _github.com/Sereal/Sereal/Go/sereal_ fails on decoding the serialized test struct

# Representative Benchmark Results

Please see the [benchmarking blog post for detailed representative results](http://ugorji.net/blog/benchmarking-serialization-in-go).

A snapshot of some results on my 2016 MacBook Pro is below.  
**Note: errors are truncated, and lines re-arranged, for readability**.

What you should notice:

- Results get better with codecgen, resulting in about 50% performance improvement.
  Users should carefully weigh the performance improvements against the 
  usability and binary-size increases, as performance is already extremely good 
  without the codecgen path.
  
See [bench.out.txt](bench.out.txt) for representative result from running `bench.sh` as below, as of 2020-11-11.
```sh
  ./bench.sh -z > bench.out.txt
```

*snippet of benchmark output, running without codecgen (2020-11-11)*  
*note that the first 5 are from codes (and the following are from other libraries (as named).*
```
BenchmarkCodecXSuite/options-false.../
Benchmark__Msgpack____Encode-8	     69523 ns/op	    3136 B/op	      44 allocs/op
Benchmark__Binc_______Encode-8	     73354 ns/op	    3152 B/op	      44 allocs/op
Benchmark__Simple_____Encode-8	     70840 ns/op	    3136 B/op	      44 allocs/op
Benchmark__Cbor_______Encode-8	     70575 ns/op	    3136 B/op	      44 allocs/op
Benchmark__Json_______Encode-8	    144646 ns/op	    3248 B/op	      44 allocs/op
Benchmark__Std_Json___Encode-8	    211380 ns/op	   74089 B/op	     444 allocs/op
Benchmark__Gob________Encode-8	    146538 ns/op	  169701 B/op	     591 allocs/op
Benchmark__JsonIter___Encode-8	    144049 ns/op	   53689 B/op	      72 allocs/op
Benchmark__Bson_______Encode-8	    286358 ns/op	  238101 B/op	    1095 allocs/op
Benchmark__Mgobson____Encode-8	    365772 ns/op	  292892 B/op	    1721 allocs/op
Benchmark__VMsgpack___Encode-8	    220757 ns/op	  164227 B/op	     354 allocs/op
Benchmark__Fxcbor_____Encode-8	    107005 ns/op	   49492 B/op	     320 allocs/op
Benchmark__Sereal_____Encode-8	    336930 ns/op	  263715 B/op	    3219 allocs/op

Benchmark__Msgpack____Decode-8	    176527 ns/op	   65656 B/op	     929 allocs/op
Benchmark__Binc_______Decode-8	    180256 ns/op	   68021 B/op	    1305 allocs/op
Benchmark__Simple_____Decode-8	    169090 ns/op	   65669 B/op	     929 allocs/op
Benchmark__Cbor_______Decode-8	    188041 ns/op	   65668 B/op	     929 allocs/op
Benchmark__Json_______Decode-8	    453760 ns/op	   90099 B/op	    1177 allocs/op
Benchmark__Std_Json___Decode-8	    941430 ns/op	  130725 B/op	    2961 allocs/op
Benchmark__Gob________Decode-8	    270804 ns/op	  150670 B/op	    2180 allocs/op
Benchmark__JsonIter___Decode-8	    451886 ns/op	  126328 B/op	    2486 allocs/op
Benchmark__Bson_______Decode-8	    481208 ns/op	  180314 B/op	    4256 allocs/op
Benchmark__Mgobson____Decode-8	    504214 ns/op	  161407 B/op	    6472 allocs/op
Benchmark__Fxcbor_____Decode-8	    237073 ns/op	   67270 B/op	    1299 allocs/op
```

* snippet of bench.out.txt, running with codecgen (2020-11-11) *
```
BenchmarkCodecXGenSuite/options-false.../
Benchmark__Msgpack____Encode-8	     37925 ns/op	     232 B/op	       2 allocs/op
Benchmark__Binc_______Encode-8	     42482 ns/op	     248 B/op	       2 allocs/op
Benchmark__Simple_____Encode-8	     40393 ns/op	     232 B/op	       2 allocs/op
Benchmark__Cbor_______Encode-8	     39480 ns/op	     232 B/op	       2 allocs/op
Benchmark__Json_______Encode-8	    110855 ns/op	     344 B/op	       2 allocs/op
Benchmark__Msgp_______Encode-8	     26674 ns/op	       0 B/op	       0 allocs/op
Benchmark__Easyjson___Encode-8	    111981 ns/op	   50536 B/op	      12 allocs/op
Benchmark__Ffjson_____Encode-8	    287862 ns/op	  124888 B/op	    1033 allocs/op

Benchmark__Msgpack____Decode-8	    105724 ns/op	   62053 B/op	     871 allocs/op
Benchmark__Binc_______Decode-8	    112872 ns/op	   64438 B/op	    1247 allocs/op
Benchmark__Simple_____Decode-8	    100715 ns/op	   62070 B/op	     871 allocs/op
Benchmark__Cbor_______Decode-8	    110085 ns/op	   62055 B/op	     871 allocs/op
Benchmark__Json_______Decode-8	    365746 ns/op	   85388 B/op	    1090 allocs/op
Benchmark__Msgp_______Decode-8	     66986 ns/op	   63846 B/op	     889 allocs/op
Benchmark__Easyjson___Decode-8	    346293 ns/op	   66953 B/op	     459 allocs/op
Benchmark__Ffjson_____Decode-8	    397967 ns/op	   90034 B/op	    1202 allocs/op
```
