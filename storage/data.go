package storage

import (
	"time"
)

type DataPoint struct {
	Id string `json:"id,omitempty"`
	Time time.Time `json:"time,omitempty"` // TODO: try tinytime, we don't need nanosecond precision...
	Value string `json:"value" required=true`
	NamedType NamedType `json:"named_type" required=true`
}

func (dataPoint DataPoint) ToReference() *DataReference {
	return &(DataReference { Id: dataPoint.Id, NamedType: dataPoint.NamedType, Time: dataPoint.Time })
}

type DataReference struct {
	Id string `json:"id"`
	Time time.Time `json:"time"`
	NamedType NamedType `json:"named_type"`
}

// Group together DataReferences
type ReferenceGroups map[string][]DataReference
