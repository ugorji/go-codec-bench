// +build alltests
// +build go1.7

package codec

// see notes in z_all_test.go

import "testing"

import . "github.com/ugorji/go/codec"

func benchmarkSuite(t *testing.B, f func(t *testing.B)) {
	// find . -name "*_test.go" | xargs grep -e 'flag.' | cut -d '&' -f 2 | cut -d ',' -f 1 | grep -e '^bench'

	testReinit() // so flag.Parse() is called first, and never called again
	benchReinit()

	testDecodeOptions = DecodeOptions{}
	testEncodeOptions = EncodeOptions{}

	testUseIoEncDec = false
	testUseReset = false
	testInternStr = false
	testJsonPreferFloat = false

	benchMapStringKeyOnly = false
	benchInitDebug = false
	benchVerify = false
	benchDepth = 2
	benchDoInitBench = true
	benchUnscientificRes = true
	testReinit()
	benchReinit()
	t.Run("init-metrics", f)

	benchDoInitBench = false
	benchUnscientificRes = false

	testReinit()
	benchReinit()
	t.Run("options-false", f)

	testUseIoEncDec = true
	testDecodeOptions.ReaderBufferSize = 128
	testEncodeOptions.WriterBufferSize = 128
	testReinit()
	benchReinit()
	t.Run("use-io-not-bytes", f)
	testDecodeOptions.ReaderBufferSize = 0
	testEncodeOptions.WriterBufferSize = 0
	testUseIoEncDec = false

	testUseReset = true
	testReinit()
	benchReinit()
	t.Run("reset-enc-dec", f)
	testUseReset = false

	testInternStr = true
	testReinit()
	benchReinit()
	t.Run("intern-strings", f)
	testInternStr = false

	// testJsonPreferFloat = true
	benchVerify = true
	testReinit()
	benchReinit()
	t.Run("verify-on-decode", f)
	benchVerify = false
	// testJsonPreferFloat = false
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
	logT(nil, "\n-------------------------------\n")
	t.Run("Benchmark__Msgpack____Encode", Benchmark__Msgpack____Encode)
	t.Run("Benchmark__Binc_______Encode", Benchmark__Binc_______Encode)
	t.Run("Benchmark__Simple_____Encode", Benchmark__Simple_____Encode)
	t.Run("Benchmark__Cbor_______Encode", Benchmark__Cbor_______Encode)
	t.Run("Benchmark__Json_______Encode", Benchmark__Json_______Encode)
	t.Run("Benchmark__Std_Json___Encode", Benchmark__Std_Json___Encode)
	t.Run("Benchmark__Gob________Encode", Benchmark__Gob________Encode)

	t.Run("Benchmark__Msgpack____Decode", Benchmark__Msgpack____Decode)
	t.Run("Benchmark__Binc_______Decode", Benchmark__Binc_______Decode)
	t.Run("Benchmark__Simple_____Decode", Benchmark__Simple_____Decode)
	t.Run("Benchmark__Cbor_______Decode", Benchmark__Cbor_______Decode)
	t.Run("Benchmark__Json_______Decode", Benchmark__Json_______Decode)
	t.Run("Benchmark__Std_Json___Decode", Benchmark__Std_Json___Decode)
	t.Run("Benchmark__Gob________Decode", Benchmark__Gob________Decode)
}

func BenchmarkCodecSuite(t *testing.B) { benchmarkSuite(t, benchmarkCodecGroup) }
