package main

import (
	"git.freeself.one/thegergo02/easyt/storage"
	"git.freeself.one/thegergo02/easyt/storage/backends/memory"

	"testing"
	"net/http"
	"net/http/httptest"
	"errors"
	"io/ioutil"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createTestStorage() storage.Storage {
	return memory.New()
}
func createTestContext() (w *httptest.ResponseRecorder, c *gin.Context) {
	w = httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ = gin.CreateTestContext(w)
	return
}
func createTestEnv() (*gin.Engine, *httptest.ResponseRecorder, storage.Storage) {
	return setupRouter(), httptest.NewRecorder(), createTestStorage()

}
func makeRequest(r *gin.Engine, w *httptest.ResponseRecorder, method string, path string) {
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
}

var testResponse = "lightning mcqueen"
var ErrTest = errors.New("stanley") 

func TestRespond(t *testing.T) {
	w, c := createTestContext()
	respond(c, testResponse, nil)
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, string(b), fmt.Sprintf("%q", testResponse))
}
func TestRespondError(t *testing.T) {
	w, c := createTestContext()
	respond(c, testResponse, ErrTest)
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, w.Code, 500)
	assert.NotEqual(t, string(b), fmt.Sprintf("%q", testResponse))
}

func TestGetCollections(t *testing.T) {
	var r *gin.Engine
	var w *httptest.ResponseRecorder
	r, w, storageBackend = createTestEnv()
	makeRequest(r, w, "GET", "/api/v1/collection/")
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, string(b), "[]")
}

/*func TestCreateCollection(t *testing.T) {
	var r *gin.Engine
	var w *httptest.ResponseRecorder
	r, w, storageBackend = createTestEnv()
	makeRequest(r, w, "GET", "/api/v1/collection/")
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, string(b), "[]")
}*/
