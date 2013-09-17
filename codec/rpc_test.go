// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Copyright (c) 2013 Max Persson <max@looplab.se>. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"fmt"
	"net"
	"net/rpc"
)

type TestServer struct{}

func (t *TestServer) Test(args *string, reply *string) error {
	*reply = "reply"
	return nil
}

func ExampleRpc() {
	// This defined the codec to use.
	var mh MsgpackHandle

	// Ready channel to wait for the server to start before calling.
	ready := make(chan struct{})

	// Run server in goroutine.
	go func() {
		// The receiver is defined as follows:
		//
		// type TestServer struct{}
		//
		// func (t *TestServer) Test(args *string, reply *string) error {
		//     *reply = "reply"
		//     return nil
		// }

		// Create and register the receiver.
		testServer := new(TestServer)
		rpc.Register(testServer)

		// Listen to TCP connections.
		listener, err := net.Listen("tcp", ":5555")
		if err != nil {
			panic(err)
		}
		defer listener.Close()

		// Signal that we are ready to serve.
		close(ready)

		// Handle one request at a time.
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			rpcCodec := MsgpackSpecRpc.ServerCodec(conn, &mh)
			rpc.ServeCodec(rpcCodec)
		}
	}()

	// Wait for server to be ready.
	<-ready

	// Connect to the server using the defined codec.
	conn, err := net.Dial("tcp", "localhost:5555")
	if err != nil {
		panic(err)
	}
	rpcCodec := MsgpackSpecRpc.ClientCodec(conn, &mh)
	client := rpc.NewClientWithCodec(rpcCodec)

	// Call the server method defined by TestServer.
	var reply string
	err = client.Call("TestServer.Test", "", &reply)
	if err != nil {
		panic(err)
	}

	fmt.Println(reply)
	// Output: reply
}
