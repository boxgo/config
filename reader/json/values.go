package json

import (
	"github.com/boxgo/config/reader"
	"github.com/boxgo/config/source"
	jsoniter "github.com/json-iterator/go"
)

type jsonValues struct {
	ch  *source.ChangeSet
	api jsoniter.API
}

func newValues(ch *source.ChangeSet) (reader.Values, error) {
	api := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "config",
	}.Froze()

	return &jsonValues{
		ch:  ch,
		api: api,
	}, nil
}

func (j *jsonValues) Bytes() []byte {
	return j.ch.Data
}

func (j *jsonValues) Get(path ...interface{}) reader.Value {
	p := []interface{}{}
	for _, item := range path {
		p = append(p, item)
	}

	return &jsonValue{
		api:   j.api,
		value: j.api.Get(j.Bytes(), p...),
	}
}

func (j *jsonValues) Map() map[string]interface{} {
	return j.api.Get(j.ch.Data).GetInterface().(map[string]interface{})
}

func (j *jsonValues) Scan(v interface{}) error {
	b, err := j.api.Marshal(j.ch.Data)
	if err != nil {
		return err
	}
	return j.api.Unmarshal(b, v)
}
