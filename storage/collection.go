package storage

type Collection struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Data DataWrappers `json:"type"`
}
