package storage

type NameReference struct {
	Id   string `json:"id" example:"237e9877-e79b-12d4-a765-321741963000"`
	Name string `json:"name" example:"some_name"`
}

type ReferenceCollection struct {
	Id   string          `json:"id" example:"237e9877-e79b-12d4-a765-321741963000"`
	Name string          `json:"name" example:"some_name"`
	Data ReferenceGroups `json:"data"`
}
