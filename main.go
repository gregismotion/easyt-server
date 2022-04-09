package main

import (
	"git.freeself.one/thegergo02/easyt/basic"
	"git.freeself.one/thegergo02/easyt/storage"

	//"git.freeself.one/thegergo02/easyt/body"
	"git.freeself.one/thegergo02/easyt/storage/backends/memory" // NOTE: temporary

	//"fmt"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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
		middleware.Recoverer,
		nethttp.OpenAPIMiddleware(apiSchema),
		request.DecoderMiddleware(decoder),
		request.ValidatorMiddleware(validator),
		response.EncoderMiddleware,
		gzip.Middleware,
	)

	r.Method(http.MethodGet, "/docs/openapi.json", apiSchema)
	r.Mount("/docs", v3cdn.NewHandler(apiSchema.Reflector().Spec.Info.Title,
		"/docs/openapi.json", "/docs"))

	r.Route("/types", func(r chi.Router) {
		r.Method(http.MethodGet, "/basic", nethttp.NewHandler(getBasicTypes()))
		r.Method(http.MethodGet, "/named", nethttp.NewHandler(getNamedTypes()))
		r.Method(http.MethodGet, "/named/{id}", nethttp.NewHandler(getNamedType()))
		r.Method(http.MethodPost, "/named", nethttp.NewHandler(createNamedType()))
		r.Method(http.MethodDelete, "/named/{id}", nethttp.NewHandler(deleteNamedType()))
	})

	r.Route("/collections", func(r chi.Router) {
		r.Method(http.MethodGet, "/", nethttp.NewHandler(getCollectionReferences()))
		r.Method(http.MethodPost, "/", nethttp.NewHandler(createCollection()))
		r.Method(http.MethodGet, "/{id}", nethttp.NewHandler(getCollection()))
		r.Method(http.MethodDelete, "/{id}", nethttp.NewHandler(deleteCollection()))
		r.Method(http.MethodPost, "/{id}", nethttp.NewHandler(addData()))
		r.Method(http.MethodGet, "/{colId}/{groupId}/{dataId}", nethttp.NewHandler(getData()))
		r.Method(http.MethodDelete, "/{colId}/{groupId}/{dataId}", nethttp.NewHandler(deleteData()))
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
		if namedTypes != nil {
			*out = *namedTypes
		}
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
			in  = input.(*getNamedTypeInput)
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
		Name      string `json:"name" required:"true"`
		BasicType string `json:"type" required:"true"`
	}
	u := usecase.NewIOI(new(createNamedTypeInput), new(storage.NamedType), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*createNamedTypeInput)
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
	u.SetTags("types")
	return u
}

func deleteNamedType() usecase.Interactor {
	type deleteNamedTypeInput struct {
		Id string `path:"id"`
	}
	u := usecase.NewIOI(new(deleteNamedTypeInput), nil, func(ctx context.Context, input, _ interface{}) error {
		var in = input.(*deleteNamedTypeInput)
		err := storageBackend.DeleteNamedTypeById(in.Id)
		if err != nil {
			return status.Wrap(err, status.NotFound)
		}
		return nil
	})
	u.SetExpectedErrors(status.NotFound)
	u.SetTags("types")
	return u
}

func getCollectionReferences() usecase.Interactor {
	type getCollectionReferencesInput struct {
		Id   string `query:"last_id" default:""`
		Size int    `query:"size" default:"10"`
	}
	u := usecase.NewIOI(new(getCollectionReferencesInput), new([]storage.NameReference), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*getCollectionReferencesInput)
			out = output.(*[]storage.NameReference)
		)
		references, err := storageBackend.GetCollectionReferences(in.Size, in.Id)
		if references != nil {
			*out = *references
		}
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
			in  = input.(*createCollectionInput)
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

func getCollection() usecase.Interactor { // TODO: add return limit of data
	type getCollectionInput struct {
		Id string `path:"id"`
	}
	u := usecase.NewIOI(new(getCollectionInput), new(storage.ReferenceCollection), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*getCollectionInput)
			out = output.(*storage.ReferenceCollection)
		)
		collection, err := storageBackend.GetReferenceCollectionById(in.Id)
		if err != nil {
			return status.Wrap(err, status.NotFound)
		}
		if collection != nil {
			*out = *collection
		}
		return err
	})
	u.SetExpectedErrors(status.NotFound)
	u.SetTags("collections")
	return u
}

func deleteCollection() usecase.Interactor {
	type deleteCollectionInput struct {
		Id string `path:"id"`
	}
	u := usecase.NewIOI(new(deleteCollectionInput), nil, func(ctx context.Context, input, _ interface{}) error {
		var in = input.(*deleteCollectionInput)
		err := storageBackend.DeleteCollectionById(in.Id)
		if err != nil {
			return status.Wrap(err, status.NotFound)
		}
		return nil
	})
	u.SetExpectedErrors(status.NotFound)
	u.SetTags("collections")
	return u
}

func addData() usecase.Interactor {
	type dataPointInput struct {
		NamedType string    `json:"named_type"`
		Time      time.Time `json:"time"`
		Value     string    `json:"value"`
	}
	type addDataInput struct {
		ColId           string           `path:"id"`
		DataPointInputs []dataPointInput `json:"data_points"`
	}
	u := usecase.NewIOI(new(addDataInput), new(storage.ReferenceGroups), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*addDataInput)
			out = output.(*storage.ReferenceGroups)
		)
		dataPoints := make([]storage.DataPoint, len(in.DataPointInputs))
		for i, dataPointInput := range in.DataPointInputs {
			dataPoints[i] = storage.DataPoint{
				NamedType: storage.NamedType{Id: dataPointInput.NamedType},
				Time:      dataPointInput.Time,
				Value:     dataPointInput.Value}
		}
		referenceGroups, err := storageBackend.AddDataPointsToCollectionById(in.ColId, dataPoints)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		*out = *referenceGroups
		return nil
	})
	u.SetTags("data")
	return u
}

func getData() usecase.Interactor {
	type getDataInput struct {
		ColId   string `path:"colId"`
		GroupId string `path:"groupId"`
		DataId  string `path:"dataId"`
	}
	u := usecase.NewIOI(new(getDataInput), new(storage.DataPoint), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*getDataInput)
			out = output.(*storage.DataPoint)
		)
		data, err := storageBackend.GetDataInCollectionById(in.ColId, in.GroupId, in.DataId)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		*out = *data
		return nil
	})
	u.SetTags("data")
	return u
}

func deleteData() usecase.Interactor {
	type deleteDataInput struct {
		ColId   string `path:"colId"`
		GroupId string `path:"groupId"`
		DataId  string `path:"dataId"`
	}
	u := usecase.NewIOI(new(deleteDataInput), nil, func(ctx context.Context, input, _ interface{}) error {
		var in = input.(*deleteDataInput)
		err := storageBackend.DeleteDataFromCollectionById(in.ColId, in.GroupId, in.DataId)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetTags("data")
	return u
}
