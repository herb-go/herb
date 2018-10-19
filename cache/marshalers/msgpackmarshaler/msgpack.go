package msgpackmarshaler

import (
	"github.com/herb-go/herb/cache"
	"github.com/vmihailenco/msgpack"
)

//MsgpackMarshaler msgpack marshaler
type MsgpackMarshaler struct {
}

//Marshal Marshal data model to  bytes.
//Return marshaled bytes and any erro rasied.
func (m *MsgpackMarshaler) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

//Unmarshal Unmarshal bytes to data model.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raseid.
func (m *MsgpackMarshaler) Unmarshal(bytes []byte, v interface{}) error {
	return msgpack.Unmarshal(bytes, v)
}

func init() {
	cache.RegisterMarshaler("msgpack", func() (cache.Marshaler, error) {
		return &MsgpackMarshaler{}, nil
	})
}
