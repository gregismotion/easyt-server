package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"bytes"
	"encoding/json"
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

type NamedType struct {
	Name string `json:"name"`
	Type BasicType `json:"type"`
}
func isNamedTypeUnique(namedType NamedType, namedTypes []NamedType) bool {
	for _, elem := range namedTypes {
		if elem.Name == namedType.Name {
			return false
		}
	}
	return true
}
var namedTypes = make([]NamedType, 0)

func getNamedTypes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, namedTypes)
}
func createNamedType(c *gin.Context) {
	typ, ok := strToBasicType(c.PostForm("type"))
	if ok {
		name := c.PostForm("name")
		if name != "" {
			namedType := NamedType {
				Name: name,
				Type: typ,
			}
			ok = isNamedTypeUnique(namedType, namedTypes)
			if ok {
				namedTypes = append(namedTypes, namedType)
				c.IndentedJSON(http.StatusOK, namedType)
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Duplicate name!"})
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unknown basic type!"})
	}
}
func getNamedType(c *gin.Context) {}
func deleteNamedType(c *gin.Context) {}

type BasicType int
const (
	num BasicType = iota
	str
)
// NOTE: this will get out of hand, FIND ALTERNATIVE!!!
var strToBasicTypes = map[string]BasicType {
	"num": num,
	"str": str,
}
var basicTypesToStr = map[BasicType]string {
	num: "num",
	str: "str",
}
func (t BasicType) String() (str string) {
	str, ok := basicTypesToStr[t]
	if !ok {
		str = "unknown"
	}
	return
}
func strToBasicType(str string) (BasicType, bool) {
	typ, ok := strToBasicTypes[str]
	return typ, ok
}
func (typ BasicType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(typ.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}
func (typ *BasicType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	readTyp, ok := strToBasicType(j)
	if ok {
		*typ = readTyp
	} 	
	return nil
}

func getBasicTypes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, strToBasicTypes)
}
