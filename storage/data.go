package storage

import (
	"git.freeself.one/thegergo02/easyt/basic"
	"time"
	"bytes"
)

type DataWrapper struct {
	Id string `json:"id"`
	Time time.Time `json:"time"` // TODO: try tinytime, we don't need nanosecond precision...
	Value string `json:"value"`
	Type basic.BasicType `json:"type"`
}

// TODO: nicer string formatting, this looks ugly rn
func (dataWrapper DataWrapper) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`{"time":"`)
	buffer.WriteString(dataWrapper.Time.String())
	buffer.WriteString(`","id":"`)
	buffer.WriteString(dataWrapper.Id)
	buffer.WriteString(`","type":"`)
	buffer.WriteString(dataWrapper.Id)
	buffer.WriteString(`","data":"`)
	buffer.WriteString(dataWrapper.Value)
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
