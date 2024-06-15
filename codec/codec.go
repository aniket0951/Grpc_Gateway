package codec

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	fmt.Println("init codec")
	encoding.RegisterCodec(NewGrpcJsonCodec())
}

const (
	HttpJsonRPCCodec = "gw-codec"
)

type WrappedBytesPb struct {
	Payload []byte
}

type GrpcJsonCodec struct {
	protojson.MarshalOptions
	protojson.UnmarshalOptions
}

func (c *GrpcJsonCodec) Name() string {
	return HttpJsonRPCCodec
}

func (c *GrpcJsonCodec) Marshal(v interface{}) ([]byte, error) {
	if raw, ok := v.(WrappedBytesPb); ok {
		return raw.Payload, nil
	}

	err := fmt.Errorf("gateway doesn't implement marshal for unknown messages, argtype: %T", v)
	logrus.Error(err)
	return nil, err
}

func (c *GrpcJsonCodec) Unmarshal(b []byte, v interface{}) error {

	if dyn, ok := v.(*WrappedBytesPb); ok {
		dyn.Payload = make([]byte, len(b))
		copy(dyn.Payload, b)
		return nil
	}
	logrus.Info("Unmarshal Error from codec on Gateway : ")
	return fmt.Errorf("unsupported type: %T", v)
}

func NewGrpcJsonCodec() *GrpcJsonCodec {
	return &GrpcJsonCodec{
		MarshalOptions:   protojson.MarshalOptions{},
		UnmarshalOptions: protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true},
	}
}
