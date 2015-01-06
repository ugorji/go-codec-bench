//+build x

// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"bytes"
	"testing"

	gcbor "code.google.com/p/cbor/go"
	"github.com/Sereal/Sereal/Go/sereal"
	"github.com/davecgh/go-xdr/xdr2"
	"github.com/philhofer/msgp/msgp"
	vmsgpack "github.com/vmihailenco/msgpack"
	"gopkg.in/mgo.v2/bson" //"labix.org/v2/mgo/bson"
)

func init() {
	testPreInitFns = append(testPreInitFns, benchXPreInit)
}

func benchXPreInit() {
	benchCheckers = append(benchCheckers,
		benchChecker{"v-msgpack", fnVMsgpackEncodeFn, fnVMsgpackDecodeFn},
		benchChecker{"bson", fnBsonEncodeFn, fnBsonDecodeFn},
		benchChecker{"msgp", fnMsgpEncodeFn, fnMsgpDecodeFn},
		// place codecs with issues at the end, so as not to make results too ugly
		benchChecker{"gcbor", fnGcborEncodeFn, fnGcborDecodeFn},
		benchChecker{"xdr", fnXdrEncodeFn, fnXdrDecodeFn},
		benchChecker{"sereal", fnSerealEncodeFn, fnSerealDecodeFn},
	)
}

func fnVMsgpackEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	return vmsgpack.Marshal(ts)
}

func fnVMsgpackDecodeFn(buf []byte, ts interface{}) error {
	return vmsgpack.Unmarshal(buf, ts)
}

func fnBsonEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	return bson.Marshal(ts)
}

func fnBsonDecodeFn(buf []byte, ts interface{}) error {
	return bson.Unmarshal(buf, ts)
}

func fnXdrEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	buf := fnBenchmarkByteBuf(bsIn)
	_, err := xdr.Marshal(buf, ts)
	return buf.Bytes(), err
}

func fnXdrDecodeFn(buf []byte, ts interface{}) error {
	_, err := xdr.Unmarshal(bytes.NewReader(buf), ts)
	return err
}

func fnSerealEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	return sereal.Marshal(ts)
}

func fnSerealDecodeFn(buf []byte, ts interface{}) error {
	return sereal.Unmarshal(buf, ts)
}

func fnMsgpEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	if benchUseIO {
		buf := fnBenchmarkByteBuf(bsIn)
		err := ts.(msgp.Encodable).EncodeMsg(msgp.NewWriter(buf))
		return buf.Bytes(), err
	}
	return ts.(msgp.Marshaler).MarshalMsg(bsIn[:0]) // msgp appends to slice.
}

func fnMsgpDecodeFn(buf []byte, ts interface{}) (err error) {
	if benchUseIO {
		err = ts.(msgp.Decodable).DecodeMsg(msgp.NewReader(bytes.NewReader(buf)))
	} else {
		_, err = ts.(msgp.Unmarshaler).UnmarshalMsg(buf)
	}
	return
}

func fnGcborEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	buf := fnBenchmarkByteBuf(bsIn)
	err = gcbor.NewEncoder(buf).Encode(ts)
	return buf.Bytes(), err
}

func fnGcborDecodeFn(buf []byte, ts interface{}) error {
	return gcbor.NewDecoder(bytes.NewReader(buf)).Decode(ts)
}

func Benchmark__Bson_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "bson", benchTs, fnBsonEncodeFn)
}

func Benchmark__Bson_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "bson", benchTs, fnBsonEncodeFn, fnBsonDecodeFn, fnBenchNewTs)
}

func Benchmark__VMsgpack___Encode(b *testing.B) {
	fnBenchmarkEncode(b, "v-msgpack", benchTs, fnVMsgpackEncodeFn)
}

func Benchmark__VMsgpack___Decode(b *testing.B) {
	fnBenchmarkDecode(b, "v-msgpack", benchTs, fnVMsgpackEncodeFn, fnVMsgpackDecodeFn, fnBenchNewTs)
}

func Benchmark__Msgp_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "msgp", benchTs, fnMsgpEncodeFn)
}

func Benchmark__Msgp_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "msgp", benchTs, fnMsgpEncodeFn, fnMsgpDecodeFn, fnBenchNewTs)
}

// Place codecs with issues at the bottom, so as not to make results look too ugly.

func Benchmark__Gcbor_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "gcbor", benchTs, fnGcborEncodeFn)
}

func Benchmark__Gcbor_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "gcbor", benchTs, fnGcborEncodeFn, fnGcborDecodeFn, fnBenchNewTs)
}

func Benchmark__Xdr________Encode(b *testing.B) {
	fnBenchmarkEncode(b, "xdr", benchTs, fnXdrEncodeFn)
}

func Benchmark__Xdr________Decode(b *testing.B) {
	fnBenchmarkDecode(b, "xdr", benchTs, fnXdrEncodeFn, fnXdrDecodeFn, fnBenchNewTs)
}

func Benchmark__Sereal_____Encode(b *testing.B) {
	fnBenchmarkEncode(b, "sereal", benchTs, fnSerealEncodeFn)
}

func Benchmark__Sereal_____Decode(b *testing.B) {
	fnBenchmarkDecode(b, "sereal", benchTs, fnSerealEncodeFn, fnSerealDecodeFn, fnBenchNewTs)
}
