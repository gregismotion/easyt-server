package storage

type ReferenceCollection struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Data ReferenceGroups `json:"data"`
}
