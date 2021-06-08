package transcoder

import (
	"encoding/json"

	"github.com/MagalixCorp/magalix-agent/v3/proto"
	"github.com/reconquest/karma-go"
)

type Encoding string

const (
	EncodingSnappy Encoding = "snappy"
	EncodingJSON   Encoding = "json"
)

type Transcoder struct {
	Encoding Encoding
}

func New(encoding Encoding) *Transcoder {
	return &Transcoder{Encoding: encoding}
}

func (t *Transcoder) Encode(in interface{}) ([]byte, error) {
	switch t.Encoding {
	case EncodingSnappy:
		return proto.EncodeSnappy(in)
	case EncodingJSON:
		return json.Marshal(in)
	default:
		return nil, karma.Format(nil, "Unsupported encoding %s", t.Encoding)
	}
}

func (t *Transcoder) Decode(in []byte, out interface{}) error {
	switch t.Encoding {
	case EncodingSnappy:
		return proto.DecodeSnappy(in, out)

	case EncodingJSON:
		return json.Unmarshal(in, out)
	default:
		return karma.Format(nil, "Unsupported encoding %s", t.Encoding)
	}
}
