// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

/*
High Performance, Feature-Rich Idiomatic Go encoding library for msgpack and binc .

Supported Serialization formats are:

  - msgpack: [https://github.com/msgpack/msgpack]
  - binc: [http://github.com/ugorji/binc]

To install:

    go get github.com/ugorji/go/codec

The idiomatic Go support is as seen in other encoding packages in
the standard library (ie json, xml, gob, etc).

Rich Feature Set includes:

  - Simple but extremely powerful and feature-rich API
  - Very High Performance.   
    Our extensive benchmarks show us outperforming Gob, Json and Bson by 2-4X.
    This was achieved by taking extreme care on:
      - managing allocation
      - stack frame size (important due to Go's use of split stacks), 
      - reflection use
      - recursion implications
      - zero-copy mode (encoding/decoding to byte slice without using temp buffers)
  - Correct.  
    Care was taken to precisely handle corner cases like: 
      overflows, nil maps and slices, nil value in stream, etc.
  - Efficient zero-copying into temporary byte buffers  
    when encoding into or decoding from a byte slice.
  - Standard field renaming via tags
  - Encoding from any value  
    (struct, slice, map, primitives, pointers, interface{}, etc)
  - Decoding into pointer to any non-nil typed value  
    (struct, slice, map, int, float32, bool, string, reflect.Value, etc)
  - Supports extension functions to handle the encode/decode of custom types
  - Support Go 1.2 encoding.BinaryMarshaler/BinaryUnmarshaler
  - Schema-less decoding  
    (decode into a pointer to a nil interface{} as opposed to a typed non-nil value).  
    Includes Options to configure what specific map or slice type to use 
    when decoding an encoded list or map into a nil interface{}
  - Provides a RPC Server and Client Codec for net/rpc communication protocol.
  - Msgpack Specific:
      - Provides extension functions to handle spec-defined extensions (binary, timestamp)
      - Options to resolve ambiguities in handling raw bytes (as string or []byte)  
        during schema-less decoding (decoding into a nil interface{})
      - RPC Server/Client Codec for msgpack-rpc protocol defined at: 
        http://wiki.msgpack.org/display/MSGPACK/RPC+specification

Extension Support

Users can register a function to handle the encoding or decoding of
their custom types. 

There are no restrictions on what the custom type can be. Extensions can
be any type: pointers, structs, custom types off arrays/slices, strings,
etc. Some examples:

    type BisSet   []int
    type BitSet64 uint64
    type UUID     string
    type MyStructWithUnexportedFields struct { a int; b bool; c []int; }
    type GifImage struct { ... }

Typically, MyStructWithUnexportedFields is encoded as an empty map because
it has no exported fields, while UUID will be encoded as a string,
etc. However, with extension support, you can encode any of these
however you like.

We provide implementations of these functions where the spec has defined
an inter-operable format. For msgpack, these are Binary and
time.Time. Library users will have to explicitly configure these as seen
in the usage below.

Usage

Typical usage model:

    var (
      mapStrIntfTyp = reflect.TypeOf(map[string]interface{}(nil))
      sliceByteTyp = reflect.TypeOf([]byte(nil))
      timeTyp = reflect.TypeOf(time.Time{})
    )
    
    // create and configure Handle
    var (
      bh codec.BincHandle
      mh codec.MsgpackHandle
    )

    mh.MapType = mapStrIntfTyp
    
    // configure extensions for msgpack, to enable Binary and Time support for tags 0 and 1
    mh.AddExt(sliceByteTyp, 0, mh.BinaryEncodeExt, mh.BinaryDecodeExt)
    mh.AddExt(timeTyp, 1, mh.TimeEncodeExt, mh.TimeDecodeExt)

    // create and use decoder/encoder
    var (
      r io.Reader
      w io.Writer
      b []byte
      h = &bh // or mh to use msgpack
    )
    
    dec = codec.NewDecoder(r, h)
    dec = codec.NewDecoderBytes(b, h)
    err = dec.Decode(&v) 
    
    enc = codec.NewEncoder(w, h)
    enc = codec.NewEncoderBytes(&b, h)
    err = enc.Encode(v)
    
    //RPC Server
    go func() {
        for {
            conn, err := listener.Accept()
            rpcCodec := codec.GoRpc.ServerCodec(conn, h)
            //OR rpcCodec := codec.MsgpackSpecRpc.ServerCodec(conn, h)
            rpc.ServeCodec(rpcCodec)
        }
    }()
    
    //RPC Communication (client side)
    conn, err = net.Dial("tcp", "localhost:5555")  
    rpcCodec := codec.GoRpc.ClientCodec(conn, h)
    //OR rpcCodec := codec.MsgpackSpecRpc.ClientCodec(conn, h)
    client := rpc.NewClientWithCodec(rpcCodec)

Representative Benchmark Results

A sample run of benchmark using "go test -bi -bench=.":

   ..............................................
   BENCHMARK INIT: 2013-10-04 14:36:50.381959842 -0400 EDT
   To run full benchmark comparing encodings (MsgPack, Binc, JSON, GOB, etc), use: "go test -bench=."
   Benchmark: 
      	Struct recursive Depth:             1
      	ApproxDeepSize Of benchmark Struct: 4694 bytes
   Benchmark One-Pass Run:
      	      bson: len: 3025 bytes
      	   msgpack: len: 1560 bytes
      	      binc: len: 1187 bytes
      	       gob: len: 1972 bytes
      	      json: len: 2538 bytes
   ..............................................
   PASS
   Benchmark__Msgpack__Encode	   50000	     61683 ns/op	   15395 B/op	      91 allocs/op
   Benchmark__Msgpack__Decode	   10000	    111090 ns/op	   15591 B/op	     437 allocs/op
   Benchmark__Binc_____Encode	   50000	     73577 ns/op	   17572 B/op	      95 allocs/op
   Benchmark__Binc_____Decode	   10000	    112656 ns/op	   16474 B/op	     318 allocs/op
   Benchmark__Gob______Encode	   10000	    138140 ns/op	   21164 B/op	     237 allocs/op
   Benchmark__Gob______Decode	    5000	    408067 ns/op	   83202 B/op	    1840 allocs/op
   Benchmark__Json_____Encode	   20000	     80442 ns/op	   13861 B/op	     102 allocs/op
   Benchmark__Json_____Decode	   10000	    249169 ns/op	   14169 B/op	     493 allocs/op
   Benchmark__Bson_____Encode	   10000	    121970 ns/op	   27717 B/op	     514 allocs/op
   Benchmark__Bson_____Decode	   10000	    161103 ns/op	   16444 B/op	     788 allocs/op

To run full benchmark suite (including against vmsgpack and bson), 
see notes in ext_dep_test.go

*/
package codec
