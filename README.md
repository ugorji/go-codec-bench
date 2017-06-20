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

  - [github.com/vmihailenco/msgpack](http://github.com/vmihailenco/msgpack)
  - [gopkg.in/mgo.v2/bson](http://gopkg.in/mgo.v2/bson)
  - [github.com/davecgh/go-xdr/xdr](http://github.com/davecgh/go-xdr/xdr)
  - [github.com/Sereal/Sereal/Go/sereal](http://github.com/Sereal/Sereal/Go/sereal)
  - [code.google.com/p/cbor/go](http://code.google.com/p/cbor/go)
  
# Data

The data being serialized is a `TestStruc` randomly generated values.
See https://github.com/ugorji/go-codec-bench/blob/master/values_test.go for the
definition of the TestStruc.

# Run Benchmarks

```
# download the code and all its dependencies 
go get -u -t github.com/ugorji/go-codec-bench
go get -u github.com/ugorji/go/codec/codecgen
go get -u github.com/philhofer/msgp
go get -u github.com/pquerna/ffjson
go get -u github.com/pquerna/ffjson/ffjson
go get -u bitbucket.org/bodhisnarkva/cbor/go

# benchmark with the default settings 
go test -bench=.

# benchmark with a larger struct, using the -bd parameter. and gather pre-info
go test -bench=. -bd=2 -bi

# see all the test parameters, using the -Z parameter (any recognized param will do)
go test -Z
```

Sample test execution, including setup for codecgen and execution:

```sh
# If you want to run the benchmarks against code generated values.
# Then first generate the code generated values from values_test.go named typed.
# we cannot normally read a _test.go file, so temporarily copy it into a readable file.

zmydir=`pwd`
zgobase=${zmydir%%/src/*}
cp values_test.go values_temp.go
msgp -tests=false -o=values_msgp.go -file=values_temp.go
$zgobase/bin/codecgen -rt codecgen -t 'x,codecgen,!unsafe' -o values_codecgen_test.go -d 19780 values_temp.go
$zgobase/bin/codecgen -u -rt codecgen -t 'x,codecgen,unsafe' -o values_codecgen_unsafe_test.go -d 19781 values_temp.go
# remove the temp file
rm -f values_temp.go
# Run the tests, using only runtime introspection support (normal mode)
go test -tm -bi -benchmem '-bench=_.*En' -tags=x
go test -tm -bi -benchmem '-bench=_.*De' -tags=x
# Run the tests using the codegeneration.
# This involves passing the tags which enable the appropriate files to be run.
go test -tm -tf -bi -benchmem '-bench=_.*En' '-tags=x codecgen unsafe'
go test -tm -tf -bi -benchmem '-bench=_.*De' '-tags=x codecgen unsafe'
```

# Issues

The following issues are seen currently (11/20/2014):

- _code.google.com/p/cbor/go_ fails on encoding and decoding the test struct
- _github.com/davecgh/go-xdr/xdr2_ fails on encoding and decoding the test struct
- _github.com/Sereal/Sereal/Go/sereal_ fails on decoding the serialized test struct

# Representative Benchmark Results

Please see the [benchmarking blog post for detailed representative results](http://ugorji.net/blog/benchmarking-serialization-in-go).

A snapshot of some results on my 2012 i7-2630QM CPU @ 2.00GHz running Ubuntu 14.04 x86_64 GNU/Linux:

```
$ go test -bench=. -benchmem -bi -bu
..............................................
prebuild done successfully
..............................................
BENCHMARK INIT: 2014-12-29 11:26:53.539634155 -0500 EST
To run full benchmark comparing encodings, use: "go test -bench=."
Benchmark: 
	Struct recursive Depth:             1
	ApproxDeepSize Of benchmark Struct: 8259 bytes
Benchmark One-Pass Run:
	   msgpack: len: 2086 bytes
	binc-nosym: len: 2102 bytes
	  binc-sym: len: 1733 bytes
	    simple: len: 2402 bytes
	      cbor: len: 2102 bytes
	      json: len: 3126 bytes
	  std-json: len: 3470 bytes
	       gob: len: 2756 bytes
	 v-msgpack: len: 2408 bytes
	      bson: len: 3997 bytes
	      msgp: len: 2456 bytes
	     gcbor: len: 2716 bytes
	       xdr: **** Error encoding benchTs: xdr:encodeInterface: can't encode nil interface
	       xdr: len: 576 bytes
	    sereal: len: 1676 bytes
..............................................

Benchmark__Noop_______Encode	   10000	     66478 ns/op	    9315 B/op	      73 allocs/op
Benchmark__Msgpack____Encode	   10000	     82843 ns/op	    9219 B/op	      70 allocs/op
Benchmark__Binc_NoSym_Encode	   10000	     82529 ns/op	    9379 B/op	      74 allocs/op
Benchmark__Binc_Sym___Encode	   10000	     96508 ns/op	   11498 B/op	      78 allocs/op
Benchmark__Simple_____Encode	   10000	     84089 ns/op	    9219 B/op	      70 allocs/op
Benchmark__Cbor_______Encode	   10000	     81773 ns/op	    9187 B/op	      70 allocs/op
Benchmark__Json_______Encode	   10000	    107768 ns/op	    9267 B/op	      70 allocs/op
Benchmark__Std_Json___Encode	    5000	    123896 ns/op	   16313 B/op	     207 allocs/op
Benchmark__Gob________Encode	    3000	    214663 ns/op	   16080 B/op	     333 allocs/op
Benchmark__Bson_______Encode	    5000	    192314 ns/op	   34224 B/op	     852 allocs/op
Benchmark__VMsgpack___Encode	   10000	     93946 ns/op	   17425 B/op	     281 allocs/op
Benchmark__Msgp_______Encode	   30000	     29468 ns/op	    1984 B/op	       8 allocs/op
Benchmark__Gcbor_______Encode	   10000	    103720 ns/op	    6768 B/op	     330 allocs/op
Benchmark__Xdr________Encode	--- FAIL: _snip_
Benchmark__Sereal_____Encode	    5000	    147433 ns/op	   27814 B/op	     481 allocs/op

Benchmark__Noop_______Decode	  100000	      5391 ns/op	    2060 B/op	       4 allocs/op
Benchmark__Msgpack____Decode	    5000	    131391 ns/op	   12880 B/op	     370 allocs/op
Benchmark__Binc_NoSym_Decode	    5000	    131739 ns/op	   12720 B/op	     358 allocs/op
Benchmark__Binc_Sym___Decode	    5000	    159744 ns/op	   17376 B/op	     265 allocs/op
Benchmark__Simple_____Decode	    5000	    131010 ns/op	   12688 B/op	     358 allocs/op
Benchmark__Cbor_______Decode	    5000	    131370 ns/op	   12800 B/op	     366 allocs/op
Benchmark__Json_______Decode	    3000	    213812 ns/op	   15360 B/op	     455 allocs/op
Benchmark__Std_Json___Decode	    2000	    362146 ns/op	   17992 B/op	     629 allocs/op
Benchmark__Gob________Decode	    2000	    507253 ns/op	   79744 B/op	    2041 allocs/op
Benchmark__Bson_______Decode	    3000	    216456 ns/op	   19976 B/op	    1246 allocs/op
Benchmark__VMsgpack___Decode	    5000	    150931 ns/op	   17744 B/op	     548 allocs/op
Benchmark__Msgp_______Decode	   10000	     51247 ns/op	    8904 B/op	     200 allocs/op
Benchmark__Gcbor_______Decode	2014/12/29 11:27:02 Error: _snip_
Benchmark__Xdr________Decode	--- FAIL: _snip_
Benchmark__Sereal_____Decode	--- FAIL: _snip_

```
