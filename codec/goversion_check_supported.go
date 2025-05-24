// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !go1.18

package codec

import "errors"

// From v1.3, this codec package will only work for go1.18 and above, as it now needs generics.

func init() {
	panic(errors.New("codec: go 1.17 and below are not supported (generics needed)"))
}
