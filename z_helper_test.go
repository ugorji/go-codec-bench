// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// All non-std package dependencies related to testing live in this file,
// so porting to different environment is easy (just update functions).
//
// Also, this file is called z_helper_test, to give a "hint" to compiler
// that its init() function should be called last. (not guaranteed by spec)

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	xcodec "github.com/ugorji/go/codec"
	// xcodec "ugorji.net/codec"
)

var (
	testLogToT = true
)

func init() {
	testBincHSym.AsSymbols = xcodec.AsSymbolAll
	testBincHNoSym.AsSymbols = xcodec.AsSymbolNone
	benchInitFlags()
	flag.Parse()
	benchInit()
}

var (
	testMsgpackH   = &xcodec.MsgpackHandle{}
	testBincH      = &xcodec.BincHandle{}
	testBincHSym   = &xcodec.BincHandle{}
	testBincHNoSym = &xcodec.BincHandle{}
	testSimpleH    = &xcodec.SimpleHandle{}
	testCborH      = &xcodec.CborHandle{}
	testJsonH      = &xcodec.JsonHandle{}
)

func fnCodecEncode(ts interface{}, h xcodec.Handle) (bs []byte, err error) {
	err = xcodec.NewEncoderBytes(&bs, h).Encode(ts)
	return
}

func fnCodecDecode(buf []byte, ts interface{}, h xcodec.Handle) error {
	return xcodec.NewDecoderBytes(buf, h).Decode(ts)
}

func logT(x interface{}, format string, args ...interface{}) {
	if t, ok := x.(*testing.T); ok && t != nil && testLogToT {
		t.Logf(format, args...)
	} else if b, ok := x.(*testing.B); ok && b != nil && testLogToT {
		b.Logf(format, args...)
	} else {
		if len(format) == 0 || format[len(format)-1] != '\n' {
			format = format + "\n"
		}
		fmt.Printf(format, args...)
	}
}

func approxDataSize(rv reflect.Value) (sum int) {
	switch rk := rv.Kind(); rk {
	case reflect.Invalid:
	case reflect.Ptr, reflect.Interface:
		sum += int(rv.Type().Size())
		sum += approxDataSize(rv.Elem())
	case reflect.Slice:
		sum += int(rv.Type().Size())
		for j := 0; j < rv.Len(); j++ {
			sum += approxDataSize(rv.Index(j))
		}
	case reflect.String:
		sum += int(rv.Type().Size())
		sum += rv.Len()
	case reflect.Map:
		sum += int(rv.Type().Size())
		for _, mk := range rv.MapKeys() {
			sum += approxDataSize(mk)
			sum += approxDataSize(rv.MapIndex(mk))
		}
	case reflect.Struct:
		//struct size already includes the full data size.
		//sum += int(rv.Type().Size())
		for j := 0; j < rv.NumField(); j++ {
			sum += approxDataSize(rv.Field(j))
		}
	default:
		//pure value types
		sum += int(rv.Type().Size())
	}
	return
}
