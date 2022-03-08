/// MAIN START ///
package main

import (
	"git.freeself.one/thegergo02/easyt/basic"
	"git.freeself.one/thegergo02/easyt/storage"
	"git.freeself.one/thegergo02/easyt/storage/backends/memory" // NOTE: temporary
	
	//"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var storageBackend storage.Storage
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

	storageBackend = memory.New()
	r.Run(host)
}

func getCollections(c *gin.Context) {
	references, ok := storageBackend.GetCollectionReferences()
	if ok {
		c.IndentedJSON(http.StatusOK, references)
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"error": "Couldn't get collection references!"})
	}
}

type CollectionRequestBody struct {
	Name 	     string `json:"name"`
	NamedTypes []string `json:"named_types"`
}
func createCollection(c *gin.Context) {
	var body CollectionRequestBody
	if err := c.BindJSON(&body); err == nil {
			if reference, ok := storageBackend.CreateCollectionByName(body.Name, body.NamedTypes); ok {
				c.IndentedJSON(http.StatusOK, reference)
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Couldn't create collection!"})
			}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.String()})
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
		if c.BindJSON(&body) != nil {
			time := time.Now() // TODO: get time from body
			if data, ok := storageBackend.AddDataToCollectionById(body.NamedType, time, body.Value, id); ok {
				c.IndentedJSON(http.StatusCreated, data)
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Couldn't create data!"})
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
				if storageBackend.DeleteDataFromCollectionById(colId, dataId) {
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
	namedTypes, ok := storageBackend.GetNamedTypes()
	if ok {
		c.IndentedJSON(http.StatusOK, namedTypes)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Couldn't get named types!"})
	}
}

type NamedTypeRequestBody struct {
	BasicType string `json:"basic_type"`
	Name string `json:"name"`
}
func createNamedType(c *gin.Context) {
	var body NamedTypeRequestBody
	if err := c.BindJSON(&body); err == nil {
		if namedType, ok := storageBackend.CreateNamedType(body.Name, body.BasicType); ok {
			c.IndentedJSON(http.StatusOK, namedType)
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Couldn't create named type!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request body!"})
	}
}

func getNamedType(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		namedType, ok := storageBackend.GetNamedTypeById(id)
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
	if id != "" {
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
