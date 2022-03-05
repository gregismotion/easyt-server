package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"bytes"
	"encoding/json"
	"time"
)

func main() {
	r := gin.Default()
	
	v1 := r.Group("/api/v1") 
	{
		col := v1.Group("/collection")
		{
			col.GET("/", getCollections)
			col.POST("/", createCollection)
			col.GET("/:name", getCollection)
			col.POST("/:name", addToCollection)
			col.DELETE("/:name", deleteCollection)
			col.GET("/data/:name/:idD", getData)
			col.DELETE("/data/:name/:idD", deleteData)
		}
		typ := v1.Group("/type")
		{
			typ.GET("/named", getNamedTypes)
			typ.POST("/named", createNamedType)
			typ.GET("/:name", getNamedType)
			typ.DELETE("/:name", deleteNamedType)
			typ.GET("/basic", getBasicTypes)
		}
	}

	host := "localhost:8080"
	r.Run(host)
}

// NOTE: might become unmanageable, find alternative
// TODO: nicer string formatting, this looks ugly rn
func (data DataWrapper) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`{"time":"`)
	buffer.WriteString(data.Time.String())
	buffer.WriteString(`","type":"`)
	buffer.WriteString(data.Type.String())
	buffer.WriteString(`","value":"`)
	switch data.Type {
		case num:
			buffer.WriteString(fmt.Sprintf("%.5f", data.Num)) 
			// TODO: what precision do we need?
		case str:
			buffer.WriteString(data.Str)
		default:
			buffer.WriteString("unknown")

	}
	buffer.WriteString(`"}`)
	return buffer.Bytes(), nil
}
func (data DataWrappers) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	first := true
	for namedType, dataWrappers := range data {
		if !first { buffer.WriteString(`,`) } else { first = false }
		buffer.WriteString(`"`)
		buffer.WriteString(namedType.Name)
		buffer.WriteString(`":[`)
		for _, dataWrapper := range dataWrappers {
			bytes, err := dataWrapper.MarshalJSON()
			if err != nil { return buffer.Bytes(), err }
			buffer.Write(bytes)
			buffer.WriteString(`,`)
		}
		buffer.WriteString(`]`)
	}
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}
type DataWrappers map[NamedType][]DataWrapper
type DataWrapper struct {
	// TODO: try tinytime, we don't need nanosecond precision...
	Time time.Time `json:"time"`
	Type BasicType `json:"type"`
	Num float64 `json:"num"`
	Str string `json:"str"`
}
type Collection struct {
	Name string `json:"name"`
	Data DataWrappers `json:"type"`
}
func (collection Collection) isUnique() bool {
	for _, elem := range collections {
		if elem.Name == collection.Name {
			return false
		}
	}
	return true
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
var collections = make([]Collection, 0)
func getCollections(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, collections)
}
type CollectionRequestBody struct {
	Name string `json:"name"`
	NamedTypes []string `json:"named_types"`
}
func createCollection(c *gin.Context) {
	var body CollectionRequestBody
	if err := c.BindJSON(&body); err == nil  {
		collection := Collection {
			Name: body.Name,
			Data: make(DataWrappers),
		}
		if collection.isUnique() {
			for _, name := range body.NamedTypes {
				namedType, ok := nameToNamedType(name)
				if ok {
					collection.Data[namedType] = make([]DataWrapper, 0)
				} else {
					// TODO: completely fail, ignore or smt else when bad named type?
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Non-existent named type!"})
					return
				}
			}
			collections = append(collections, collection)
			c.IndentedJSON(http.StatusOK, collection)
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Duplicate name!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request body!"})
	}
}
func getCollection(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		collection, ok := nameToCollection(name)
		if ok {
			c.IndentedJSON(http.StatusOK, collection)
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Couldn't find collection with this name!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name passed!"})
	}
}
func addToCollection(c *gin.Context) {}
func deleteCollection(c *gin.Context) {}
func getData(c *gin.Context) {}
func deleteData(c *gin.Context) {}

type NamedType struct {
	Name string `json:"name"`
	Type BasicType `json:"type"`
}
func (namedType NamedType) isUnique() bool {
	for _, elem := range namedTypes {
		if elem.Name == namedType.Name {
			return false
		}
	}
	return true
}
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
		if elem != namedType {
			namedTypes[i] = elem
			i++
		}
	}
	namedTypes = namedTypes[:i]
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
			if namedType.isUnique() {
				namedTypes = append(namedTypes, namedType)
				c.IndentedJSON(http.StatusCreated, namedType)
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
func getNamedType(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		namedType, ok := nameToNamedType(name)
		if ok {
			c.IndentedJSON(http.StatusOK, namedType)
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find named type!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}
func deleteNamedType(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		namedType, ok := nameToNamedType(name)
		if ok {
			removeNamedType(namedType)
			// NOTE: maybe some message would be appropiate? consult the do- oh wait
			c.String(http.StatusOK, "")
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find named type!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}

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
