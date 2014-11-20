# go-codec-bench

This is a comparison of different binary and text encodings.

We compare the codecs provided by github.com/ugorji/go/codec package,
against other libraries:

github.com/ugorji/go/codec [http://godoc.org/github.com/ugorji/go/codec] provides:

  - msgpack: [https://github.com/msgpack/msgpack] 
  - binc:    [http://github.com/ugorji/binc]
  - cbor:    [http://cbor.io] [http://tools.ietf.org/html/rfc7049]
  - simple: 
  - json:    [http://json.org] [http://tools.ietf.org/html/rfc7159] 

Other codecs compared include:

  - [http://godoc.org/github.com/ugorji/go/codec] github.com/vmihailenco/msgpack
  - [http://godoc.org/labix.org/v2/mgo/bson] labix.org/v2/mgo/bson


# Data

The data being serialized is the following structure with randomly generated values:

The -bd parameter tells how many times the structure contains itself recursively.

```
type TestStruc struct {
	S    string
	I64  int64
	I16  int16
	Ui64 uint64
	Ui8  uint8
	B    bool
	By   byte

	Sslice    []string
	I64slice  []int64
	I16slice  []int16
	Ui64slice []uint64
	Ui8slice  []uint8
	Bslice    []bool
	Byslice   []byte

	Islice    []interface{}
	Iptrslice []*int64

	AnonInTestStruc

	//M map[interface{}]interface{}  `json:"-",bson:"-"`
	Ms    map[string]interface{}
	Msi64 map[string]int64

	Nintf      interface{} //don't set this, so we can test for nil
	T          time.Time
	Nmap       map[string]bool //don't set this, so we can test for nil
	Nslice     []byte          //don't set this, so we can test for nil
	Nint64     *int64          //don't set this, so we can test for nil
	Mtsptr     map[string]*TestStruc
	Mts        map[string]TestStruc
	Its        []*TestStruc
	Nteststruc *TestStruc
}

type AnonInTestStruc struct {
	AS        string
	AI64      int64
	AI16      int16
	AUi64     uint64
	ASslice   []string
	AI64slice []int64
}

```

# Run Benchmarks

```
go get -u -t 
go test -bench=.

# To test with a larger struct, use the -bd parameter. e.g.
go test -bench=. -bd=2 -bi

# To see all the test parameters, use the -Z parameter e.g.
go test -Z

```

# Test Results

Results on my 2012 i7-2630QM CPU @ 2.00GHz running Ubuntu 14.04 x86_64 GNU/Linux:

```
$ go test -bench=. -benchmem -bi -bu
..............................................
BENCHMARK INIT: 2014-11-19 21:01:46.072265189 -0500 EST
__ snip __
Benchmark: 
	Struct recursive Depth:             1
	ApproxDeepSize Of benchmark Struct: 4462 bytes
Benchmark One-Pass Run (with Unscientific Encode/Decode times): 
	   msgpack: len: 1580 bytes, encode: 289.606us, decode: 301.962us
	binc-nosym: len: 1552 bytes, encode: 137.693us, decode: 217.097us
	  binc-sym: len: 1179 bytes, encode: 138.842us, decode: 217.388us
	    simple: len: 1889 bytes, encode: 126.381us, decode: 236.93us
	      cbor: len: 1584 bytes, encode: 117.908us, decode: 210.059us
	      json: len: 2554 bytes, encode: 188.474us, decode: 425.691us
	  std-json: len: 2546 bytes, encode: 369.61us, decode: 597.577us
	       gob: len: 1992 bytes, encode: 514.737us, decode: 728.763us
	 v-msgpack: len: 1628 bytes, encode: 178.522us, decode: 150.268us
	      bson: len: 3025 bytes, encode: 228.887us, decode: 242.065us
..............................................
PASS
Benchmark__Msgpack____Encode	   50000	     68560 ns/op	   16628 B/op	      93 allocs/op
Benchmark__Msgpack____Decode	   10000	    116583 ns/op	   14462 B/op	     258 allocs/op
Benchmark__Binc_NoSym_Encode	   50000	     68423 ns/op	   16592 B/op	      93 allocs/op
Benchmark__Binc_NoSym_Decode	   10000	    113283 ns/op	   13208 B/op	     242 allocs/op
Benchmark__Binc_Sym___Encode	   20000	     82808 ns/op	   18780 B/op	      97 allocs/op
Benchmark__Binc_Sym___Decode	   10000	    125974 ns/op	   14877 B/op	     204 allocs/op
Benchmark__Simple_____Encode	   50000	     70421 ns/op	   16629 B/op	      93 allocs/op
Benchmark__Simple_____Decode	   10000	    116058 ns/op	   13630 B/op	     246 allocs/op
Benchmark__Cbor_______Encode	   50000	     68901 ns/op	   16628 B/op	      93 allocs/op
Benchmark__Cbor_______Decode	   10000	    112324 ns/op	   13630 B/op	     246 allocs/op
Benchmark__Json_______Encode	   20000	     97063 ns/op	   21665 B/op	     102 allocs/op
Benchmark__Json_______Decode	   10000	    212053 ns/op	   14469 B/op	     273 allocs/op
Benchmark__Std_Json___Encode	   20000	     92538 ns/op	   14807 B/op	     132 allocs/op
Benchmark__Std_Json___Decode	   10000	    273706 ns/op	   12840 B/op	     265 allocs/op
Benchmark__Gob________Encode	   10000	    149301 ns/op	   22100 B/op	     222 allocs/op
Benchmark__Gob________Decode	    5000	    420676 ns/op	   79633 B/op	    1656 allocs/op
Benchmark__Bson_______Encode	   10000	    130600 ns/op	   26119 B/op	     405 allocs/op
Benchmark__Bson_______Decode	   10000	    157508 ns/op	   14768 B/op	     422 allocs/op
Benchmark__VMsgpack___Encode	   50000	     64604 ns/op	   10510 B/op	     107 allocs/op
Benchmark__VMsgpack___Decode	   10000	    128225 ns/op	   14274 B/op	     270 allocs/op
```
