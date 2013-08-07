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
	"bufio"
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
	encbuf *bufio.Writer 
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
	encbuf := bufio.NewWriter(conn)
	return rpcCodec{
		rwc: conn,
		encbuf: encbuf,
		enc: NewEncoder(encbuf, h),
		dec: NewDecoder(bufio.NewReader(conn), h),
		//enc: NewEncoder(conn, h),
		//dec: NewDecoder(conn, h),
	}
}

// /////////////// RPC Codec Shared Methods ///////////////////
func (c rpcCodec) write(obj1, obj2 interface{}, writeObj2, doFlush bool) (err error) {
	if err = c.enc.Encode(obj1); err != nil {
		return
	}
	if writeObj2 {
		if err = c.enc.Encode(obj2); err != nil {
			return
		}
	}
	if doFlush && c.encbuf != nil {
		//println("rpc flushing")
		return c.encbuf.Flush()
	}
	return
}


func (c rpcCodec) read(obj interface{}) (err error) {
	//If nil is passed in, we should still attempt to read content to nowhere.
	if obj == nil {
		var obj2 interface{}
		return c.dec.Decode(&obj2)
	}
	return c.dec.Decode(obj)
}

func (c rpcCodec) Close() error {
	return c.rwc.Close()
}

func (c rpcCodec) ReadResponseBody(body interface{}) error {
	return c.read(body)
}

func (c rpcCodec) ReadRequestBody(body interface{}) error {
	return c.read(body)
}

// /////////////// Go RPC Codec ///////////////////
func (c goRpcCodec) WriteRequest(r *rpc.Request, body interface{}) error {
	return c.write(r, body, true, true)
}

func (c goRpcCodec) WriteResponse(r *rpc.Response, body interface{}) error {
	return c.write(r, body, true, true)
}

func (c goRpcCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.read(r)
}

func (c goRpcCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.read(r)
}
