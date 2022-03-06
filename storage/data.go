package storage

import (
	"git.freeself.one/thegergo02/easyt/basic"
	"fmt"
	"time"
	"bytes"
)

type DataWrapper struct {
	Id string `json:"id"`
	Time time.Time `json:"time"` // TODO: try tinytime, we don't need nanosecond precision...
	Type basic.BasicType `json:"type"`
	Num float64 `json:"num"`
	Str string `json:"str"`
}

// TODO: nicer string formatting, this looks ugly rn
func (data DataWrapper) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`{"time":"`)
	buffer.WriteString(data.Time.String())
	buffer.WriteString(`","type":"`)
	buffer.WriteString(data.Type.String())
	buffer.WriteString(`","id":"`)
	buffer.WriteString(data.Id)
	buffer.WriteString(`","value":"`)
	switch data.Type {
		case basic.Num:
			buffer.WriteString(fmt.Sprintf("%.5f", data.Num)) // TODO: what precision do we need?
		case basic.Str:
			buffer.WriteString(data.Str)
		default:
			buffer.WriteString("unknown")

	}
	buffer.WriteString(`"}`)
	return buffer.Bytes(), nil
}

type DataWrappers map[NamedType][]DataWrapper

func (data DataWrappers) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	first := true
	for namedType, dataWrappers := range data {
		if !first { buffer.WriteString(`,`) } else { first = false }
		buffer.WriteString(`"`)
		buffer.WriteString(namedType.Name)
		buffer.WriteString(`":[`)
		dFirst := true
		for _, dataWrapper := range dataWrappers {
			if !dFirst { buffer.WriteString(`,`) } else { dFirst = false }
			bytes, err := dataWrapper.MarshalJSON()
			if err != nil { return buffer.Bytes(), err }
			buffer.Write(bytes)
		}
		buffer.WriteString(`]`)
	}
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}
