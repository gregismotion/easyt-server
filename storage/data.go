package storage

import (
	"time"
)

type DataPoint struct {
	Id        string    `json:"id,omitempty" example:"237e9877-e79b-12d4-a765-321741963000"`
	Time      time.Time `json:"time,omitempty"` // TODO: try tinytime, we don't need nanosecond precision...
	Value     string    `json:"value" example:"some string data..."`
	NamedType NamedType `json:"named_type" example:"str"`
}

func (dataPoint DataPoint) ToReference() *DataReference {
	return &(DataReference{Id: dataPoint.Id, NamedType: dataPoint.NamedType, Time: dataPoint.Time})
}

type DataReference struct {
	Id        string    `json:"id" example:"237e9877-e79b-12d4-a765-321741963000"`
	Time      time.Time `json:"time"`
	NamedType NamedType `json:"named_type" example:"num"`
}

// Group together DataReferences
type ReferenceGroups map[string][]DataReference
