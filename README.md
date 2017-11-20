# go-codec-bench

This is a comparison of different binary and text encodings.

We compare the codecs provided by github.com/ugorji/go/codec package,
against other libraries:

[github.com/ugorji/go/codec](http://github.com/ugorji/go) provides:

  - msgpack: [http://github.com/msgpack/msgpack] 
  - binc:    [http://github.com/ugorji/binc]
  - cbor:    [http://cbor.io] [http://tools.ietf.org/html/rfc7049]
  - simple: 
  - json:    [http://json.org] [http://tools.ietf.org/html/rfc7159] 

Other codecs compared include:

  - [gopkg.in/vmihailenco/msgpack.v2](http://gopkg.in/vmihailenco/msgpack.v2)
  - [gopkg.in/mgo.v2/bson](http://gopkg.in/mgo.v2/bson)
  - [github.com/davecgh/go-xdr/xdr2](https://godoc.org/github.com/davecgh/go-xdr/xdr)
  - [github.com/Sereal/Sereal/Go/sereal](https://godoc.org/github.com/Sereal/Sereal/Go/sereal)
  - [code.google.com/p/cbor/go](http://code.google.com/p/cbor/go)
  - [github.com/tinylib/msgp](http://github.com/tinylib/msgp)
  - [github.com/tinylib/msgp](http://godoc.org/github.com/tinylib/msgp)
  - [github.com/pquerna/ffjson/ffjson](http://godoc.org/github.com/pquerna/ffjson/ffjson)
  - [bitbucket.org/bodhisnarkva/cbor/go](http://godoc.org/bitbucket.org/bodhisnarkva/cbor/go)
  - [github.com/json-iterator/go](http://godoc.org/github.com/json-iterator/go)
  - [github.com/mailru/easyjson](http://godoc.org/github.com/mailru/easyjson)
  
# Data

The data being serialized is a `TestStruc` randomly generated values.
See https://github.com/ugorji/go-codec-bench/blob/master/values_test.go for the
definition of the TestStruc.

# Run Benchmarks

See  https://github.com/ugorji/go-codec-bench/blob/master/bench.sh 
for how to download the external libraries which we benchmark against,
generate the files for the types when needed, 
and run the suite of tests.

The 3 suite of benchmarks are

  - CodecSuite
  - XSuite
  - CodecXSuite

```
# download the code and all its dependencies
./bench.sh -d

# code-generate files needed for benchmarks against ffjson, easyjson, msgp, etc
./bench.sh -c

# run the full suite of tests
./bench.sh -s

# Below, see how to just run some specific suite of tests, knowing the right tags and flags ...
# See bench.sh for different iterations

# Run suite of tests in default mode (selectively using unsafe in specific areas)
go test -tags "alltests x" -bench "CodecXSuite" -benchmem 
# Run suite of tests in safe mode (no usage of unsafe)
go test -tags "alltests x safe" -bench "CodecXSuite" -benchmem 
# Run suite of tests in codecgen mode, including all tests which are generated (msgp, ffjson, etc)
go test -tags "alltests x generated" -bench "CodecXGenSuite" -benchmem 

```

# Issues

The following issues are seen currently (11/20/2014):

- _code.google.com/p/cbor/go_ fails on encoding and decoding the test struct
- _github.com/davecgh/go-xdr/xdr2_ fails on encoding and decoding the test struct
- _github.com/Sereal/Sereal/Go/sereal_ fails on decoding the serialized test struct

# Representative Benchmark Results

Please see the [benchmarking blog post for detailed representative results](http://ugorji.net/blog/benchmarking-serialization-in-go).

A snapshot of some results on my 2016 MacBook Pro is below.  
**Note: errors are truncated, and lines re-arranged, for readability**.

Below are results of running the entire suite on 2017-11-20 (ie running ./bench.sh -s).

What you should notice:

- Results get better with codecgen, showing about 20-50% performance improvement.
  Users should carefully weigh the performance improvements against the 
  usability and binary-size increases, as performance is already extremely good 
  without the codecgen path.
  
See  https://github.com/ugorji/go-codec-bench/blob/master/bench.out for latest run of bench.sh as of 2017-11-20

* snippet of bench.out, running without codecgen *
```
BenchmarkCodecXSuite/options-false/Benchmark__Msgpack____Encode-8        	   10000	    174305 ns/op	   10224 B/op	      75 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Binc_______Encode-8        	   10000	    195490 ns/op	   12553 B/op	      80 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Simple_____Encode-8        	   10000	    184454 ns/op	   10224 B/op	      75 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Cbor_______Encode-8        	   10000	    178140 ns/op	   10224 B/op	      75 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Json_______Encode-8        	    3000	    483744 ns/op	   10352 B/op	      75 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Std_Json___Encode-8        	    3000	    532382 ns/op	  256049 B/op	     835 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Gob________Encode-8        	    5000	    281319 ns/op	  333545 B/op	     959 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__JsonIter___Encode-8        	    3000	    485147 ns/op	  183552 B/op	    3262 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Bson_______Encode-8        	    2000	    768477 ns/op	  715539 B/op	    5629 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__VMsgpack___Encode-8        	    2000	    663463 ns/op	  320385 B/op	     542 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Sereal_____Encode-8        	    5000	    380179 ns/op	  297532 B/op	    4286 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Msgpack____Decode-8        	    5000	    374246 ns/op	  120352 B/op	    1210 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Binc_______Decode-8        	    3000	    433275 ns/op	  126144 B/op	    1263 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Simple_____Decode-8        	    5000	    386328 ns/op	  120352 B/op	    1210 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Cbor_______Decode-8        	    5000	    381346 ns/op	  120352 B/op	    1210 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Json_______Decode-8        	    2000	    741081 ns/op	  159288 B/op	    1478 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Std_Json___Decode-8        	    1000	   2234843 ns/op	  276336 B/op	    6959 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Gob________Decode-8        	    5000	    405576 ns/op	  256681 B/op	    3261 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__JsonIter___Decode-8        	    2000	    913076 ns/op	  301457 B/op	    7769 allocs/op
BenchmarkCodecXSuite/options-false/Benchmark__Bson_______Decode-8        	    2000	   1163687 ns/op	  373121 B/op	   15703 allocs/op
```

* snippet of bench.out, running with codecgen *
```
BenchmarkCodecXGenSuite/options-false/Benchmark__Msgpack____Encode-8        	   10000	    122495 ns/op	    6224 B/op	       7 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Binc_______Encode-8        	   10000	    117633 ns/op	    6256 B/op	       7 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Simple_____Encode-8        	   10000	    124100 ns/op	    6224 B/op	       7 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Cbor_______Encode-8        	   10000	    126178 ns/op	    6224 B/op	       7 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Json_______Encode-8        	    3000	    424374 ns/op	    6352 B/op	       7 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Std_Json___Encode-8        	    3000	    832128 ns/op	  256049 B/op	     835 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Gob________Encode-8        	    5000	    305695 ns/op	  333551 B/op	     959 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__JsonIter___Encode-8        	    2000	   2074985 ns/op	  183552 B/op	    3262 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Bson_______Encode-8        	     500	   3573898 ns/op	  715538 B/op	    5629 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__VMsgpack___Encode-8        	     500	   3073024 ns/op	  320385 B/op	     542 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Msgp_______Encode-8        	    5000	    272777 ns/op	       0 B/op	       0 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Easyjson___Encode-8        	    1000	   1944819 ns/op	   92826 B/op	      14 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Ffjson_____Encode-8        	     500	   2781966 ns/op	  221909 B/op	    1569 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Sereal_____Encode-8        	    1000	   1704682 ns/op	  300409 B/op	    4285 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Msgpack____Decode-8        	    2000	   1067568 ns/op	  131912 B/op	    1112 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Binc_______Decode-8        	    1000	   1900440 ns/op	  131944 B/op	    1112 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Simple_____Decode-8        	    5000	    248293 ns/op	  131912 B/op	    1112 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Cbor_______Decode-8        	    5000	    248065 ns/op	  131912 B/op	    1112 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Json_______Decode-8        	    3000	    571795 ns/op	  170824 B/op	    1376 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Std_Json___Decode-8        	    1000	   2276679 ns/op	  276337 B/op	    6959 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Gob________Decode-8        	    5000	    391046 ns/op	  256681 B/op	    3261 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__JsonIter___Decode-8        	    2000	    921069 ns/op	  301489 B/op	    7769 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Bson_______Decode-8        	    2000	   1158740 ns/op	  373121 B/op	   15703 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Msgp_______Decode-8        	   10000	    125038 ns/op	  112688 B/op	    1058 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Easyjson___Decode-8        	    3000	    516648 ns/op	  184176 B/op	    1371 allocs/op
BenchmarkCodecXGenSuite/options-false/Benchmark__Ffjson_____Decode-8        	    2000	    934455 ns/op	  161806 B/op	    1927 allocs/op
```
