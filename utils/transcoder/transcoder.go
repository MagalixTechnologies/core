package transcoder

import (
	"encoding/json"
	"github.com/MagalixCorp/magalix-agent/proto"
	"github.com/reconquest/karma-go"
)

type Encoding string
const (
	EncodingGob Encoding = "gob"
	EncodingSnappy Encoding = "snappy"
	EncodingJSON Encoding = "json"
)

type Transcoder struct {
	Encoding Encoding
}

func (t *Transcoder) Encode(in interface{}) ([]byte, error) {
	switch t.Encoding {
	case EncodingSnappy:
		return proto.EncodeSnappy(in)
	case EncodingGob:
		return proto.EncodeGOB(in)
	case EncodingJSON:
		return json.Marshal(in)
	default:
		return nil, karma.Format(nil, "Unsupported encoding %s", t.Encoding)
	}
}

func (t *Transcoder)  Decode(in []byte, out interface{}) error {
	switch t.Encoding {
	case EncodingSnappy:
		return proto.DecodeSnappy(in, out)
	case EncodingGob:
		return proto.DecodeGOB(in, out)
	case EncodingJSON:
		return json.Unmarshal(in, out)
	default:
		return karma.Format(nil, "Unsupported encoding %s", t.Encoding)
	}
}