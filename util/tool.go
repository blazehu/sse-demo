package util

import (
	"bytes"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
	"io"
)

type CustomTranscoder struct {
	marshaler runtime.Marshaler
}

func NewCustomTranscoder(marshaler runtime.Marshaler) *CustomTranscoder {
	return &CustomTranscoder{marshaler: marshaler}
}

func (c *CustomTranscoder) ContentType(_ interface{}) string {
	return "text/event-stream"
}

func (c *CustomTranscoder) Marshal(v interface{}) ([]byte, error) {
	var jsonBytes []byte
	var err error

	if pb, ok := v.(proto.Message); ok {
		// Marshal message to JSON
		jsonBytes, err = c.marshaler.Marshal(pb)
	} else {
		// If not a proto.Message, try to marshal it as a regular JSON object
		jsonBytes, err = json.Marshal(v)
	}

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("data: ")
	buf.Write(jsonBytes)
	buf.WriteString("\n\n")

	return buf.Bytes(), nil
}

func (c *CustomTranscoder) Unmarshal(data []byte, v interface{}) error {
	return c.marshaler.Unmarshal(data, v)
}

func (c *CustomTranscoder) NewDecoder(r io.Reader) runtime.Decoder {
	return c.marshaler.NewDecoder(r)
}

func (c *CustomTranscoder) NewEncoder(w io.Writer) runtime.Encoder {
	return c.marshaler.NewEncoder(w)
}
