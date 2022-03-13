package storage

import (
	"time"
	//"bytes"
)

type DataPoint struct {
	Id string `json:"id"`
	Time time.Time `json:"time"` // TODO: try tinytime, we don't need nanosecond precision...
	Value string `json:"value"`
	NamedType NamedType `json:"named_type"`
}
type DataReference struct {
	Id string `json:"id"`
	Time time.Time `json:"time"`
	NamedType NamedType `json:"named_type"`
}

/*func (data DataPoint) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("")), nil
}*/

/*func (dataWrapper DataWrapper) MarshalJSON() ([]byte, error) {
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
}*/

// Group together DataReferences
type ReferenceGroups map[string][]DataReference

/*func (data DataWrappers) MarshalJSON() ([]byte, error) {
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
}*/
