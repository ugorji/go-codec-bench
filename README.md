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

```
# download the code and all its dependencies 
    go get -u -t github.com/ugorji/go-codec-bench
    go get -u \
           github.com/tinylib/msgp/msgp github.com/tinylib/msgp \
           github.com/pquerna/ffjson/ffjson github.com/pquerna/ffjson \
           github.com/Sereal/Sereal/Go/sereal \
           bitbucket.org/bodhisnarkva/cbor/go \
           github.com/davecgh/go-xdr/xdr2 \
           gopkg.in/mgo.v2/bson \
           gopkg.in/vmihailenco/msgpack.v2 \
           github.com/json-iterator/go \
           github.com/mailru/easyjson/...

# benchmark with the default settings 
go test -bench=.

# benchmark with a larger struct, using the -bd parameter. and gather pre-info
go test -bench=. -bd=2 -bi

# see all the test parameters, using the -Z parameter (any recognized param will do)
go test -Z

```

To run the full suite of benchmarks, including executing against the external frameworks
listed above, you MUST first run code generation for the frameworks that support it.

```sh
# If you want to run the benchmarks against code generated values.
# Then first generate the code generated values from values_test.go named typed.
# we cannot normally read a _test.go file, so temporarily copy it into a readable file.

z=`pwd`
z=${z%%/src/*}
cp values_test.go values_temp.go

$z/bin/msgp -tests=false -o=values_msgp_test.go -file=values_temp.go

$z/bin/ffjson -force-regenerate -reset-fields -w values_ffjson_test.go values_temp.go
sed -i '' -e 's+ MarshalJSON(+ _MarshalJSON(+g' values_ffjson_test.go
sed -i '' -e 's+ UnmarshalJSON(+ _UnmarshalJSON(+g' values_ffjson_test.go

$z/bin/easyjson -all -no_std_marshalers -omit_empty -output_filename easyjson123.go values_temp.go
mv easyjson123.go values_easyjson.go

$z/bin/codecgen -rt codecgen -t 'codecgen' -o values_codecgen_test.go -d 19780 values_temp.go

rm -f values_temp.go
```

Then you can run the tests. The fastest way is to use the bench.sh script, 
which is a simple script that runs the benchmark suite.

```sh
./bench.sh [AllSuite|XSuite|CodecSuite]
```

Feel free to run selected sets also:

```sh
# Run the tests, using only runtime introspection support (normal mode)
go test -tm -bi -benchmem '-bench=_.*En' -tags=x
go test -tm -bi -benchmem '-bench=_.*De' -tags=x

# Run the tests using the codegeneration mode for codecgen.
# This involves passing the tags which enable the appropriate files to be run.
go test -tm -tf -bi -benchmem '-bench=_.*En' '-tags=x codecgen'
go test -tm -tf -bi -benchmem '-bench=_.*De' '-tags=x codecgen'
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

Below are results in the default execution phase, without codecgen.

```
$ go test -bench=. -benchmem -tags x -bi -bu -bd 2
BENCHMARK INIT: 2017-09-18 18:03:06.563956051 -0400 EDT m=+0.003244434
Benchmark: 
	Struct recursive Depth:             2
	ApproxDeepSize Of benchmark Struct: 27046 bytes
Benchmark One-Pass Run (with Unscientific Encode/Decode times): 
	   msgpack: len: 7966 bytes,	 encode: 294.597µs,	 decode: 294.391µs
	      binc: len: 6599 bytes,	 encode: 150.412µs,	 decode: 229.514µs
	    simple: len: 8921 bytes,	 encode: 137.634µs,	 decode: 205.453µs
	      cbor: len: 7966 bytes,	 encode: 146.296µs,	 decode: 211.413µs
	      json: len: 11839 bytes,	 encode: 228.125µs,	 decode: 314.687µs
	  std-json: len: 13346 bytes,	 encode: 339.198µs,	 decode: 716.138µs
	       gob: len: 7990 bytes,	 encode: 439.171µs,	 decode: 383.415µs
	 json-iter: len: 13515 bytes,	 encode: 895.545µs,	 decode: 269.785µs
	 v-msgpack: len: 9233 bytes,	 encode: 281.534µs,	 decode: 376.967µs
	      bson: len: 14817 bytes,	 encode: 324.643µs,	 decode: 434.032µs
	      msgp: **** Error encoding benchTs: msgp: input is not a msgp.Encodable
	      msgp: **** Error decoding into new TestStruc: msgp: input is not a msgp.Decodable
	      msgp: len: 0 bytes,	 encode: 20.848µs,	 decode: 12.936µs
	  easyjson: **** Error encoding benchTs: easyjson: input is not a easyjson.Marshaler
	  easyjson: **** Error decoding into new TestStruc: easyjson: input is not a easyjson.Unmarshaler
	  easyjson: len: 0 bytes,	 encode: 16.915µs,	 decode: 13.123µs
	    ffjson: len: 13346 bytes,	 encode: 223.291µs,	 decode: 641.494µs <error decoding map val: ... snip>
	     gcbor: **** Error decoding into new TestStruc: can't read map into *codec.AnonInTestStrucIntf
	     gcbor: len: 9103 bytes,	 encode: 230.39µs,	 decode: 116.141µs
	       xdr: **** Error encoding benchTs: xdr:encodeInterface: can't encode nil interface
	       xdr: **** Error decoding into new TestStruc: xdr:decodeInterface: can't decode to nil interface
	       xdr: len: 672 bytes,	 encode: 70.214µs,	 decode: 61.248µs
	    sereal: **** Error decoding into new TestStruc: reflect: call of reflect.Value.Set on zero Value
	    sereal: len: 3717 bytes,	 encode: 191.948µs,	 decode: 89.543µs
..............................................
goos: darwin
goarch: amd64
pkg: github.com/ugorji/go-codec-bench
Benchmark__Msgpack____Encode-8   	   20000	     64973 ns/op	    5743 B/op	      45 allocs/op
Benchmark__VMsgpack___Encode-8   	   10000	    185104 ns/op	   47072 B/op	     282 allocs/op
Benchmark__Binc_______Encode-8   	   20000	     73947 ns/op	    6576 B/op	      47 allocs/op
Benchmark__Simple_____Encode-8   	   20000	     70273 ns/op	    5743 B/op	      45 allocs/op
Benchmark__Cbor_______Encode-8   	   20000	     66264 ns/op	    5744 B/op	      45 allocs/op
Benchmark__Json_______Encode-8   	   10000	    125379 ns/op	    6285 B/op	      45 allocs/op
Benchmark__Std_Json___Encode-8   	   10000	    154258 ns/op	   65248 B/op	     560 allocs/op
Benchmark__JsonIter___Encode-8   	   10000	    138191 ns/op	   40496 B/op	    1013 allocs/op
Benchmark__Bson_______Encode-8   	   10000	    217081 ns/op	  121366 B/op	    1496 allocs/op
Benchmark__Gob________Encode-8   	   10000	    161064 ns/op	   47229 B/op	     523 allocs/op
Benchmark__Sereal_____Encode-8   	   10000	    112755 ns/op	   55285 B/op	    1035 allocs/op
--- FAIL: Benchmark__Xdr________Encode <snip - error - xdr: xdr:encodeInterface: can't encode nil interface>
Benchmark__Gcbor______Encode-8   	   <snip - error>
Benchmark__Msgpack____Decode-8   	   10000	    128177 ns/op	   34688 B/op	     687 allocs/op
Benchmark__Binc_______Decode-8   	   10000	    141204 ns/op	   36576 B/op	     715 allocs/op
Benchmark__Simple_____Decode-8   	   10000	    125495 ns/op	   34064 B/op	     687 allocs/op
Benchmark__Cbor_______Decode-8   	   10000	    132491 ns/op	   34064 B/op	     687 allocs/op
Benchmark__Json_______Decode-8   	   10000	    199486 ns/op	   43008 B/op	     756 allocs/op
Benchmark__Std_Json___Decode-8   	    3000	    556246 ns/op	   58320 B/op	    1967 allocs/op
Benchmark__JsonIter___Decode-8   	   10000	    217277 ns/op	   60128 B/op	    2244 allocs/op
Benchmark__Gob________Decode-8   	    5000	    275268 ns/op	   87471 B/op	    2028 allocs/op
Benchmark__Bson_______Decode-8   	    5000	    345356 ns/op	   77664 B/op	    4224 allocs/op
Benchmark__VMsgpack___Decode-8   	    5000	    275950 ns/op	   51504 B/op	    1717 allocs/op
--- FAIL: Benchmark__Gcbor______Decode <snip - error - gcbor: can't read map into *codec.AnonInTestStrucIntf>
--- FAIL: Benchmark__Xdr________Decode <snip - error - xdr: xdr:encodeInterface: can't encode nil interface
--- FAIL: Benchmark__Sereal_____Decode <snip - error - sereal: reflect: call of reflect.Value.Set on zero Value
```

These results get better with codecgen. Our numerous tests show about 20-50% performance improvement
with codecgen. We think that users should carefully weight the performance improvements against the 
usability and binary-size increases. Already, the performance is extremely good without the codecgen path.

```
$ go test -bench=. -benchmem -tags "codecgen x" -bi -bu -bd 2
BENCHMARK INIT: 2017-09-18 18:53:12.918740626 -0400 EDT m=+0.002959732
Benchmark: 
	Struct recursive Depth:             2
	ApproxDeepSize Of benchmark Struct: 27046 bytes
Benchmark One-Pass Run (with Unscientific Encode/Decode times): 
	   msgpack: len: 7862 bytes,	 encode: 131.676µs,	 decode: 165.316µs
	      binc: len: 7584 bytes,	 encode: 67.326µs,	 decode: 183.64µs
	    simple: len: 8830 bytes,	 encode: 75.497µs,	 decode: 121.013µs
	      cbor: len: 7888 bytes,	 encode: 68.159µs,	 decode: 160.067µs
	      json: len: 11839 bytes,	 encode: 143.843µs,	 decode: 208.065µs
	  std-json: len: 13346 bytes,	 encode: 302.067µs,	 decode: 604.291µs
	       gob: len: 7990 bytes,	 encode: 392.578µs,	 decode: 324.268µs
	 json-iter: len: 13515 bytes,	 encode: 662.717µs,	 decode: 343.62µs
	 v-msgpack: len: 9233 bytes,	 encode: 262.835µs,	 decode: 349.634µs
	      bson: len: 14817 bytes,	 encode: 359.744µs,	 decode: 394.882µs
	      msgp: len: 9376 bytes,	 encode: 79.804µs,	 decode: 92.386µs
	  easyjson: len: 11839 bytes,	 encode: 161.526µs,	 decode: 165.198µs
panic: runtime error: invalid memory address or nil pointer dereference
	     gcbor: **** Error decoding into new TestStruc: can't read map into *codec.AnonInTestStrucIntf
	     gcbor: len: 9103 bytes,	 encode: 190.042µs,	 decode: 103.035µs
	       xdr: **** Error encoding benchTs: xdr:encodeInterface: can't encode nil interface
	       xdr: **** Error decoding into new TestStruc: xdr:decodeInterface: can't decode to nil interface
	       xdr: len: 672 bytes,	 encode: 58.712µs,	 decode: 63.731µs
	    sereal: **** Error decoding into new TestStruc: reflect: call of reflect.Value.Set on zero Value
	    sereal: len: 3730 bytes,	 encode: 195.442µs,	 decode: 56.763µs
..............................................
goos: darwin
goarch: amd64
pkg: ugorji.net/codec
Benchmark__Msgpack____Encode-8   	   50000	     35865 ns/op	    2528 B/op	      20 allocs/op
Benchmark__VMsgpack___Encode-8   	   10000	    205336 ns/op	   47072 B/op	     282 allocs/op
Benchmark__Msgp_______Encode-8   	   50000	     34997 ns/op	       0 B/op	       0 allocs/op
Benchmark__Binc_______Encode-8   	   50000	     35811 ns/op	    3248 B/op	      22 allocs/op
Benchmark__Simple_____Encode-8   	   50000	     40413 ns/op	    2528 B/op	      20 allocs/op
Benchmark__Cbor_______Encode-8   	   50000	     36796 ns/op	    2320 B/op	       7 allocs/op
Benchmark__Gcbor______Encode-8   	   10000	    168929 ns/op	   23336 B/op	    1413 allocs/op
Benchmark__Json_______Encode-8   	   20000	     86954 ns/op	    3072 B/op	      20 allocs/op
Benchmark__JsonIter___Encode-8   	   10000	    138147 ns/op	   40496 B/op	    1013 allocs/op
Benchmark__Std_Json___Encode-8   	   10000	    159704 ns/op	   65248 B/op	     560 allocs/op
Benchmark__Easyjson___Encode-8   	   20000	    100054 ns/op	   32292 B/op	     153 allocs/op
Benchmark__Ffjson_____Encode-8   	   10000	    171808 ns/op	   65445 B/op	     652 allocs/op
Benchmark__Bson_______Encode-8   	   10000	    221200 ns/op	  121380 B/op	    1496 allocs/op
Benchmark__Gob________Encode-8   	   10000	    157812 ns/op	   47230 B/op	     523 allocs/op
Benchmark__Sereal_____Encode-8   	   10000	    111683 ns/op	   55270 B/op	    1035 allocs/op
Benchmark__Msgpack____Decode-8   	   20000	     88023 ns/op	   33696 B/op	     671 allocs/op
Benchmark__VMsgpack___Decode-8   	    5000	    301914 ns/op	   51504 B/op	    1717 allocs/op
Benchmark__Msgp_______Decode-8   	   30000	     50168 ns/op	   30440 B/op	     616 allocs/op
Benchmark__Binc_______Decode-8   	   10000	    102534 ns/op	   34216 B/op	     595 allocs/op
Benchmark__Simple_____Decode-8   	   20000	     90866 ns/op	   33688 B/op	     671 allocs/op
Benchmark__Cbor_______Decode-8   	   10000	    131086 ns/op	   34104 B/op	     697 allocs/op
Benchmark__Json_______Decode-8   	   10000	    156032 ns/op	   42608 B/op	     736 allocs/op
Benchmark__JsonIter___Decode-8   	   10000	    211163 ns/op	   60128 B/op	    2244 allocs/op
Benchmark__Std_Json___Decode-8   	    3000	    552957 ns/op	   58320 B/op	    1967 allocs/op
Benchmark__Easyjson___Decode-8   	   10000	    116978 ns/op	   41008 B/op	     695 allocs/op
Benchmark__Ffjson_____Decode-8   	panic: runtime error: invalid memory address or nil pointer dereference
Benchmark__Bson_______Decode-8   	    5000	    347238 ns/op	   77664 B/op	    4224 allocs/op
Benchmark__Gob________Decode-8   	    5000	    271352 ns/op	   87485 B/op	    2028 allocs/op
--- FAIL: Benchmark__Gcbor______Decode gcbor: can't read map into *codec.AnonInTestStrucIntf
--- FAIL: Benchmark__Xdr________Encode xdr: xdr:encodeInterface: can't encode nil interface
--- FAIL: Benchmark__Xdr________Decode xdr: xdr:encodeInterface: can't encode nil interface
--- FAIL: Benchmark__Sereal_____Decode sereal: reflect: call of reflect.Value.Set on zero Value

```

Selected results for a different run using benchcmp is below:


```sh
   go test -bench "_(Json|Cbor)" -benchmem -bd 2 -tags "x" > no-codecgen.bench.out.txt
   go test -bench "_(Json|Cbor)" -benchmem -bd 2 -tags "x codecgen" > with-codecgen.bench.out.txt
   benchcmp no-codecgen.bench.out.txt with-codecgen.bench.out.txt > benchcmp.out.txt
```

Results:
```
benchmark                          old ns/op     new ns/op     delta
Benchmark__Cbor_______Encode-8     69672         37051         -46.82%
Benchmark__Json_______Encode-8     115274        87775         -23.86%
Benchmark__JsonIter___Encode-8     133558        133635        +0.06%
Benchmark__Std_Json___Encode-8     165240        157007        -4.98%
Benchmark__Cbor_______Decode-8     164396        131192        -20.20%
Benchmark__Json_______Decode-8     207014        156589        -24.36%
Benchmark__JsonIter___Decode-8     209418        208993        -0.20%
Benchmark__Std_Json___Decode-8     553650        557783        +0.75%

benchmark                          old allocs     new allocs     delta
Benchmark__Cbor_______Encode-8     32             7              -78.12%
Benchmark__Json_______Encode-8     45             20             -55.56%
Benchmark__JsonIter___Encode-8     1013           1013           +0.00%
Benchmark__Std_Json___Encode-8     560            560            +0.00%
Benchmark__Cbor_______Decode-8     713            697            -2.24%
Benchmark__Json_______Decode-8     756            736            -2.65%
Benchmark__JsonIter___Decode-8     2244           2244           +0.00%
Benchmark__Std_Json___Decode-8     1967           1967           +0.00%

benchmark                          old bytes     new bytes     delta
Benchmark__Cbor_______Encode-8     5529          2320          -58.04%
Benchmark__Json_______Encode-8     6290          3072          -51.16%
Benchmark__JsonIter___Encode-8     40496         40496         +0.00%
Benchmark__Std_Json___Encode-8     65248         65248         +0.00%
Benchmark__Cbor_______Decode-8     34384         34104         -0.81%
Benchmark__Json_______Decode-8     42976         42608         -0.86%
Benchmark__JsonIter___Decode-8     60128         60128         +0.00%
Benchmark__Std_Json___Decode-8     58320         58320         +0.00%
```
