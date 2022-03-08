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

type any interface{} // NOTE: remove in Go 1.18, default behaviour there
func respond(c *gin.Context, response any, err error) {
	if err == nil {
		c.IndentedJSON(http.StatusOK, response)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func getCollections(c *gin.Context) {
	references, err := storageBackend.GetCollectionReferences()
	respond(c, &references, err)
}

type CollectionRequestBody struct {
	Name 	     string `json:"name"`
	NamedTypes []string `json:"named_types"`
}
func createCollection(c *gin.Context) {
	var body CollectionRequestBody
	if err := c.BindJSON(&body); err == nil {
		reference, err := storageBackend.CreateCollectionByName(body.Name, body.NamedTypes)
		respond(c, &reference, err)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func getCollection(c *gin.Context) { // TODO: add return limit of data, maybe only send references to data
	id := c.Param("id")
	if id != "" {
		collection, err := storageBackend.GetCollectionById(id)
		respond(c, &collection, err)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No ID passed!"})
	}
}

func deleteCollection(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		respond(c, "ok", storageBackend.DeleteCollectionById(id))	
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
	id := c.Param("colId")
	if id != "" {
		if c.BindJSON(&body) != nil {
			time := time.Now() // TODO: get time from body
			data, err := storageBackend.AddDataToCollectionById(body.NamedType, time, body.Value, id)
			respond(c, &data, err)
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
		dataId := c.Param("dataId")
		if dataId != "" {
			data, err := storageBackend.GetDataInCollectionById(colId, dataId)
			respond(c, &data, err)
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No data ID specified!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No collection ID specified!"})
	}
}

func deleteData(c *gin.Context) {
	colId := c.Param("colId")
	if colId != "" {
		dataId := c.Param("dataId")
		if dataId != "" {
			respond(c, "ok", storageBackend.DeleteDataFromCollectionById(colId, dataId))
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No data ID specified!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}


func getNamedTypes(c *gin.Context) {
	namedTypes, err := storageBackend.GetNamedTypes()
	respond(c, &namedTypes, err)
}

type NamedTypeRequestBody struct {
	BasicType string `json:"basic_type"`
	Name string `json:"name"`
}
func createNamedType(c *gin.Context) {
	var body NamedTypeRequestBody
	if err := c.BindJSON(&body); err == nil && body.Name != "" {
		namedType, err := storageBackend.CreateNamedType(body.Name, body.BasicType)
		respond(c, &namedType, err)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request body!"})
	}
}

func getNamedType(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		namedType, err := storageBackend.GetNamedTypeById(id)
		respond(c, &namedType, err)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No ID specified!"})
	}
}

func deleteNamedType(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		respond(c, "ok", storageBackend.DeleteNamedTypeById(id))
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}

func getBasicTypes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, basic.GetBasicTypes())
}
