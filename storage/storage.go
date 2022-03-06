package storage

type Storage interface {
	GetCollectionReferences()      ([]NameReference, bool)
	IsCollectionExistentById(string) (bool)
	CreateCollection(Collection)   (bool)
	AddToCollectionById(DataWrapper, NamedType, string) (bool)
	
	GetDataInCollectionById(string, string) (DataWrapper, bool)

	GetNamedTypeById(string)     (NamedType, bool)
	DeleteNamedTypeById(string) (bool)
}


