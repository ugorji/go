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
See https://github.com/ugorji/go-codec-bench/blob/master/codec/values_test.go for the
definition of the TestStruc.

# Run Benchmarks

See  https://github.com/ugorji/go-codec-bench/blob/master/codec/bench.sh 
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

Below are results of running the entire suite on 2017-11-20 (ie running ./bench.sh -s).

What you should notice:

- Results get better with codecgen, showing about 20-50% performance improvement.
  Users should carefully weigh the performance improvements against the 
  usability and binary-size increases, as performance is already extremely good 
  without the codecgen path.
  
See [bench.out.txt] for representative result from running `bench.sh` as below, as of 2020-11-16.
```sh
  ./bench.sh -z > bench.out.txt
```

* snippet of benchmark output, running without codecgen (2020-11-16)*
```
BenchmarkCodecXSuite/options-false.../Benchmark__Msgpack____Encode-8         	   16327	     73793 ns/op	    3136 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Binc_______Encode-8         	   15442	     78241 ns/op	    3152 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Simple_____Encode-8         	   15921	     75216 ns/op	    3136 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Cbor_______Encode-8         	   15985	     75456 ns/op	    3136 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Json_______Encode-8         	    7635	    153859 ns/op	    3248 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Std_Json___Encode-8         	    5432	    210771 ns/op	   74090 B/op	     444 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Gob________Encode-8         	    7918	    160088 ns/op	  169704 B/op	     592 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__JsonIter___Encode-8         	    7881	    153445 ns/op	   54483 B/op	     106 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Bson_______Encode-8         	    3996	    306484 ns/op	  238102 B/op	    1095 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Mgobson____Encode-8         	    3168	    382501 ns/op	  292893 B/op	    1721 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__VMsgpack___Encode-8         	    5764	    212583 ns/op	  164228 B/op	     354 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Fxcbor_____Encode-8         	   10000	    102383 ns/op	   49540 B/op	     320 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Sereal_____Encode-8         	    3577	    335741 ns/op	  259062 B/op	    3220 allocs/op

BenchmarkCodecXSuite/options-false.../Benchmark__Msgpack____Decode-8         	    6619	    179925 ns/op	   65685 B/op	     929 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Binc_______Decode-8         	    6454	    183535 ns/op	   65667 B/op	     929 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Simple_____Decode-8         	    6878	    176912 ns/op	   65686 B/op	     929 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Cbor_______Decode-8         	    6502	    190310 ns/op	   65669 B/op	     929 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Json_______Decode-8         	    3646	    326112 ns/op	   89793 B/op	    1088 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Std_Json___Decode-8         	    1532	    790513 ns/op	  130730 B/op	    2961 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Gob________Decode-8         	    4432	    274886 ns/op	  150651 B/op	    2180 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__JsonIter___Decode-8         	    3694	    321904 ns/op	  126302 B/op	    2486 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Bson_______Decode-8         	    2392	    504140 ns/op	  180300 B/op	    4256 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Mgobson____Decode-8         	    2320	    529322 ns/op	  161410 B/op	    6472 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Fxcbor_____Decode-8         	    5084	    235041 ns/op	   67253 B/op	    1299 allocs/op
```

* snippet of bench.out.txt, running with codecgen (2020-11-16) *
```
BenchmarkCodecXGenSuite/options-false.../Benchmark__Msgpack____Encode-8         	   29899	     40000 ns/op	     232 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Binc_______Encode-8         	   27520	     42680 ns/op	     248 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Simple_____Encode-8         	   28378	     42314 ns/op	     232 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Cbor_______Encode-8         	   29206	     41489 ns/op	     232 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Json_______Encode-8         	   10000	    116918 ns/op	     344 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Msgp_______Encode-8         	   44971	     27162 ns/op	       0 B/op	       0 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Easyjson___Encode-8         	   10000	    112439 ns/op	   50526 B/op	      12 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Ffjson_____Encode-8         	    4302	    280621 ns/op	  124716 B/op	    1033 allocs/op

BenchmarkCodecXGenSuite/options-false.../Benchmark__Msgpack____Decode-8         	   10000	    103152 ns/op	   62101 B/op	     871 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Binc_______Decode-8         	   10000	    107485 ns/op	   62087 B/op	     871 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Simple_____Decode-8         	   10000	    102391 ns/op	   62102 B/op	     871 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Cbor_______Decode-8         	   10000	    112386 ns/op	   62103 B/op	     871 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Json_______Decode-8         	    5061	    233754 ns/op	   85043 B/op	    1001 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Msgp_______Decode-8         	   17287	     69502 ns/op	   63894 B/op	     889 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Easyjson___Decode-8         	    5581	    213565 ns/op	   66946 B/op	     459 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Ffjson_____Decode-8         	    3033	    405474 ns/op	   90046 B/op	    1202 allocs/op
```
