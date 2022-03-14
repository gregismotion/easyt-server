package storage

import (
	"time"
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

// Group together DataReferences
type ReferenceGroups map[string][]DataReference
