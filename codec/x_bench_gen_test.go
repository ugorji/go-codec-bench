//go:build x && generated

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/mailru/easyjson"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/tinylib/msgp/msgp"
)

/*
 To update all these, use:
 go get -u github.com/tinylib/msgp/msgp github.com/tinylib/msgp \
           github.com/pquerna/ffjson/ffjson github.com/pquerna/ffjson \
           github.com/mailru/easyjson/...

 Known Issues with external libraries:
 - msgp io.R/W support doesn't work. It throws error
 - ffjson: generated code fails to compile as of latest commit checked on 2022-12-26

*/

// benchXGenSkipFFJSON triggers whether we throw a panic based
// on whether ffjson generated code compiles.
//
// MARKER: change if ffjson is updated and compiles successfully
const benchXGenSkipFFJSON = true

func init() {
	testPreInitFns = append(testPreInitFns, benchXGenPreInit)
}

func benchXGenPreInit() {
	benchCheckers = append(benchCheckers,
		benchChecker{"msgp", fnMsgpEncodeFn, fnMsgpDecodeFn},
		benchChecker{"easyjson", fnEasyjsonEncodeFn, fnEasyjsonDecodeFn},
		benchChecker{"ffjson", fnFfjsonEncodeFn, fnFfjsonDecodeFn},
	)
}

func fnEasyjsonEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	ts2, ok := ts.(easyjson.Marshaler)
	if !ok || ts2 == nil {
		return nil, errors.New("easyjson: input is not a easyjson.Marshaler")
	}
	if testUseIO() {
		buf := bytes.NewBuffer(bsIn[:0]) // new(bytes.Buffer)
		_, err := easyjson.MarshalToWriter(ts2, buf)
		return buf.Bytes(), err
	}
	return easyjson.Marshal(ts2)
	// return ts.(json.Marshaler).MarshalJSON()
}

func fnEasyjsonDecodeFn(buf []byte, ts interface{}) error {
	ts2, ok := ts.(easyjson.Unmarshaler)
	if !ok {
		return errors.New("easyjson: input is not a easyjson.Unmarshaler")
	}
	if testUseIO() {
		return easyjson.UnmarshalFromReader(bytes.NewReader(buf), ts2)
	}
	return easyjson.Unmarshal(buf, ts2)
	// return ts.(json.Unmarshaler).UnmarshalJSON(buf)
}

func fnFfjsonEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	if benchXGenSkipFFJSON {
		panic(errors.New("ffjson: generated code fails to compile; checked 2022-12-26"))
	}

	return ffjson.Marshal(ts)
	// return ts.(json.Marshaler).MarshalJSON()
}

func fnFfjsonDecodeFn(buf []byte, ts interface{}) error {
	if benchXGenSkipFFJSON {
		panic(errors.New("ffjson: generated code fails to compile; checked 2022-12-26"))
	}

	return ffjson.Unmarshal(buf, ts)
	// return ts.(json.Unmarshaler).UnmarshalJSON(buf)
}

func fnMsgpEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	if _, ok := ts.(msgp.Encodable); !ok {
		return nil, fmt.Errorf("msgp: input of type %T is not a msgp.Encodable", ts)
	}
	if testUseIO() {
		buf := fnBenchmarkByteBuf(bsIn)
		err := ts.(msgp.Encodable).EncodeMsg(msgp.NewWriter(buf))
		return buf.Bytes(), err
	}
	return ts.(msgp.Marshaler).MarshalMsg(bsIn[:0]) // msgp appends to slice.
}

func fnMsgpDecodeFn(buf []byte, ts interface{}) (err error) {
	if _, ok := ts.(msgp.Decodable); !ok {
		return fmt.Errorf("msgp: input of type %T is not a msgp.Decodable", ts)
	}
	if testUseIO() {
		err = ts.(msgp.Decodable).DecodeMsg(msgp.NewReader(bytes.NewReader(buf)))
		return
	}
	_, err = ts.(msgp.Unmarshaler).UnmarshalMsg(buf)
	return
}

func Benchmark__Msgp_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "msgp", benchTs, fnMsgpEncodeFn)
}

func Benchmark__Msgp_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "msgp", benchTs, fnMsgpEncodeFn, fnMsgpDecodeFn, fnBenchNewTs)
}

func Benchmark__Easyjson___Encode(b *testing.B) {
	fnBenchmarkEncode(b, "easyjson", benchTs, fnEasyjsonEncodeFn)
}

func Benchmark__Easyjson___Decode(b *testing.B) {
	fnBenchmarkDecode(b, "easyjson", benchTs, fnEasyjsonEncodeFn, fnEasyjsonDecodeFn, fnBenchNewTs)
}

func Benchmark__Ffjson_____Encode(b *testing.B) {
	fnBenchmarkEncode(b, "ffjson", benchTs, fnFfjsonEncodeFn)
}

func Benchmark__Ffjson_____Decode(b *testing.B) {
	fnBenchmarkDecode(b, "ffjson", benchTs, fnFfjsonEncodeFn, fnFfjsonDecodeFn, fnBenchNewTs)
}
