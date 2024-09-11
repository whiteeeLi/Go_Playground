package codec

import "io"

type Header struct {
	ServiceMethod string // format "Service.Method"
	Seq           uint64 // sequence number chosen by client
	Error         string
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error //这里可以看到rpc本质上是对网络通信进行了封装，所以会有对heade bo
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
