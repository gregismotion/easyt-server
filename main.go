package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	
	v1 := r.Group("/api/v1") 
	{
		col := v1.Group("/collection")
		{
			col.GET("/", getCollections)
			col.POST("/", createCollection)
			col.GET("/:id", getCollection)
			col.POST("/:id", addToCollection)
			col.DELETE("/:id", deleteCollection)
			col.GET("/data/:idC/:idD", getData)
			col.DELETE("/data/:idC/:idD", deleteData)
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
	r.Run(host)
}

func getCollections(c *gin.Context) {}
func createCollection(c *gin.Context) {}
func getCollection(c *gin.Context) {}
func addToCollection(c *gin.Context) {}
func deleteCollection(c *gin.Context) {}
func getData(c *gin.Context) {}
func deleteData(c *gin.Context) {}
func getNamedTypes(c *gin.Context) {}
func createNamedType(c *gin.Context) {}
func getNamedType(c *gin.Context) {}
func deleteNamedType(c *gin.Context) {}

var basicTypes = []string{"num", "str"}
func getBasicTypes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, basicTypes)
}
