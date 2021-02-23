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
go test -tags "alltests x codec.safe" -bench "CodecXSuite" -benchmem 
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

A snapshot of some results on my Core-i5 2018 Dell Inspiron laptop is below.  
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

*snippet of benchmark output, running without codecgen (2021-02-04)*  
*note that the first 5 are from codes (and the following are from other libraries (as named).*
```
BenchmarkCodecXSuite/options-false.../
Benchmark__Msgpack____Encode-8         	   17764	     67504 ns/op	      24 B/op	       1 allocs/op
Benchmark__Binc_______Encode-8         	   16442	     72011 ns/op	      24 B/op	       1 allocs/op
Benchmark__Simple_____Encode-8         	   17060	     69965 ns/op	      24 B/op	       1 allocs/op
Benchmark__Cbor_______Encode-8         	   16450	     70783 ns/op	      24 B/op	       1 allocs/op
Benchmark__Json_______Encode-8         	    7303	    145232 ns/op	      24 B/op	       1 allocs/op
Benchmark__Std_Json___Encode-8         	    5308	    229397 ns/op	   79216 B/op	     556 allocs/op
Benchmark__Gob________Encode-8         	    7278	    165572 ns/op	  174169 B/op	     683 allocs/op
Benchmark__JsonIter___Encode-8         	    7950	    154060 ns/op	   54437 B/op	      80 allocs/op
Benchmark__Bson_______Encode-8         	    3780	    305520 ns/op	  242263 B/op	    1181 allocs/op
Benchmark__Mgobson____Encode-8         	    3040	    397415 ns/op	  300345 B/op	    1877 allocs/op
Benchmark__VMsgpack___Encode-8         	    5572	    218910 ns/op	  164607 B/op	     360 allocs/op
Benchmark__Fxcbor_____Encode-8         	   10000	    114321 ns/op	   51947 B/op	     400 allocs/op
Benchmark__Sereal_____Encode-8         	    3063	    350328 ns/op	  277718 B/op	    3388 allocs/op

Benchmark__Msgpack____Decode-8         	    8372	    144322 ns/op	   37412 B/op	     315 allocs/op
Benchmark__Binc_______Decode-8         	    8487	    141145 ns/op	   37413 B/op	     315 allocs/op
Benchmark__Simple_____Decode-8         	    8542	    139446 ns/op	   37413 B/op	     315 allocs/op
Benchmark__Cbor_______Decode-8         	    7989	    154240 ns/op	   37414 B/op	     315 allocs/op
Benchmark__Json_______Decode-8         	    2880	    416769 ns/op	   71880 B/op	     592 allocs/op
Benchmark__Std_Json___Decode-8         	    1207	    988487 ns/op	  137146 B/op	    3113 allocs/op
Benchmark__Gob________Decode-8         	    4083	    300238 ns/op	  159394 B/op	    2329 allocs/op
Benchmark__JsonIter___Decode-8         	    2607	    476016 ns/op	  132690 B/op	    2654 allocs/op
Benchmark__Bson_______Decode-8         	    2418	    508113 ns/op	  190987 B/op	    4544 allocs/op
Benchmark__Mgobson____Decode-8         	    2217	    543875 ns/op	  168741 B/op	    6776 allocs/op
Benchmark__VMsgpack___Decode-8         	    3133	    391164 ns/op	  100126 B/op	    2030 allocs/op
Benchmark__Fxcbor_____Decode-8         	    4645	    251379 ns/op	   72515 B/op	    1395 allocs/op
```

* snippet of bench.out.txt, running with codecgen (2021-02-04) *
```
BenchmarkCodecXGenSuite/options-false.../
Benchmark__Msgpack____Encode-8         	   26810	     45254 ns/op	      24 B/op	       1 allocs/op
Benchmark__Binc_______Encode-8         	   25396	     48931 ns/op	      24 B/op	       1 allocs/op
Benchmark__Simple_____Encode-8         	   28516	     46036 ns/op	      24 B/op	       1 allocs/op
Benchmark__Cbor_______Encode-8         	   26923	     44930 ns/op	      24 B/op	       1 allocs/op
Benchmark__Json_______Encode-8         	    9980	    124222 ns/op	      24 B/op	       1 allocs/op
Benchmark__Msgp_______Encode-8         	   36180	     33802 ns/op	       0 B/op	       0 allocs/op
Benchmark__Easyjson___Encode-8         	   10000	    114262 ns/op	   50690 B/op	      12 allocs/op
Benchmark__Ffjson_____Encode-8         	    4051	    299364 ns/op	  128399 B/op	    1121 allocs/op

Benchmark__Msgpack____Decode-8         	   13057	     94631 ns/op	   35909 B/op	     308 allocs/op
Benchmark__Binc_______Decode-8         	   12678	     91177 ns/op	   35909 B/op	     308 allocs/op
Benchmark__Simple_____Decode-8         	   13144	     90958 ns/op	   35910 B/op	     308 allocs/op
Benchmark__Cbor_______Decode-8         	   10000	    102062 ns/op	   35911 B/op	     308 allocs/op
Benchmark__Json_______Decode-8         	    3019	    365442 ns/op	   70469 B/op	     589 allocs/op
Benchmark__Msgp_______Decode-8         	   14877	     79750 ns/op	   68788 B/op	     969 allocs/op
Benchmark__Easyjson___Decode-8         	    3150	    389608 ns/op	   70531 B/op	     475 allocs/op
Benchmark__Ffjson_____Decode-8         	    2758	    435859 ns/op	   95178 B/op	    1290 allocs/op
```
