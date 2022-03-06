/// MAIN START ///
package main

import (
	"git.freeself.one/thegergo02/easyt/basic"
	"git.freeself.one/thegergo02/easyt/storage"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var storageBackend Storage
func main() {
	r := gin.Default()
	
	v1 := r.Group("/api/v1") 
	{
		col := v1.Group("/collection")
		{
			col.GET("/", getCollections)
			col.POST("/", createCollection)
			col.GET("/:id", getCollection)
			col.DELETE("/:id", deleteCollection)
			col.POST("/:colId", addData)
			col.GET("/data/:colId/:dataId", getData)
			col.DELETE("/data/:colId/:dataId", deleteData)
		}
		typ := v1.Group("/type")
		{
			typ.GET("/named", getNamedTypes)
			typ.POST("/named", createNamedType)
			typ.GET("/:id", getNamedType)
			typ.DELETE("/:id", deleteNamedType)
			typ.GET("/basic", getBasicTypes)
		}
	}

	host := "localhost:8080"

	//storageBackend = // TODO: init storageBackend backend
	r.Run(host)
}

func getCollections(c *gin.Context) {
	/*collectionNames := make([]string, 0)
	for _, collection := range collections {
		collectionNames = append(collectionNames, collection.Name)
	}
	c.IndentedJSON(http.StatusOK, collectionNames)*/
	c.IndentedJSON(http.StatusOK, storageBackend.GetCollectionReferences())
}

type CollectionRequestBody struct {
	Name 	     string `json:"name"`
	NamedTypes []string `json:"named_types"`
}
func createCollection(c *gin.Context) {
	var body CollectionRequestBody
	if err := c.BindJSON(&body); err == nil  {
		collection := Collection {
			Name: body.Name,
			Data: make(DataWrappers),
		}
		if collection.isUnique(storageBackend) {
			for _, id := range body.NamedTypes {
				namedType, ok := storageBackend.GetNamedTypeById(id)
				if ok {
					collection.Data[namedType] = make([]DataWrapper, 0)
				} else {
					// TODO: completely fail, ignore or smt else when bad named type?
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Non-existent named type!"})
					return
				}
			}
			//collections = append(collections, collection)
			if storageBackend.createCollection(collection) {
				c.IndentedJSON(http.StatusOK, collection)
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Couldn't create collection!"})
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Duplicate name!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request body!"})
	}
}

func getCollection(c *gin.Context) { // TODO: add return limit of data
	id := c.Param("name")
	if id != "" {
		collection, ok := storageBackend.GetCollectionById(id)
		if ok {
			c.IndentedJSON(http.StatusOK, collection)
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Couldn't find collection with this ID!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No ID passed!"})
	}
}

func deleteCollection(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		if storageBackend.DeleteCollectionById(id) {
			c.String(http.StatusOK, "") // NOTE: maybe some message would be appropiate? consult the do- oh wait
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't delete collection!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No ID specified!"})
	}
}

type DataRequestBody struct {
	Time 	  string `json:"time,omitempty"`
	NamedType string `json:"named_type"`
	Value 	  string `json:"value"`
}
func addData(c *gin.Context) {
	var body DataRequestBody
	id := c.Param("id")
	if id != "" {
		if c.BindJSON(&body) == nil {
			if storageBackend.IsCollectionExistentById(id) {
				time := time.Now() // TODO: get time from body
				namedType, okTyp := storageBackend.GetNamedTypeById(body.NamedType)
				if okTyp {
					dataWrapper := DataWrapper {
						Id: uuid.New().String(),
						Time: time,
						Type: namedType.Type,
					}
					value := body.Value
					switch namedType.Type { // TODO: should return error at unparseable values
						case basic.Num:
							if n, err := strconv.ParseFloat(value, 64); err == nil {
								dataWrapper.Num = n
							} else {
								dataWrapper.Str = value
							}
						case basic.Str:
							dataWrapper.Str = value
						default:
							dataWrapper.Str = value
					}
					storageBackend.AddDataToCollectionById(dataWrapper, namedType, id)
					c.IndentedJSON(http.StatusCreated, dataWrapper)
				} else {
					c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Named type does not exist!"})
				}
				
			} else {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find collection!"})
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request body!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No ID specified!"})
	}
}

func getData(c *gin.Context) {
	colId := c.Param("colId")
	if colId != "" {
		if storageBackend.IsCollectionExistentById(colId) {
			dataId := c.Param("dataId")
			if dataId != "" {
				data, ok := storageBackend.GetDataInCollectionById(colId, dataId)
				if ok {
					c.IndentedJSON(http.StatusOK, data)
				} else {
					c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find data!"})
				}
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No data ID specified!"})
			}
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find collection!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No collection ID specified!"})
	}
}

func deleteData(c *gin.Context) {
	colId := c.Param("colId")
	if colId != "" {
		if storageBackend.IsCollectionExistentById(colId) {
			dataId := c.Param("dataId")
			if dataId != "" {
				if storageBackend.RemoveDataFromCollectionById(colId, dataId) {
					c.IndentedJSON(http.StatusOK, "")
				} else {
					c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find data!"})
				}
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No data ID specified!"})
			}
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find collection!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}


func getNamedTypes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, storageBackend.GetNamedTypes())
}

type NamedTypeRequestBody struct {
	BasicType string `json:"basic_type"`
	Name string `json:"name"`
}
func createNamedType(c *gin.Context) {
	var body NamedTypeRequestBody
	if c.BindJSON(&body) {
		typ, ok := basic.StrToBasicType(body.NamedType)
		if ok {
			namedType := NamedType {
				Id: uuid().New().String(),
				Name: body.Name,
				Type: typ,
			}
			if namedType.isUnique(storageBackend) {
				//namedTypes = append(namedTypes, namedType)
				storageBackend.CreateNamedType(namedType)
				c.IndentedJSON(http.StatusCreated, namedType)
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Duplicate name!"})
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unknown basic type!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request body!"})
	}
}

func getNamedType(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		namedType, ok := GetNamedTypeById(id)
		if ok {
			c.IndentedJSON(http.StatusOK, namedType)
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find named type!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No ID specified!"})
	}
}

func deleteNamedType(c *gin.Context) {
	id := c.Param("id")
	if name != "" {
		if storageBackend.DeleteNamedTypeById(id) {
			c.String(http.StatusOK, "") // NOTE: maybe some message would be appropiate? consult the do- oh wait
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't delete named type!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}

func getBasicTypes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, basic.GetBasicTypes())
}
/// MAIN END ///



/// STORAGE START ///
// NOTE: might become unmanageable, find alternative
func idToData(collection Collection, id string) (DataWrapper, bool){
	for _, dataWrappers := range collection.Data {
		for _, data := range dataWrappers {
			if data.Id == id {
				return data, true
			}
		}
	}
	return DataWrapper{}, false
}
// TODO: optimize, but will probably disappear with storageBackend backend refactor...
func removeData(collection Collection, data DataWrapper) {
	var targetType NamedType
	done := false
	for namedType, dataWrappers := range collection.Data {
		for _, dataWrapper := range dataWrappers {
			if dataWrapper.Id == data.Id {
				targetType = namedType
				done = true
				break
			}
		}
		if done { break }
	}
	i := 0
	for _, elem := range collection.Data[targetType] {
		if elem.Id != data.Id {
			collection.Data[targetType][i] = elem
			i++
		}
	}
	collection.Data[targetType] = collection.Data[targetType][:i]
}

func nameToCollection(name string) (collection Collection, ok bool) {
	for _, elem := range collections {
		if elem.Name == name {
			collection = elem
			ok = true
			return
		}
	}
	ok = false
	return
}
func removeCollection(collection Collection) {
	i := 0
	for _, elem := range collections {
		if elem.Name != collection.Name {
			collections[i] = elem
			i++
		}
	}
	collections = collections[:i]
}
func addToCollection(dataWrapper DataWrapper, namedType NamedType, collection *Collection) {
	(*collection).Data[namedType] = append((*collection).Data[namedType], dataWrapper)
}
var collections = make([]Collection, 0)
func nameToNamedType(name string) (namedType NamedType, ok bool) {
	for _, elem := range namedTypes {
		if elem.Name == name {
			namedType = elem
			ok = true
			return
		}
	}
	ok = false
	return
}
func removeNamedType(namedType NamedType) {
	i := 0
	for _, elem := range namedTypes {
		// TODO: check only for name
		if elem != namedType {
			namedTypes[i] = elem
			i++
		}
	}
	namedTypes = namedTypes[:i]
}
var namedTypes = make([]NamedType, 0)
/// STORAGE END ///
