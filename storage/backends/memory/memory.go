package memory

import (
	"git.freeself.one/thegergo02/easyt/basic"
	"git.freeself.one/thegergo02/easyt/storage"

	"fmt"
	//"time"

	"github.com/google/uuid"
)

// FIXME: if a group becomes empty, it will not get deleted!

// Group together DataPoints
type DataGroups map[string][]storage.DataPoint
type Collection struct {
	Id   string
	Name string
	Data DataGroups
}

func (collection *Collection) GetData(groupId, dataId string) (*storage.DataPoint, error) {
	for _, data := range (*collection).Data[groupId] {
		if data.Id == dataId {
			return &data, nil
		}
	}
	return nil, fmt.Errorf("get data: %q: %w", dataId, storage.ErrFailedSearch)
}

func (collection *Collection) DeleteData(groupId, dataId string) error {
	i := 0
	found := false
	for _, elem := range (*collection).Data[groupId] {
		if elem.Id != dataId {
			(*collection).Data[groupId][i] = elem
			i++
		} else {
			found = true
		}
	}
	(*collection).Data[groupId] = (*collection).Data[groupId][:i]
	if !found {
		return fmt.Errorf("delete data: %q: %w", dataId, storage.ErrFailedDeletion)
	} else {
		return nil
	}
}

type MemoryStorage struct {
	collections []Collection
	namedTypes  []storage.NamedType
}

func New() *MemoryStorage {
	return &(MemoryStorage{
		collections: make([]Collection, 0),
		namedTypes:  make([]storage.NamedType, 0),
	})
}

func (memory MemoryStorage) GetCollectionReferences(size int, lastId string) (*[]storage.NameReference, error) {
	var references []storage.NameReference = make([]storage.NameReference, 0)
	var throughLast = false
	if len(lastId) == 0 {
		throughLast = true
	}
	for _, collection := range memory.collections {
		if throughLast {
			references = append(references, storage.NameReference{Id: collection.Id, Name: collection.Name})
		} else if collection.Id == lastId {
			throughLast = true
		}
		if len(references) >= size {
			break
		}
	}
	return &references, nil
}

func (memory *MemoryStorage) CreateCollectionByName(name string) (*storage.NameReference, error) {
	collection := Collection{
		Id:   uuid.New().String(),
		Name: name,
		Data: make(DataGroups),
	}
	memory.collections = append(memory.collections, collection)
	return &(storage.NameReference{Id: collection.Id, Name: collection.Name}), nil
}

func (memory MemoryStorage) GetReferenceCollectionById(id string) (*storage.ReferenceCollection, error) {
	collectionPointer, ok := memory.getCollectionPointerById(id)
	if !ok {
		return nil, fmt.Errorf("get collection: %q: %v", id, storage.ErrFailedSearch)
	} else {
		collection := *collectionPointer
		referenceGroups := make(storage.ReferenceGroups)
		for groupId, dataGroup := range collection.Data {
			dataReferences := make([]storage.DataReference, len(dataGroup))
			for i, dataPoint := range dataGroup {
				dataReferences[i] = *dataPoint.ToReference()
			}
			referenceGroups[groupId] = dataReferences
		}
		return &(storage.ReferenceCollection{Id: collection.Id, Name: collection.Name, Data: referenceGroups}), nil
	}
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
	if !ok {
		err = fmt.Errorf("delete collection: %q: %w", id, storage.ErrFailedDeletion)
	}
	memory.collections = memory.collections[:i]
	return
}

// TODO: should follow the same access format as other methods on collection.
func (memory *MemoryStorage) AddDataPointsToCollectionById(colId string, dataPoints []storage.DataPoint) (*storage.ReferenceGroups, error) {
	groupId := uuid.New().String()
	references := make([]storage.DataReference, len(dataPoints))
	var err error
	for i, dataPoint := range dataPoints {
		var reference *storage.DataReference
		reference, err = (*memory).addDataPointToCollectionById(colId, groupId, dataPoint)
		if err != nil {
			return nil, err
		} else {
			references[i] = *reference
		}
	}
	groupReferences := storage.ReferenceGroups{groupId: references}
	return &groupReferences, nil
}

func (memory MemoryStorage) GetDataInCollectionById(colId, groupId, dataId string) (*storage.DataPoint, error) {
	collection, ok := memory.getCollectionPointerById(colId)
	if ok {
		data, err := collection.GetData(groupId, dataId)
		return data, err
	} else {
		return nil, fmt.Errorf("get data: collection: %q: %w", colId, storage.ErrFailedSearch)
	}
}

func (memory *MemoryStorage) DeleteDataFromCollectionById(colId, groupId, dataId string) error {
	collection, ok := memory.getCollectionPointerById(colId)
	if ok {
		return collection.DeleteData(groupId, dataId)
	} else {
		return fmt.Errorf("delete data: collection: %q: %w", colId, storage.ErrFailedSearch)
	}
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
		namedType := storage.NamedType{
			Id:   uuid.New().String(),
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
	if !ok {
		err = fmt.Errorf("delete namedtype: %q: %w", id, storage.ErrFailedDeletion)
	}
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

func (memory MemoryStorage) getCollectionPointerById(id string) (*Collection, bool) {
	for _, collection := range memory.collections {
		if collection.Id == id {
			return &collection, true
		}
	}
	return nil, false
}

func (memory *MemoryStorage) addDataPointToCollectionById(colId, groupId string, dataPoint storage.DataPoint) (*storage.DataReference, error) {
	collection, ok := memory.getCollectionPointerById(colId)
	namedType, err := memory.GetNamedTypeById(dataPoint.NamedType.Id)
	if ok {
		if err == nil {
			dataPoint.Id = uuid.New().String()
			dataPoint.NamedType = *namedType
			(*collection).Data[groupId] = append((*collection).Data[groupId], dataPoint)
			return dataPoint.ToReference(), nil
		} else {
			return nil, fmt.Errorf("add data: namedtype: %q: %w", dataPoint.NamedType.Id, storage.ErrFailedSearch)
		}
	} else {
		return nil, fmt.Errorf("add data: collection: %q: %w", colId, storage.ErrFailedSearch)
	}
}
