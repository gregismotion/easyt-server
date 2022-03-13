package body

type CollectionRequestBody struct {
	Name 	     string `json:"name"`
}

type DataRequestBody struct {
	Time 	  string `json:"time,omitempty"`
	NamedType string `json:"named_type"`
	Value 	  string `json:"value"`
}

type NamedTypeRequestBody struct {
	BasicType string `json:"basic_type"`
	Name string `json:"name"`
}
