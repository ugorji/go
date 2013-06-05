# Codec

High Performance and Feature-Rich Idiomatic Go Library providing
encode/decode support for different serialization formats.

Supported Serialization formats are:

  - msgpack: [http://wiki.msgpack.org/display/MSGPACK/Format+specification]
  - binc: [http://github.com/ugorji/binc]

To install:

    go get github.com/ugorji/go/codec

Online documentation: [http://godoc.org/github.com/ugorji/go/codec]

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

## Extension Support

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

## Usage

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
    rpcCodec := rpcH.ClientCodec(conn, h)  
    client := rpc.NewClientWithCodec(rpcCodec)

## Representative Benchmark Results

A sample run of benchmark using "go test -bi -bench=.":

    ..............................................
    Benchmark: 
    	Struct recursive Depth:             1
    	ApproxDeepSize Of benchmark Struct: 4786
    Benchmark One-Pass Run:
    	   msgpack: len: 1564
    	      binc: len: 1191
    	       gob: len: 1972
    	      json: len: 2538
    	 v-msgpack: len: 1600
    	      bson: len: 3025
    ..............................................
    PASS
    Benchmark__Msgpack__Encode	   50000	     61731 ns/op
    Benchmark__Msgpack__Decode	   10000	    115947 ns/op
    Benchmark__Binc_____Encode	   50000	     64568 ns/op
    Benchmark__Binc_____Decode	   10000	    113843 ns/op
    Benchmark__Gob______Encode	   10000	    143956 ns/op
    Benchmark__Gob______Decode	    5000	    431889 ns/op
    Benchmark__Json_____Encode	   10000	    158662 ns/op
    Benchmark__Json_____Decode	    5000	    310744 ns/op
    Benchmark__Bson_____Encode	   10000	    172905 ns/op
    Benchmark__Bson_____Decode	   10000	    228564 ns/op
    Benchmark__VMsgpack_Encode	   20000	     81752 ns/op
    Benchmark__VMsgpack_Decode	   10000	    160050 ns/op

To run full benchmark suite (including against vmsgpack and bson), 
see notes in ext\_dep\_test.go

