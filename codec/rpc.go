// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

/*
RPC

RPC Client and Server Codecs are implemented, so the codecs can be used
with the standard net/rpc package.
*/
package codec

import (
	"io"
	"net/rpc"
)

// GoRpc implements Rpc using the communication protocol defined in net/rpc package.
var GoRpc goRpc

// Rpc interface provides a rpc Server or Client Codec for rpc communication.
type Rpc interface {
	ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec
	ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec
}

type rpcCodec struct {
	rwc io.ReadWriteCloser
	dec *Decoder
	enc *Encoder
}

type goRpcCodec struct {
	rpcCodec
}

// goRpc is the implementation of Rpc that uses the communication protocol
// as defined in net/rpc package.
type goRpc struct{}

func (x goRpc) ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec {
	return goRpcCodec{newRPCCodec(conn, h)}
}

func (x goRpc) ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec {
	return goRpcCodec{newRPCCodec(conn, h)}
}

func newRPCCodec(conn io.ReadWriteCloser, h Handle) rpcCodec {
	return rpcCodec{
		rwc: conn,
		dec: NewDecoder(conn, h),
		enc: NewEncoder(conn, h),
	}
}

// /////////////// RPC Codec Shared Methods ///////////////////
func (c rpcCodec) write(objs ...interface{}) (err error) {
	for _, obj := range objs {
		if err = c.enc.Encode(obj); err != nil {
			return
		}
	}
	return
}

func (c rpcCodec) read(objs ...interface{}) (err error) {
	for _, obj := range objs {
		//If nil is passed in, we should still attempt to read content to nowhere.
		if obj == nil {
			//obj = &obj //This bombs/uses all memory up. Dunno why (maybe because obj is not addressable???).
			var n interface{}
			obj = &n
		}
		if err = c.dec.Decode(obj); err != nil {
			return
		}
	}
	return
}

func (c rpcCodec) Close() error {
	return c.rwc.Close()
}

func (c rpcCodec) ReadResponseBody(body interface{}) (err error) {
	err = c.read(body)
	return
}

func (c rpcCodec) ReadRequestBody(body interface{}) error {
	return c.read(body)
}

// /////////////// Go RPC Codec ///////////////////
func (c goRpcCodec) WriteRequest(r *rpc.Request, body interface{}) error {
	return c.write(r, body)
}

func (c goRpcCodec) WriteResponse(r *rpc.Response, body interface{}) error {
	return c.write(r, body)
}

func (c goRpcCodec) ReadResponseHeader(r *rpc.Response) (err error) {
	err = c.read(r)
	return
}

func (c goRpcCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.read(r)
}
