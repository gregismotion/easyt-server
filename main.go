package main

import (
	"git.freeself.one/thegergo02/easyt/basic"
	"git.freeself.one/thegergo02/easyt/storage"
	//"git.freeself.one/thegergo02/easyt/body"
	"git.freeself.one/thegergo02/easyt/storage/backends/memory" // NOTE: temporary
	
	//"fmt"
	"net/http"
	//"time"
	"log"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/middleware"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/chirouter"
	"github.com/swaggest/rest/jsonschema"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/openapi"
	"github.com/swaggest/rest/request"
	"github.com/swaggest/rest/response"
	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/swgui/v3cdn"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

/*func setupRouter() (r *gin.Engine) {
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
				data.GET("/:colId/:groupId/:dataId", getData)
				data.POST("/:colId", addData)
				data.DELETE("/:colId/:groupId/:dataId", deleteData)
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
}*/

func startRouter(host string, r *chirouter.Wrapper) error { 
	log.Printf("Server started on %s", host) // BUG: gets printed even if error
	return http.ListenAndServe(host, *r) 
}

func setupApiSchema() (apiSchema *openapi.Collector) {
	apiSchema = new(openapi.Collector)
	apiSchema.Reflector().SpecEns().Info.Title = "EasyTracker"
	apiSchema.Reflector().SpecEns().Info.WithDescription("A service (with a REST API) to create data-points in an organized manner.")
	apiSchema.Reflector().SpecEns().Info.Version = "v0.1.0"
	return
}
func setupValidator(apiSchema *openapi.Collector) (validator jsonschema.Factory) {
	validator = jsonschema.NewFactory(apiSchema, apiSchema)
	return
}
func setupDecoder() (decoder *request.DecoderFactory) {
	decoder = request.NewDecoderFactory()
	decoder.ApplyDefaults = true
	decoder.SetDecoderFunc(rest.ParamInPath, chirouter.PathToURLValues)
	return
}
func setupRouter() (r *chirouter.Wrapper) {
	r = chirouter.NewWrapper(chi.NewRouter())

	apiSchema := setupApiSchema()
	validator := setupValidator(apiSchema)
	decoder := setupDecoder()
	r.Use(
		middleware.Recoverer,                          // Panic recovery.
		nethttp.OpenAPIMiddleware(apiSchema),          // Documentation collector.
		request.DecoderMiddleware(decoder),     // Request decoder setup.
		request.ValidatorMiddleware(validator),		// Request validator setup.
		response.EncoderMiddleware,                    	// Response encoder setup.
		gzip.Middleware,                               // Response compression with support for direct gzip pass through.
	)

	r.Method(http.MethodGet, "/docs/openapi.json", apiSchema)
	r.Mount("/docs", v3cdn.NewHandler(apiSchema.Reflector().Spec.Info.Title,
	"/docs/openapi.json", "/docs"))
	
	r.Route("/types", func (r chi.Router) {
		r.Method(http.MethodGet, "/basic", nethttp.NewHandler(getBasicTypes()))
		r.Method(http.MethodGet, "/named", nethttp.NewHandler(getNamedTypes()))
		r.Method(http.MethodGet, "/named/{id}", nethttp.NewHandler(getNamedType()))
		r.Method(http.MethodPost, "/named", nethttp.NewHandler(createNamedType()))
	})

	r.Route("/collections", func (r chi.Router) {
		r.Method(http.MethodGet, "/", nethttp.NewHandler(getCollectionReferences()))
		r.Method(http.MethodPost, "/", nethttp.NewHandler(createCollection()))
	})
	return
}


var storageBackend storage.Storage
func main() {
	storageBackend = memory.New()

	r := setupRouter()

	host := "localhost:8080"
	if err := startRouter(host, r); err != nil {
		log.Fatal(err) 
	}
}

func getBasicTypes() usecase.Interactor {
	u := usecase.NewIOI(nil, new([]string), func(ctx context.Context, _, output interface{}) error {
		var out = output.(*[]string)
		*out = basic.GetBasicTypes()
		return nil
	})
	u.SetTags("types")
	return u
}

func getNamedTypes() usecase.Interactor {
	u := usecase.NewIOI(nil, new([]storage.NamedType), func(ctx context.Context, _, output interface{}) error {
		var out = output.(*[]storage.NamedType)
		namedTypes, err := storageBackend.GetNamedTypes()
		if namedTypes != nil { *out = *namedTypes }
		return err
	})
	u.SetTags("types")
	return u
}

func getNamedType() usecase.Interactor {
	type getNamedTypeInput struct {
		Id string `path:"id"`
	}
	u := usecase.NewIOI(new(getNamedTypeInput), new(storage.NamedType), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*getNamedTypeInput)
			out = output.(*storage.NamedType)
		)
		namedType, err := storageBackend.GetNamedTypeById(in.Id)
		if err != nil {
			return status.Wrap(err, status.NotFound)
		}
		*out = *namedType
		return nil
	})
	u.SetExpectedErrors(status.NotFound)
	u.SetTags("types")
	return u
}
func createNamedType() usecase.Interactor {
	type createNamedTypeInput struct {
		Name string `json:"name" required:"true"`
		BasicType string `json:"type" required:"true"`
	}
	u := usecase.NewIOI(new(createNamedTypeInput), new(storage.NamedType), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*createNamedTypeInput)
			out = output.(*storage.NamedType)
		)
		namedType, err := storageBackend.CreateNamedType(in.Name, in.BasicType)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		*out = *namedType
		return nil
	})
	u.SetExpectedErrors(status.Internal)
	u.SetTags("collections")
	return u
}

func getCollectionReferences() usecase.Interactor {
	u := usecase.NewIOI(nil, new([]storage.NameReference), func(ctx context.Context, _, output interface{}) error {
		var out = output.(*[]storage.NameReference)
		references, err := storageBackend.GetCollectionReferences()
		if references != nil { *out = *references }
		return err
	})
	u.SetTags("collections")
	return u
}

func createCollection() usecase.Interactor {
	type createCollectionInput struct {
		Name string `json:"name" required:"true"`
	}
	u := usecase.NewIOI(new(createCollectionInput), new(storage.NameReference), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*createCollectionInput)
			out = output.(*storage.NameReference)
		)
		reference, err := storageBackend.CreateCollectionByName(in.Name)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		*out = *reference
		return nil
	})
	u.SetExpectedErrors(status.Internal)
	u.SetTags("collections")
	return u
}

/*
func getCollection(c *gin.Context) { // TODO: add return limit of data
	id := c.Param("id")
	if id != "" {
		collection, err := storageBackend.GetReferenceCollectionById(id)
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
	var dataPoints []body.DataRequestBody
	id := c.Param("colId")
	if id != "" {
		if err := c.BindJSON(&dataPoints); err == nil {
			//time := time.Now() // TODO: get time from body
			namedTypeIds := make([]string, len(dataPoints))
			values := make([]string, len(dataPoints))
			for i, dataPoint := range dataPoints { // TODO: move this to another func
				namedTypeIds[i] = dataPoint.NamedType
				values[i] = dataPoint.Value
			}
			data, groupId, err := storageBackend.AddDataPointsToCollectionById(id, namedTypeIds, values)
			var dataGroup storage.ReferenceGroups
			if data != nil {
				dataGroup = storage.ReferenceGroups {
					groupId: *data,
				}
			}
			respond(c, dataGroup, err, http.StatusCreated)
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
			groupId := c.Param("groupId")
			if groupId != "" {
				data, err := storageBackend.GetDataInCollectionById(colId, groupId, dataId)
				respond(c, &data, err)
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No group ID specified!"})
			}
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
			groupId := c.Param("groupId")
			if groupId != "" {
				respond(c, "ok", storageBackend.DeleteDataFromCollectionById(colId, groupId, dataId))
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No group ID specified!"})
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No data ID specified!"})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No name specified!"})
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
*/
