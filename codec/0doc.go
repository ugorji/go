// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

/*
High Performance, Feature-Rich Idiomatic Go codec/encoding library for binc, msgpack and cbor.

Supported Serialization formats are:

  - msgpack: [https://github.com/msgpack/msgpack]
  - binc:    [http://github.com/ugorji/binc]
  - cbor:    [http://cbor.io] [http://tools.ietf.org/html/rfc7049]
  - simple:  

To install:

    go get github.com/ugorji/go/codec

The idiomatic Go support is as seen in other encoding packages in
the standard library (ie json, xml, gob, etc).

Rich Feature Set includes:

  - Simple but extremely powerful and feature-rich API
  - Very High Performance.
    Our extensive benchmarks show us outperforming Gob, Json and Bson by 2-4X.
    Achieved by extreme care on allocations, recursions, bypassing reflection, zero-copy, etc.
  - Multiple conversions:
    Package co-erces types where appropriate e.g. decode an int in the stream into a float, etc
  - Corner Cases: Overflows, nil maps/slices, nil values in streams are handled correctly
  - Standard field renaming via tags
  - Encoding from any value
    (struct, slice, map, primitives, pointers, interface{}, etc)
  - Decoding into pointer to any value
    (struct, slice, map, int, float32, bool, string, reflect.Value, etc)
  - Supports extension functions to handle the encode/decode of custom types
  - Support Go 1.2 encoding.BinaryMarshaler/BinaryUnmarshaler
  - Schema-less decoding
    (decode into a pointer to a nil interface{} as opposed to a typed value).
    Includes Options to configure what specific map or slice type to use
    when decoding an encoded list or map into a nil interface{}
  - Provides a RPC Server and Client Codec for net/rpc communication protocol.
  - Fast Paths for some container types:
    For some container types, we circumvent reflection and its associated overhead
    and allocation costs, and encode/decode directly. These types are:
	    Slice of all builtin types and interface{},
	    map of all builtin types and interface{} to string, interface{}, int, int64, uint64
	    symetrical maps of all builtin types and interface{}
  - Msgpack Specific:
      - Options to resolve ambiguities in handling raw bytes (as string or []byte)
        during schema-less decoding (decoding into a nil interface{})
      - RPC Server/Client Codec for msgpack-rpc protocol defined at:
        https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md

Extension Support

Users can register a function to handle the encoding or decoding of
their custom types.

There are no restrictions on what the custom type can be. Some examples:

    type BisSet   []int
    type BitSet64 uint64
    type UUID     string
    type MyStructWithUnexportedFields struct { a int; b bool; c []int; }
    type GifImage struct { ... }

As an illustration, MyStructWithUnexportedFields would normally be
encoded as an empty map because it has no exported fields, while UUID
would be encoded as a string. However, with extension support, you can
encode any of these however you like.

RPC

RPC Client and Server Codecs are implemented, so the codecs can be used
with the standard net/rpc package.

Usage

Typical usage model:

    // create and configure Handle
    var (
      bh codec.BincHandle
      mh codec.MsgpackHandle
      ch codec.CborHandle
    )

    mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

    // configure extensions
    // e.g. for msgpack, define functions and enable Time support for tag 1
    // mh.SetExt(reflect.TypeOf(time.Time{}), 1, myExt)

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

Run the benchmark suite using:
   go test -bi -bench=. -benchmem

To run full benchmark suite (including against vmsgpack and bson),
see notes in ext_dep_test.go

*/
package codec

