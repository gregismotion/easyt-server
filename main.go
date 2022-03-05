/// MAIN START ///
package main

// TODO: move these to right places
import (
	"git.freeself.one/thegergo02/easyt/basic"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fmt"
	"net/http"
	"time"
	"strconv"

	"bytes"
)

func main() {
	r := gin.Default()
	
	v1 := r.Group("/api/v1") 
	{
		col := v1.Group("/collection")
		{
			// TODO: use id for collections too, names can get weird in urls...
			col.GET("/", getCollections)
			col.POST("/", createCollection)
			col.GET("/:name", getCollection)
			col.DELETE("/:name", deleteCollection)
			col.POST("/:name", addData)
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

func getCollections(c *gin.Context) {
	collectionNames := make([]string, 0)
	for _, collection := range collections {
		collectionNames = append(collectionNames, collection.Name)
	}
	c.IndentedJSON(http.StatusOK, collectionNames)
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

func deleteCollection(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		collection, ok := nameToCollection(name)
		if ok {
			removeCollection(collection)
			// NOTE: maybe some message would be appropiate? consult the do- oh wait
			c.String(http.StatusOK, "")
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find collection!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}


type DataRequestBody struct {
	Time 	  string `json:"time,omitempty"`
	NamedType string `json:"named_type"`
	Value 	  string `json:"value"`
}
func addData(c *gin.Context) {
	var body DataRequestBody
	name := c.Param("name")
	if name != "" {
		if err := c.BindJSON(&body); err == nil  {
			collection, ok := nameToCollection(name)
			if ok {
				time := time.Now() // TODO: get time from body
				namedType, okTyp := nameToNamedType(body.NamedType)
				if okTyp {
					dataWrapper := DataWrapper {
						Id: uuid.New().String(),
						Time: time,
						Type: namedType.Type,
					}
					value := body.Value
					// TODO: should return error at unparseable values
					switch namedType.Type {
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
					addToCollection(dataWrapper, namedType, &collection)
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}

func getData(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		collection, ok := nameToCollection(name)
		if ok {
			id := c.Param("idD")
			if id != "" {
				data, ok := idToData(collection, id)
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
	}
}

func deleteData(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		collection, ok := nameToCollection(name)
		if ok {
			id := c.Param("idD")
			if id != "" {
				data, ok := idToData(collection, id)
				if ok {
					removeData(collection, data)
					// TODO: return smt?
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
	c.IndentedJSON(http.StatusOK, namedTypes)
}

// TODO: use json
func createNamedType(c *gin.Context) {
	typ, ok := basic.StrToBasicType(c.PostForm("type"))
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

func getBasicTypes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, basic.StrToBasicTypes)
}
/// MAIN END ///



/// STORAGE START ///
// NOTE: might become unmanageable, find alternative
// TODO: nicer string formatting, this looks ugly rn
func (data DataWrapper) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`{"time":"`)
	buffer.WriteString(data.Time.String())
	buffer.WriteString(`","type":"`)
	buffer.WriteString(data.Type.String())
	buffer.WriteString(`","id":"`)
	buffer.WriteString(data.Id)
	buffer.WriteString(`","value":"`)
	switch data.Type {
		case basic.Num:
			buffer.WriteString(fmt.Sprintf("%.5f", data.Num)) 
			// TODO: what precision do we need?
		case basic.Str:
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
		dFirst := true
		for _, dataWrapper := range dataWrappers {
			if !dFirst { buffer.WriteString(`,`) } else { dFirst = false }
			bytes, err := dataWrapper.MarshalJSON()
			if err != nil { return buffer.Bytes(), err }
			buffer.Write(bytes)
		}
		buffer.WriteString(`]`)
	}
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}
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
// TODO: optimize, but will probably disappear with storage backend refactor...
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

type DataWrappers map[NamedType][]DataWrapper
type DataWrapper struct {
	Id string `json:"id"`
	// TODO: try tinytime, we don't need nanosecond precision...
	Time time.Time `json:"time"`
	Type basic.BasicType `json:"type"`
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


type NamedType struct {
	Name string `json:"name"`
	Type basic.BasicType `json:"type"`
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
