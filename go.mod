module git.freeself.one/thegergo02/easyt

go 1.17

require (
	git.freeself.one/thegergo02/easyt/body v0.0.0-00010101000000-000000000000
	git.freeself.one/thegergo02/easyt/storage v0.0.0-00010101000000-000000000000
	git.freeself.one/thegergo02/easyt/storage/backends/memory v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.7.7
	github.com/go-chi/chi v1.5.4
	github.com/stretchr/testify v1.7.0
	github.com/swaggest/rest v0.2.22
	github.com/swaggest/swgui v1.4.4
	github.com/swaggest/usecase v1.1.2
)

require (
	git.freeself.one/thegergo02/easyt/basic v0.0.0-00010101000000-000000000000 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-chi/chi/v5 v5.0.7 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v3 v3.1.0 // indirect
	github.com/swaggest/form/v5 v5.0.1 // indirect
	github.com/swaggest/jsonschema-go v0.3.24 // indirect
	github.com/swaggest/openapi-go v0.2.15 // indirect
	github.com/swaggest/refl v1.0.1 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	google.golang.org/protobuf v1.23.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	git.freeself.one/thegergo02/easyt/basic => ./basic
	git.freeself.one/thegergo02/easyt/body => ./body
	git.freeself.one/thegergo02/easyt/storage => ./storage
	git.freeself.one/thegergo02/easyt/storage/backends/memory => ./storage/backends/memory
)
