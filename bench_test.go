// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// Sample way to run:
// go test -bi -bv -bd=1 -benchmem -bench=.

func init() {
	testPreInitFns = append(testPreInitFns, benchPreInit)
	testPostInitFns = append(testPostInitFns, benchPostInit)
}

var (
	benchTs *TestStruc

	approxSize int

	benchCheckers []benchChecker
)

type benchEncFn func(interface{}, []byte) ([]byte, error)
type benchDecFn func([]byte, interface{}) error
type benchIntfFn func() interface{}

type benchChecker struct {
	name     string
	encodefn benchEncFn
	decodefn benchDecFn
}

func benchPreInit() {
	benchTs = newTestStruc(benchDepth, true, !testSkipIntf, benchMapStringKeyOnly)
	approxSize = approxDataSize(reflect.ValueOf(benchTs)) * 3 / 2 // multiply by 1.5 to appease msgp, and prevent alloc
	// bytesLen := 1024 * 4 * (benchDepth + 1) * (benchDepth + 1)
	// if bytesLen < approxSize {
	// 	bytesLen = approxSize
	// }

	benchCheckers = append(benchCheckers,
		// benchChecker{"noop", fnNoopEncodeFn, fnNoopDecodeFn}, // TODO: why comment this out?
		benchChecker{"msgpack", fnMsgpackEncodeFn, fnMsgpackDecodeFn},
		benchChecker{"binc", fnBincEncodeFn, fnBincDecodeFn},
		benchChecker{"simple", fnSimpleEncodeFn, fnSimpleDecodeFn},
		benchChecker{"cbor", fnCborEncodeFn, fnCborDecodeFn},
		benchChecker{"json", fnJsonEncodeFn, fnJsonDecodeFn},
		benchChecker{"std-json", fnStdJsonEncodeFn, fnStdJsonDecodeFn},
		benchChecker{"gob", fnGobEncodeFn, fnGobDecodeFn},
	)
}

func benchPostInit() {
	if benchDoInitBench {
		runBenchInit()
	}
}

func runBenchInit() {
	logT(nil, "..............................................")
	logT(nil, "BENCHMARK INIT: %v", time.Now())
	logT(nil, "To run full benchmark comparing encodings, "+
		"use: \"go test -bench=.\"")
	logT(nil, "Benchmark: ")
	logT(nil, "\tStruct recursive Depth:             %d", benchDepth)
	if approxSize > 0 {
		logT(nil, "\tApproxDeepSize Of benchmark Struct: %d bytes", approxSize)
	}
	if benchUnscientificRes {
		logT(nil, "Benchmark One-Pass Run (with Unscientific Encode/Decode times): ")
	} else {
		logT(nil, "Benchmark One-Pass Run:")
	}
	for _, bc := range benchCheckers {
		doBenchCheck(bc.name, bc.encodefn, bc.decodefn)
	}
	logT(nil, "..............................................")
	if benchInitDebug {
		logT(nil, "<<<<====>>>> depth: %v, ts: %#v\n", benchDepth, benchTs)
	}
}

var vBenchTs = TestStruc{}

func fnBenchNewTs() interface{} {
	vBenchTs = TestStruc{}
	return &vBenchTs
	// return new(TestStruc)
}

func doBenchCheck(name string, encfn benchEncFn, decfn benchDecFn) {
	runtime.GC()
	tnow := time.Now()
	buf, err := encfn(benchTs, nil)
	if err != nil {
		logT(nil, "\t%10s: **** Error encoding benchTs: %v", name, err)
	}
	encDur := time.Now().Sub(tnow)
	encLen := len(buf)
	runtime.GC()
	if !benchUnscientificRes {
		logT(nil, "\t%10s: len: %d bytes\n", name, encLen)
		return
	}
	tnow = time.Now()
	if err = decfn(buf, new(TestStruc)); err != nil {
		logT(nil, "\t%10s: **** Error decoding into new TestStruc: %v", name, err)
	}
	decDur := time.Now().Sub(tnow)
	logT(nil, "\t%10s: len: %d bytes, encode: %v, decode: %v\n", name, encLen, encDur, decDur)
}

func fnBenchmarkEncode(b *testing.B, encName string, ts interface{}, encfn benchEncFn) {
	testOnce.Do(testInitAll)
	var err error
	bs := make([]byte, 0, approxSize)
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = encfn(ts, bs); err != nil {
			break
		}
	}
	if err != nil {
		logT(b, "Error encoding benchTs: %s: %v", encName, err)
		b.FailNow()
	}
}

func fnBenchmarkDecode(b *testing.B, encName string, ts interface{},
	encfn benchEncFn, decfn benchDecFn, newfn benchIntfFn,
) {
	testOnce.Do(testInitAll)
	bs := make([]byte, 0, approxSize)
	buf, err := encfn(ts, bs)
	if err != nil {
		logT(b, "Error encoding benchTs: %s: %v", encName, err)
		b.FailNow()
	}
	if benchVerify {
		// ts2 := newfn()
		ts1 := ts.(*TestStruc)
		ts2 := new(TestStruc)
		if err = decfn(buf, ts2); err != nil {
			logT(b, "BenchVerify: Error decoding benchTs: %s: %v", encName, err)
			b.FailNow()
		}
		if err = deepEqual(ts1, ts2); err != nil {
			logT(b, "BenchVerify: Error comparing benchTs: %s: %v", encName, err)
			b.FailNow()
		}
	}
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ts = newfn()
		if err = decfn(buf, ts); err != nil {
			break
		}
	}
	if err != nil {
		logT(b, "Error decoding into new TestStruc: %s: %v", encName, err)
		b.FailNow()
	}
}

// ------------ tests below

func fnMsgpackEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return benchFnCodecEncode(ts, bsIn, testMsgpackH)
}

func fnMsgpackDecodeFn(buf []byte, ts interface{}) error {
	return benchFnCodecDecode(buf, ts, testMsgpackH)
}

func fnBincEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return benchFnCodecEncode(ts, bsIn, testBincH)
}

func fnBincDecodeFn(buf []byte, ts interface{}) error {
	return benchFnCodecDecode(buf, ts, testBincH)
}

func fnSimpleEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return benchFnCodecEncode(ts, bsIn, testSimpleH)
}

func fnSimpleDecodeFn(buf []byte, ts interface{}) error {
	return benchFnCodecDecode(buf, ts, testSimpleH)
}

func fnNoopEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return benchFnCodecEncode(ts, bsIn, testNoopH)
}

func fnNoopDecodeFn(buf []byte, ts interface{}) error {
	return benchFnCodecDecode(buf, ts, testNoopH)
}

func fnCborEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return benchFnCodecEncode(ts, bsIn, testCborH)
}

func fnCborDecodeFn(buf []byte, ts interface{}) error {
	return benchFnCodecDecode(buf, ts, testCborH)
}

func fnJsonEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return benchFnCodecEncode(ts, bsIn, testJsonH)
}

func fnJsonDecodeFn(buf []byte, ts interface{}) error {
	return benchFnCodecDecode(buf, ts, testJsonH)
}

func fnGobEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	buf := fnBenchmarkByteBuf(bsIn)
	err := gob.NewEncoder(buf).Encode(ts)
	return buf.Bytes(), err
}

func fnGobDecodeFn(buf []byte, ts interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(buf)).Decode(ts)
}

func fnStdJsonEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	return json.Marshal(ts)
}

func fnStdJsonDecodeFn(buf []byte, ts interface{}) error {
	return json.Unmarshal(buf, ts)
}

// ----------- DECODE ------------------

// Re-enable NoopHandle tests when fixed. TODO: Oct 16, 2015
func __Benchmark__Noop_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "noop", benchTs, fnNoopEncodeFn)
}

func Benchmark__Msgpack____Encode(b *testing.B) {
	fnBenchmarkEncode(b, "msgpack", benchTs, fnMsgpackEncodeFn)
}

func Benchmark__Binc_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "binc", benchTs, fnBincEncodeFn)
}

func Benchmark__Simple_____Encode(b *testing.B) {
	fnBenchmarkEncode(b, "simple", benchTs, fnSimpleEncodeFn)
}

func Benchmark__Cbor_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "cbor", benchTs, fnCborEncodeFn)
}

func Benchmark__Json_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "json", benchTs, fnJsonEncodeFn)
}

func Benchmark__Std_Json___Encode(b *testing.B) {
	fnBenchmarkEncode(b, "std-json", benchTs, fnStdJsonEncodeFn)
}

func Benchmark__Gob________Encode(b *testing.B) {
	fnBenchmarkEncode(b, "gob", benchTs, fnGobEncodeFn)
}

// ----------- DECODE ------------------

func __Benchmark__Noop_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "noop", benchTs, fnNoopEncodeFn, fnNoopDecodeFn, fnBenchNewTs)
}

func Benchmark__Msgpack____Decode(b *testing.B) {
	fnBenchmarkDecode(b, "msgpack", benchTs, fnMsgpackEncodeFn, fnMsgpackDecodeFn, fnBenchNewTs)
}

func Benchmark__Binc_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "binc", benchTs, fnBincEncodeFn, fnBincDecodeFn, fnBenchNewTs)
}

func Benchmark__Simple_____Decode(b *testing.B) {
	fnBenchmarkDecode(b, "simple", benchTs, fnSimpleEncodeFn, fnSimpleDecodeFn, fnBenchNewTs)
}

func Benchmark__Cbor_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "cbor", benchTs, fnCborEncodeFn, fnCborDecodeFn, fnBenchNewTs)
}

func Benchmark__Json_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "json", benchTs, fnJsonEncodeFn, fnJsonDecodeFn, fnBenchNewTs)
}

func Benchmark__Std_Json___Decode(b *testing.B) {
	fnBenchmarkDecode(b, "std-json", benchTs, fnStdJsonEncodeFn, fnStdJsonDecodeFn, fnBenchNewTs)
}

func Benchmark__Gob________Decode(b *testing.B) {
	fnBenchmarkDecode(b, "gob", benchTs, fnGobEncodeFn, fnGobDecodeFn, fnBenchNewTs)
}

// ---------- NOOP -----------

func TestBenchNoop(t *testing.T) {
	testOnce.Do(testInitAll)
}
