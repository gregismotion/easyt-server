package storage

type NameReference struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type ReferenceCollection struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Data ReferenceGroups `json:"data"`
}
