package memory

import (
	"git.freeself.one/thegergo02/easyt/storage"
	"git.freeself.one/thegergo02/easyt/basic"
	
	//"fmt"
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

func (memory MemoryStorage) GetCollectionReferences() (references []storage.NameReference, ok bool) {
	ok = true
	for _, collection := range memory.collections {
		references = append(references, storage.NameReference { Id: collection.Id, Name: collection.Name })
	}
	return
}

func (memory MemoryStorage) IsCollectionExistentById(id string) bool {
	for _, collection := range memory.collections {
		if collection.Id == id {
			return true
		}
	}
	return false
}

func (memory *MemoryStorage) CreateCollectionByName(name string, namedTypeIds []string) (storage.NameReference, bool) {
	namedTypes, ok := memory.getNamedTypesByIds(namedTypeIds)
	if ok {
		collection := storage.Collection {
			Id: uuid.New().String(),
			Name: name,
		}
		for _, namedType := range namedTypes {
			collection.Data[namedType] = make([]storage.DataWrapper, 0)
		}
		memory.collections = append(memory.collections, collection)
		return storage.NameReference { Id: collection.Id, Name: collection.Name }, true
	} else { return storage.NameReference{}, false }
}

func (memory MemoryStorage) GetCollectionById(id string) (storage.Collection, bool) {
	collectionPointer, ok := memory.getCollectionPointerById(id)
	return *collectionPointer, ok
}

func (memory *MemoryStorage) DeleteCollectionById(id string) (ok bool) {
	i := 0
	for _, elem := range memory.collections {
		if elem.Id != id {
			memory.collections[i] = elem
			i++
		} else {
			ok = true
		}
	}
	memory.collections = memory.collections[:i]
	return
}


func (memory *MemoryStorage) AddDataToCollectionById(namedTypeId string, time time.Time, value string, id string) (storage.DataWrapper, bool) {
	collection, ok := memory.getCollectionPointerById(id)
	namedType, ok1 := memory.GetNamedTypeById(namedTypeId)
	if ok && ok1 {
		dataWrapper := storage.DataWrapper {
			Id: uuid.New().String(),
			Time: time,
			Value: value,
		}
		(*collection).Data[namedType] = append((*collection).Data[namedType], dataWrapper)
		return dataWrapper, true
	} else {
		return storage.DataWrapper{}, false
	}
}

func (memory MemoryStorage) GetDataInCollectionById(colId string, dataId string) (storage.DataWrapper, bool) {
	collection, ok := memory.GetCollectionById(colId)
	if ok {
		for _, dataWrappers := range collection.Data {
			for _, data := range dataWrappers {
				if data.Id == dataId {
					return data, true
				}
			}	
		}
	}
	return storage.DataWrapper{}, false
}

func (memory *MemoryStorage) DeleteDataFromCollectionById(colId string, dataId string) (bool) {
	last := false
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
					}
				}
				(*collection).Data[namedType] = (*collection).Data[namedType][:i]
			} else { return true }
		}
	}
	return false
}


func (memory MemoryStorage) GetNamedTypes() ([]storage.NamedType, bool) {
	return memory.namedTypes, true
}

func (memory MemoryStorage) GetNamedTypeById(id string) (storage.NamedType, bool) {
	for _, namedType := range memory.namedTypes {
		if namedType.Id == id {
			return namedType, true
		}
	}
	return storage.NamedType{}, false
}

func (memory *MemoryStorage) CreateNamedType(name string, basicName string) (namedType storage.NamedType, ok bool) {
	var basicType basic.BasicType
	basicType, ok = basic.StrToBasicType(basicName)
	if ok {
		namedType = storage.NamedType {
			Id: uuid.New().String(),
			Name: name,
			Type: basicType,
		}
		memory.namedTypes = append(memory.namedTypes, namedType)
	}
	return

}

func (memory *MemoryStorage) DeleteNamedTypeById(id string) (ok bool) {
	i := 0
	for _, elem := range memory.namedTypes {
		if elem.Id != id {
			memory.namedTypes[i] = elem
			i++
		} else {
			ok = true
		}
	}
	memory.namedTypes = memory.namedTypes[:i]
	return
}


func (memory MemoryStorage) getNamedTypesByIds(ids []string) (namedTypes []storage.NamedType, ok bool) {
	ok = true
	for _, id := range ids {
		namedType, ok1 := memory.GetNamedTypeById(id)
		if ok1 {
			namedTypes = append(namedTypes, namedType)
		} else {
			ok = false
			return
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
