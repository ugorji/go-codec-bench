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

