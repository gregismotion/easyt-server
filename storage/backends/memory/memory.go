package memory

import (
	"git.freeself.one/thegergo02/easyt/storage"
	"git.freeself.one/thegergo02/easyt/basic"
	
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	collections []storage.Collection
	namedTypes []storage.NamedType
}

func New() *MemoryStorage {
	return &(
		MemoryStorage {
			collections: make([]storage.Collection, 0),
			namedTypes: make([]storage.NamedType, 0),
		})
}

func (memory MemoryStorage) GetCollectionReferences() (*[]storage.NameReference, error) {
	var references []storage.NameReference
	for _, collection := range memory.collections {
		references = append(references, storage.NameReference { Id: collection.Id, Name: collection.Name })
	}
	if references == nil { references = make([]storage.NameReference, 0) }
	return &references, nil
}

func (memory *MemoryStorage) CreateCollectionByName(name string, namedTypeIds []string) (*storage.NameReference, error) {
	namedTypes, ok := memory.getNamedTypesByIds(namedTypeIds)
	if ok {
		collection := storage.Collection {
			Id: uuid.New().String(),
			Name: name,
			Data: make(storage.DataWrappers),
		}
		for _, namedType := range namedTypes {
			collection.Data[namedType] = make([]storage.DataWrapper, 0)
		}
		memory.collections = append(memory.collections, collection)
		return &(storage.NameReference { Id: collection.Id, Name: collection.Name }), nil
	} else { return nil, fmt.Errorf("create collection: %q: Failed to get namedtype!", name) }
}

func (memory MemoryStorage) GetCollectionById(id string) (*storage.Collection, error) {
	var err error
	collectionPointer, ok := memory.getCollectionPointerById(id)
	if !ok { err = fmt.Errorf("get collection: %q: %v", id, storage.ErrFailedSearch) }
	return collectionPointer, err
}

func (memory *MemoryStorage) DeleteCollectionById(id string) (err error) {
	i := 0
	ok := false
	for _, elem := range memory.collections {
		if elem.Id != id {
			memory.collections[i] = elem
			i++
		} else {
			ok = true
		}
	}
	if !ok { err = fmt.Errorf("delete collection: %q: %w", id, storage.ErrFailedDeletion) }
	memory.collections = memory.collections[:i]
	return
}


func (memory *MemoryStorage) AddDataToCollectionById(namedTypeId string, time time.Time, value string, id string) (*storage.DataWrapper, error) {
	collection, ok := memory.getCollectionPointerById(id)
	namedType, err := memory.GetNamedTypeById(namedTypeId)
	if ok {
		if err == nil {
			dataWrapper := storage.DataWrapper {
				Id: uuid.New().String(),
				Time: time,
				Value: value,
			}
			(*collection).Data[*namedType] = append((*collection).Data[*namedType], dataWrapper)
			return &dataWrapper, nil
		} else {
			return nil, fmt.Errorf("add data: namedtype: %q: %w", namedTypeId, storage.ErrFailedSearch)
		}
	} else {
		return nil, fmt.Errorf("add data: collection: %q: %w", id, storage.ErrFailedSearch)
	}
}

func (memory MemoryStorage) GetDataInCollectionById(colId string, dataId string) (*storage.DataWrapper, error) {
	collection, err := memory.GetCollectionById(colId)
	if err != nil { return nil, fmt.Errorf("get data: collection: %q: %w", colId, storage.ErrFailedSearch) }
	for _, dataWrappers := range collection.Data {
		for _, data := range dataWrappers {
			if data.Id == dataId {
				return &data, nil
			}
		}	
	}
	return nil, fmt.Errorf("get data: %q: %w", dataId, storage.ErrFailedSearch)
}

func (memory *MemoryStorage) DeleteDataFromCollectionById(colId string, dataId string) (error) {
	last := false
	found := false
	collection, ok := memory.getCollectionPointerById(colId)
	if ok {
		for namedType, _ := range (*collection).Data {
			if !last {
				i := 0
				for _, elem := range (*collection).Data[namedType] {
					if elem.Id != dataId {
						(*collection).Data[namedType][i] = elem
						i++
						last = true
					} else {
						found = true
					}
				}
				(*collection).Data[namedType] = (*collection).Data[namedType][:i]
			} else { if found { return nil } else { return fmt.Errorf("delete data: %q: %w", dataId, storage.ErrFailedDeletion) } }
		}
	}
	return fmt.Errorf("delete data: collection: %q: %w", colId, storage.ErrFailedSearch)
}


func (memory MemoryStorage) GetNamedTypes() (*[]storage.NamedType, error) {
	return &(memory.namedTypes), nil
}

func (memory MemoryStorage) GetNamedTypeById(id string) (*storage.NamedType, error) {
	for _, namedType := range memory.namedTypes {
		if namedType.Id == id {
			return &namedType, nil
		}
	}
	return nil, fmt.Errorf("get namedtype: %q: %w", id, storage.ErrFailedSearch)
}

func (memory *MemoryStorage) CreateNamedType(name string, basicName string) (*storage.NamedType, error) {
	basicType, ok := basic.StrToBasicType(basicName)
	if ok {
		namedType := storage.NamedType {
			Id: uuid.New().String(),
			Name: name,
			Type: basicType,
		}
		memory.namedTypes = append(memory.namedTypes, namedType)
		return &namedType, nil
	} else {
		return nil, fmt.Errorf("create namedtype: %q, %q: Failed to convert str to basictype!", name, basicName)
	}
}

func (memory *MemoryStorage) DeleteNamedTypeById(id string) (err error) {
	i := 0
	ok := false
	for _, elem := range memory.namedTypes {
		if elem.Id != id {
			memory.namedTypes[i] = elem
			i++
		} else {
			ok = true
		}
	}
	memory.namedTypes = memory.namedTypes[:i]
	if !ok { err = fmt.Errorf("delete namedtype: %q: %w", id, storage.ErrFailedDeletion) }
	return
}


func (memory MemoryStorage) getNamedTypesByIds(ids []string) (namedTypes []storage.NamedType, ok bool) {
	for _, id := range ids {
		namedType, err := memory.GetNamedTypeById(id)
		if err == nil {
			namedTypes = append(namedTypes, *namedType)
			ok = true
		}
	}
	return
}

func (memory MemoryStorage) getCollectionPointerById(id string) (*storage.Collection, bool) {
	for _, collection := range memory.collections {
		if collection.Id == id {
			return &collection, true
		}
	}
	return nil, false
}
