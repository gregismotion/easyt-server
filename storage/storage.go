package storage

import "time"

type Storage interface {
	GetCollectionReferences()               	    ([]NameReference, bool)
	IsCollectionExistentById(string)        	    (bool)
	CreateCollectionByName(string, []string)   	    (NameReference, bool)
	GetCollectionById(string) 			    (Collection, bool)
	DeleteCollectionById(string)   			    (bool)
	
	AddDataToCollectionById(string, time.Time, string, string) (DataWrapper, bool)
	GetDataInCollectionById(string, string) 	    	   (DataWrapper, bool)
	DeleteDataFromCollectionById(string, string) 		   (bool)

	GetNamedTypes()             ([]NamedType, bool)
	GetNamedTypeById(string)    (NamedType, bool)
	CreateNamedType(string, string) (NamedType, bool)
	DeleteNamedTypeById(string) (bool)
}


