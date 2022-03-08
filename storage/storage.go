package storage

import (
	"time"
	"errors"
)

var ErrFailedDeletion = errors.New("Failed deletion!")
var ErrFailedSearch = errors.New("Failed search!")

type Storage interface {
	GetCollectionReferences()               	    (*[]NameReference, error)
	CreateCollectionByName(string, []string)   	    (*NameReference, error)
	GetCollectionById(string) 			    (*Collection, error)
	DeleteCollectionById(string)   			    (error)
	
	AddDataToCollectionById(string, time.Time, string, string) (*DataWrapper, error)
	GetDataInCollectionById(string, string) 	    	   (*DataWrapper, error)
	DeleteDataFromCollectionById(string, string) 		   (error)

	GetNamedTypes()             (*[]NamedType, error)
	GetNamedTypeById(string)    (*NamedType, error)
	CreateNamedType(string, string) (*NamedType, error)
	DeleteNamedTypeById(string) (error)
}


