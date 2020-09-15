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
  
See https://github.com/ugorji/go-codec-bench/blob/master/bench.out.txt for representative result from running `bench.sh` as below, as of 2020-09-11.
```sh
  printf "**** STATS ****\n\n"
  ./bench.sh -tx  # (stats)
  printf "**** SUITE ****\n\n"
  ./bench.sh -sx  # (not codecgen)
  printf "**** SUITE (WITH CODECGEN) ****\n\n"
  ./bench.sh -sgx # (codecgen)
```

* snippet of benchmark output, running without codecgen (2020-09-11)*
```
BenchmarkCodecXSuite/options-false.../Benchmark__Msgpack____Encode-8         	   14905	     77946 ns/op	    3192 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Binc_______Encode-8         	   14835	     80659 ns/op	    3192 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Simple_____Encode-8         	   15114	     80583 ns/op	    3192 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Cbor_______Encode-8         	   15292	     79259 ns/op	    3192 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Json_______Encode-8         	    7834	    157497 ns/op	    3256 B/op	      44 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Std_Json___Encode-8         	    5710	    205372 ns/op	   74473 B/op	     444 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Gob________Encode-8         	    7972	    153002 ns/op	  170413 B/op	     591 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__JsonIter___Encode-8         	    7356	    162228 ns/op	   57063 B/op	     142 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Bson_______Encode-8         	    5158	    237240 ns/op	  222827 B/op	     364 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Mgobson____Encode-8         	    3129	    382564 ns/op	  300578 B/op	    1721 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__VMsgpack___Encode-8         	    5796	    209100 ns/op	  164483 B/op	     354 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Fxcbor_____Encode-8         	   10000	    105667 ns/op	   49572 B/op	     320 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Sereal_____Encode-8         	    3667	    327725 ns/op	  274359 B/op	    3218 allocs/op

BenchmarkCodecXSuite/options-false.../Benchmark__Msgpack____Decode-8         	    6062	    199709 ns/op	   67372 B/op	     913 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Binc_______Decode-8         	    5906	    195737 ns/op	   67375 B/op	     913 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Simple_____Decode-8         	    6327	    191206 ns/op	   67387 B/op	     913 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Cbor_______Decode-8         	    6082	    197908 ns/op	   67372 B/op	     913 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Json_______Decode-8         	    3531	    345377 ns/op	   90930 B/op	    1049 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Std_Json___Decode-8         	    1486	    788927 ns/op	  138561 B/op	    3032 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Gob________Decode-8         	    4482	    268897 ns/op	  156145 B/op	    2242 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__JsonIter___Decode-8         	    3741	    317500 ns/op	  129234 B/op	    2504 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Bson_______Decode-8         	    2574	    467859 ns/op	  183828 B/op	    4085 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Mgobson____Decode-8         	    2353	    510441 ns/op	  165787 B/op	    6472 allocs/op
BenchmarkCodecXSuite/options-false.../Benchmark__Fxcbor_____Decode-8         	    5160	    229965 ns/op	   71069 B/op	    1326 allocs/op
```

* snippet of bench.out.txt, running with codecgen (2020-09-11) *
```
BenchmarkCodecXGenSuite/options-false.../Benchmark__Msgpack____Encode-8         	   28167	     40910 ns/op	     288 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Binc_______Encode-8         	   27782	     43138 ns/op	     288 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Simple_____Encode-8         	   27820	     42704 ns/op	     288 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Cbor_______Encode-8         	   28224	     41826 ns/op	     288 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Json_______Encode-8         	    9912	    119889 ns/op	     352 B/op	       2 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Msgp_______Encode-8         	   42363	     29064 ns/op	       0 B/op	       0 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Easyjson___Encode-8         	    8216	    147945 ns/op	   50498 B/op	      12 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Ffjson_____Encode-8         	    4210	    279269 ns/op	  128124 B/op	    1033 allocs/op

BenchmarkCodecXGenSuite/options-false.../Benchmark__Msgpack____Decode-8         	   10000	    109836 ns/op	   64120 B/op	     871 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Binc_______Decode-8         	   10000	    114463 ns/op	   64087 B/op	     871 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Simple_____Decode-8         	   10000	    108317 ns/op	   64069 B/op	     871 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Cbor_______Decode-8         	   10000	    113531 ns/op	   64087 B/op	     871 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Json_______Decode-8         	    4846	    244543 ns/op	   87473 B/op	    1002 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Msgp_______Decode-8         	   17607	     68284 ns/op	   65862 B/op	     889 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Easyjson___Decode-8         	    4732	    250434 ns/op	   96870 B/op	    1000 allocs/op
BenchmarkCodecXGenSuite/options-false.../Benchmark__Ffjson_____Decode-8         	    2594	    412920 ns/op	   92223 B/op	    1202 allocs/op
```
