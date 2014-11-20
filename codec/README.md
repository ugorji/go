# Codec

High Performance, Feature-Rich Idiomatic Go codec/encoding library for binc, msgpack, cbor, json.

Supported Serialization formats are:

  - msgpack: [https://github.com/msgpack/msgpack]
  - binc:    [http://github.com/ugorji/binc]
  - cbor:    [http://cbor.io] [http://tools.ietf.org/html/rfc7049]
  - simple: 
  - json:    [http://json.org] [http://tools.ietf.org/html/rfc7159] 

To install:

    go get github.com/ugorji/go/codec

Online documentation: [http://godoc.org/github.com/ugorji/go/codec]

The idiomatic Go support is as seen in other encoding packages in
the standard library (ie json, xml, gob, etc).

Rich Feature Set includes:

  - Simple but extremely powerful and feature-rich API
  - Very High Performance.
    Our extensive benchmarks show us outperforming Gob, Json and Bson by 2-4X.
    Achieved by extreme care on allocations, recursions, bypassing reflection, zero-copy, etc.
  - Multiple conversions:
    Package coerces types where appropriate e.g. decode an int in the stream into a float, etc
  - Corner Cases: Overflows, nil maps/slices, nil values in streams are handled correctly
  - Standard field renaming via tags
  - Encoding from any value and decoding into pointer to any value
    (struct, slice, map, primitives, pointers, interface{}, etc)
  - Supports extension functions to handle the encode/decode of custom types
  - Support encoding.(Binary|Text)(M|Unm)arshaler interfaces
  - Schema-less decoding
    (decode into a pointer to a nil interface{} as opposed to a typed value).
    Includes Options to configure what specific map or slice type to use
    when decoding an encoded list or map into a nil interface{}
  - Provides a RPC Server and Client Codec for net/rpc communication protocol.
  - Fast Paths for common maps and slices of built-in types (numbers, string, bool).
    Reflection (and its associated overhead) is bypassed.
  - Handle unique idiosynchracies of codecs e.g. 
    - For messagepack, configure how ambiguities in handling raw bytes are resolved 
    - For messagepack, provide rpc server/client codec to support  msgpack-rpc protocol defined at:
      https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md

## Extension Support

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

## RPC

RPC Client and Server Codecs are implemented, so the codecs can be used
with the standard net/rpc package.

## Usage

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
    // mh.AddExt(reflect.TypeOf(time.Time{}), 1, myMsgpackTimeEncodeExtFn, myMsgpackTimeDecodeExtFn)

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

