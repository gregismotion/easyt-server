/// MAIN START ///
package main

import (
	"git.freeself.one/thegergo02/easyt/basic"
	"git.freeself.one/thegergo02/easyt/storage"
	"git.freeself.one/thegergo02/easyt/body"
	"git.freeself.one/thegergo02/easyt/storage/backends/memory" // NOTE: temporary
	
	//"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter() (r *gin.Engine) {
	r = gin.Default()
	v1 := r.Group("/api/v1") 
	{
		col := v1.Group("/collections")
		{
			col.GET("/", getCollections)
			col.POST("/", createCollection)
			col.GET("/:id", getCollection)
			col.DELETE("/:id", deleteCollection)
			data := col.Group("/data")
			{
			      data.GET("/:colId/:dataId", getData)
			      data.POST("/:colId", addData)
			      data.DELETE("/:colId/:dataId", deleteData)
			}
		}
		typ := v1.Group("/types")
		{
			named := typ.Group("/named")
			{
				named.GET("/", getNamedTypes)
				named.POST("/", createNamedType)
				named.GET("/:id", getNamedType)
				named.DELETE("/:id", deleteNamedType)
			}
			typ.GET("/basic", getBasicTypes)
		}
	}
	return
}

var storageBackend storage.Storage
func main() {
	storageBackend = memory.New()
	host := "localhost:8080"
	setupRouter().Run(host)
}

type any interface{} // NOTE: remove in Go 1.18, default behaviour there
func respond(c *gin.Context, response any, err error, customSuccessStatus ...int) {
	if err == nil {
		status := http.StatusOK
		if customSuccessStatus != nil { status = customSuccessStatus[0] }
		c.IndentedJSON(status, response)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func getCollections(c *gin.Context) {
	references, err := storageBackend.GetCollectionReferences()
	respond(c, &references, err)
}

func createCollection(c *gin.Context) {
	var body body.CollectionRequestBody
	if err := c.BindJSON(&body); err == nil {
		reference, err := storageBackend.CreateCollectionByName(body.Name, body.NamedTypes)
		respond(c, &reference, err, http.StatusCreated)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func addData(c *gin.Context) {
	var body body.DataRequestBody
	id := c.Param("colId")
	if id != "" {
		if err := c.BindJSON(&body); err == nil {
			time := time.Now() // TODO: get time from body
			data, err := storageBackend.AddDataToCollectionById(body.NamedType, time, body.Value, id)
			respond(c, &data, err)
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func createNamedType(c *gin.Context) {
	var body body.NamedTypeRequestBody
	if err := c.BindJSON(&body); err == nil {
		namedType, err := storageBackend.CreateNamedType(body.Name, body.BasicType)
		respond(c, &namedType, err, http.StatusCreated)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
