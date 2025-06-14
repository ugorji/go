// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"testing"
)

func TestGH417(t *testing.T) {
	const numMapEntries = 1024 // 1000
	const numGoroutines = 64   // 100

	type S string
	type E struct {
		Map1 map[string]S      `codec:"map"`
		Map2 map[string]string `codec:"map"` // not error happen, when using builtin type
	}

	e := E{
		Map1: map[string]S{},
		Map2: map[string]string{},
	}

	for i := 0; i < numMapEntries; i++ {
		key1 := make([]byte, 10)
		rand.Read(key1)

		b1 := make([]byte, 10)
		rand.Read(b1)
		e.Map1[hex.EncodeToString(key1)] = S(hex.EncodeToString(b1))
		e.Map2[hex.EncodeToString(key1)] = hex.EncodeToString(b1)
	}

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var b bytes.Buffer
			NewEncoder(&b, testSimpleH).MustEncode(e)
			// if err := NewEncoder(&b, &msgpackHandle).Encode(e); err != nil {
			// 	fmt.Printf("error: %v\n", err)
			// }
		}()
	}

	wg.Wait()
}
