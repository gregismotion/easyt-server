package main

import (
	"git.freeself.one/thegergo02/easyt/storage"
	"git.freeself.one/thegergo02/easyt/storage/backends/memory"
	
	"io"
	"testing"
	"net/http"
	"net/http/httptest"
	"errors"
	"io/ioutil"
	"fmt"
	"strings"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getStorageBackend() storage.Storage {
	return memory.New()
}
func createTestContext() (w *httptest.ResponseRecorder, c *gin.Context) {
	w = httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ = gin.CreateTestContext(w)
	return
}

func makeRequest(method string, path string, body ...io.Reader) (w *httptest.ResponseRecorder) {
	var req *http.Request
	w = httptest.NewRecorder()
	if len(body) > 0 {
		req, _ = http.NewRequest(method, path, body[0])
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	r.ServeHTTP(w, req)
	return
}

func createTestCollection(name string, id0 string, id1 string, id2 string) (*httptest.ResponseRecorder, string, error) {
	w := makeRequest( "POST", "/api/v1/collection/", strings.NewReader(fmt.Sprintf("{\"name\":%q,\"named_types\":[%q, %q, %q]}", name, id0, id1, id2)))
	var resp map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	return w, resp["id"], err
}
func deleteTestCollection(id string) (*httptest.ResponseRecorder, error) {
	w := makeRequest("DELETE", fmt.Sprintf("/api/v1/collection/%s", id))
	return w, nil
}

func createTestNamedType(name string, basicType string) (*httptest.ResponseRecorder, string, error) {
	w := makeRequest( "POST", "/api/v1/type/named", strings.NewReader(fmt.Sprintf("{\"name\":%q,\"basic_type\":%q}", name, basicType)))
	var resp map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	return w, resp["id"], err
}


var r = setupRouter()

var testResponse = "lightning mcqueen"
var testResponse1 = "test1"
var testResponse2 = "test2"
var ErrTest = errors.New("stanley") 
func TestRespond(t *testing.T) {
	w, c := createTestContext()
	respond(c, testResponse, nil)
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, fmt.Sprintf("%q", testResponse), string(b))
}
func TestRespondError(t *testing.T) {
	w, c := createTestContext()
	respond(c, testResponse, ErrTest)
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, w.Code, 500)
	assert.NotEqual(t, fmt.Sprintf("%q", testResponse), string(b))
	assert.True(t, strings.Contains(string(b), ErrTest.Error()))
}

func TestGetCollections(t *testing.T) {
	storageBackend = getStorageBackend()
	w := makeRequest( "GET", "/api/v1/collection/")
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", string(b))
}

func TestCollection(t *testing.T) { // Try somehow decoupling tests
	storageBackend = getStorageBackend()
	// have to create named types
	_, id0, err0 := createTestNamedType("weight", "num")
	_, id1, err1 := createTestNamedType("height", "num")
	_, id2, err2 := createTestNamedType("comment", "str")
	assert.Nil(t, err0)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	w, id, err := createTestCollection("body", id0, id1, id2)
	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	// check if collection got created
	w = makeRequest( "GET", fmt.Sprintf("/api/v1/collection/%s", id))
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, 200, w.Code)
	assert.True(t, strings.Contains(string(b), "body"))
	// delete collection
	w, err = deleteTestCollection(id)
	assert.Equal(t, 200, w.Code)
	// check if collection got created
	w = makeRequest( "GET", "/api/v1/collection/")
	b, _ = ioutil.ReadAll(w.Body)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", string(b))
}
