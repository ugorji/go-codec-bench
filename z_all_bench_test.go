// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build alltests codecgen
// +build go1.7

package codec

// see notes in z_all_test.go

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

import . "github.com/ugorji/go/codec"

var benchmarkGroupOnce sync.Once

var benchmarkGroupSave struct {
	testUseIoEncDec int
	testUseReset    bool
	testInternStr   bool

	benchDepth            int
	benchMapStringKeyOnly bool
	benchInitDebug        bool
	benchVerify           bool
	benchDoInitBench      bool
	benchUnscientificRes  bool
}

func benchmarkGroupInitAll() {
	testInitAll() // calls flag.Parse
	benchmarkGroupSave.testUseIoEncDec = testUseIoEncDec
	benchmarkGroupSave.testUseReset = testUseReset
	benchmarkGroupSave.testInternStr = testInternStr

	benchmarkGroupSave.benchDepth = benchDepth
	benchmarkGroupSave.benchMapStringKeyOnly = benchMapStringKeyOnly
	benchmarkGroupSave.benchInitDebug = benchInitDebug
	benchmarkGroupSave.benchVerify = benchVerify
	benchmarkGroupSave.benchDoInitBench = benchDoInitBench
	benchmarkGroupSave.benchUnscientificRes = benchUnscientificRes
}

func benchmarkGroupReset() {
	testUseIoEncDec = benchmarkGroupSave.testUseIoEncDec
	testUseReset = benchmarkGroupSave.testUseReset
	testInternStr = benchmarkGroupSave.testInternStr

	benchDepth = benchmarkGroupSave.benchDepth
	benchMapStringKeyOnly = benchmarkGroupSave.benchMapStringKeyOnly
	benchInitDebug = benchmarkGroupSave.benchInitDebug
	benchVerify = benchmarkGroupSave.benchVerify
	benchDoInitBench = benchmarkGroupSave.benchDoInitBench
	benchUnscientificRes = benchmarkGroupSave.benchUnscientificRes
}

func benchmarkOneFn(fns []func(*testing.B)) func(*testing.B) {
	switch len(fns) {
	case 0:
		return nil
	case 1:
		return fns[0]
	default:
		return func(t *testing.B) {
			for _, f := range fns {
				f(t)
			}
		}
	}
}

func benchmarkSuiteNoop(b *testing.B) {
	testOnce.Do(testInitAll)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Sleep(1 * time.Millisecond)
	}
}

func benchmarkSuite(t *testing.B, fns ...func(t *testing.B)) {
	benchmarkGroupOnce.Do(benchmarkGroupInitAll)

	f := benchmarkOneFn(fns)
	// find . -name "*_test.go" | xargs grep -e 'flag.' | cut -d '&' -f 2 | cut -d ',' -f 1 | grep -e '^bench'

	testReinit() // so flag.Parse() is called first, and never called again
	benchReinit()

	testDecodeOptions = DecodeOptions{}
	testEncodeOptions = EncodeOptions{}

	benchmarkGroupReset()

	benchVerify = true
	benchDoInitBench = true
	benchUnscientificRes = true
	testReinit()
	benchReinit()
	t.Run("init-metrics....", func(t *testing.B) { t.Run("Benchmark__Noop.............", benchmarkSuiteNoop) })

	benchVerify = false
	benchDoInitBench = false
	benchUnscientificRes = false

	testReinit()
	benchReinit()
	t.Run("options-false...", f)

	testUseIoEncDec = 128
	testReinit()
	benchReinit()
	t.Run("use-io-not-bytes", f)
	testUseIoEncDec = -1

	testUseReset = true
	testReinit()
	benchReinit()
	t.Run("reset-enc-dec...", f)
	testUseReset = false

	// intern string only applies to binc: don't do a full run of it
	// testInternStr = true
	// testReinit()
	// benchReinit()
	// t.Run("intern-strings", f)
	// testInternStr = false

	// benchVerify is kinda lame - serves no real purpose.
	// benchVerify = true
	// testReinit()
	// benchReinit()
	// t.Run("verify-on-decode", f)
	// benchVerify = false
}

func benchmarkQuickSuite(t *testing.B, fns ...func(t *testing.B)) {
	benchmarkGroupOnce.Do(benchmarkGroupInitAll)
	f := benchmarkOneFn(fns)
	benchmarkGroupReset()

	// bd=1 2 | ti=-1, 1024 |

	testUseReset = true
	testUseIoEncDec = -1
	// benchDepth = depth
	testReinit()
	benchReinit()

	t.Run("json-all-bd"+strconv.Itoa(benchDepth)+"........", f)

	testUseReset = true
	testUseIoEncDec = 0
	// benchDepth = depth
	testReinit()
	benchReinit()
	t.Run("json-all-bd"+strconv.Itoa(benchDepth)+"-io.....", f)

	testUseReset = true
	testUseIoEncDec = 1024
	// benchDepth = depth
	testReinit()
	benchReinit()
	t.Run("json-all-bd"+strconv.Itoa(benchDepth)+"-buf1024", f)

	benchmarkGroupReset()
}

/*
z='bench_test.go'
find . -name "$z" | xargs grep -e '^func Benchmark.*Encode' | \
    cut -d '(' -f 1 | cut -d ' ' -f 2 | \
    while read f; do echo "t.Run(\"$f\", $f)"; done &&
echo &&
find . -name "$z" | xargs grep -e '^func Benchmark.*Decode' | \
    cut -d '(' -f 1 | cut -d ' ' -f 2 | \
    while read f; do echo "t.Run(\"$f\", $f)"; done
*/

func benchmarkCodecGroup(t *testing.B) {
	logT(nil, "-------------------------------\n")
	t.Run("Benchmark__Msgpack____Encode", Benchmark__Msgpack____Encode)
	t.Run("Benchmark__Binc_______Encode", Benchmark__Binc_______Encode)
	t.Run("Benchmark__Simple_____Encode", Benchmark__Simple_____Encode)
	t.Run("Benchmark__Cbor_______Encode", Benchmark__Cbor_______Encode)
	t.Run("Benchmark__Json_______Encode", Benchmark__Json_______Encode)
	t.Run("Benchmark__Std_Json___Encode", Benchmark__Std_Json___Encode)
	t.Run("Benchmark__Gob________Encode", Benchmark__Gob________Encode)
	// t.Run("Benchmark__Std_Xml____Encode", Benchmark__Std_Xml____Encode)
	logT(nil, "-------------------------------\n")
	t.Run("Benchmark__Msgpack____Decode", Benchmark__Msgpack____Decode)
	t.Run("Benchmark__Binc_______Decode", Benchmark__Binc_______Decode)
	t.Run("Benchmark__Simple_____Decode", Benchmark__Simple_____Decode)
	t.Run("Benchmark__Cbor_______Decode", Benchmark__Cbor_______Decode)
	t.Run("Benchmark__Json_______Decode", Benchmark__Json_______Decode)
	t.Run("Benchmark__Std_Json___Decode", Benchmark__Std_Json___Decode)
	t.Run("Benchmark__Gob________Decode", Benchmark__Gob________Decode)
	// t.Run("Benchmark__Std_Xml____Decode", Benchmark__Std_Xml____Decode)
}

func BenchmarkCodecSuite(t *testing.B) { benchmarkSuite(t, benchmarkCodecGroup) }

func benchmarkJsonEncodeGroup(t *testing.B) {
	t.Run("Benchmark__Json_______Encode", Benchmark__Json_______Encode)
}

func benchmarkJsonDecodeGroup(t *testing.B) {
	t.Run("Benchmark__Json_______Decode", Benchmark__Json_______Decode)
}

func BenchmarkCodecQuickJsonSuite(t *testing.B) {
	benchmarkQuickSuite(t, benchmarkJsonEncodeGroup)
	benchmarkQuickSuite(t, benchmarkJsonDecodeGroup)

	// depths := [...]int{1, 4}
	// for _, d := range depths {
	// 	benchmarkQuickSuite(t, d, benchmarkJsonEncodeGroup)
	// 	benchmarkQuickSuite(t, d, benchmarkJsonDecodeGroup)
	// }

	// benchmarkQuickSuite(t, 1, benchmarkJsonEncodeGroup)
	// benchmarkQuickSuite(t, 4, benchmarkJsonEncodeGroup)
	// benchmarkQuickSuite(t, 1, benchmarkJsonDecodeGroup)
	// benchmarkQuickSuite(t, 4, benchmarkJsonDecodeGroup)

	// benchmarkQuickSuite(t, 1, benchmarkJsonEncodeGroup, benchmarkJsonDecodeGroup)
	// benchmarkQuickSuite(t, 4, benchmarkJsonEncodeGroup, benchmarkJsonDecodeGroup)
	// benchmarkQuickSuite(t, benchmarkJsonEncodeGroup)
	// benchmarkQuickSuite(t, benchmarkJsonDecodeGroup)
}
