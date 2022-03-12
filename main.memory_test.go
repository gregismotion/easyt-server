package main

import (
	"git.freeself.one/thegergo02/easyt/body"
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
	"bytes"
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

func createTestCollection(name string, namedTypes []string) (*httptest.ResponseRecorder, string, error) {
	b, err := json.Marshal(body.CollectionRequestBody{Name: name, NamedTypes: namedTypes})
	if err == nil {
		w := makeRequest( "POST", "/api/v1/collections/", bytes.NewReader(b))
		var resp map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &resp)
		return w, resp["id"], err
	} else {
		return nil, "", err
	}
}
func deleteTestCollection(id string) (*httptest.ResponseRecorder, error) {
	w := makeRequest("DELETE", fmt.Sprintf("/api/v1/collections/%s", id))
	return w, nil
}

func createTestNamedType(name string, basicType string) (*httptest.ResponseRecorder, string, error) {
	b, err := json.Marshal(body.NamedTypeRequestBody{Name: name, BasicType: basicType})
	if err == nil {
		w := makeRequest( "POST", "/api/v1/types/named/", bytes.NewReader(b))
		var resp map[string]string
		err = json.Unmarshal([]byte(w.Body.String()), &resp)
		return w, resp["id"], err
	} else {
		return nil, "", err
	}
}
func deleteTestNamedType(id string) (*httptest.ResponseRecorder, error) {
	w := makeRequest("DELETE", fmt.Sprintf("/api/v1/types/named/%s", id))
	return w, nil
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
	w := makeRequest( "GET", "/api/v1/collections/")
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", string(b))
}
func TestGetNamedTypes(t *testing.T) {
	w := makeRequest( "GET", "/api/v1/types/named/")
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", string(b))
}
func TestGetBasicTypes(t *testing.T) {
	w := makeRequest( "GET", "/api/v1/types/basic")
	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, 200, w.Code)
	assert.True(t, strings.Contains(string(b), "num"))
	assert.True(t, strings.Contains(string(b), "str"))
}

func TestNamedType(t *testing.T) {
	// create named type
	w, id, err := createTestNamedType(testResponse, "num")
	assert.Nil(t, err)
	assert.Equal(t, 201, w.Code)
	// get named type
	w = makeRequest("GET", fmt.Sprintf("/api/v1/types/named/%s", id))
	assert.Equal(t, 200, w.Code)
	b, _ := ioutil.ReadAll(w.Body)
	assert.True(t, strings.Contains(string(b), testResponse))
	// delete named type
	w, err = deleteTestNamedType(id)
	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	// check if named type got deleted
	w = makeRequest("GET", fmt.Sprintf("/api/v1/types/named/"))
	assert.Equal(t, 200, w.Code)
	b, _ = ioutil.ReadAll(w.Body)
	assert.False(t, strings.Contains(string(b), testResponse))
}
func TestCollection(t *testing.T) {
	var err error
	storageBackend = getStorageBackend()
	// have to create named types
	ids := make([]string, 3)
	_, ids[0], err = createTestNamedType("weight", "num")
	_, ids[1], err = createTestNamedType("height", "num")
	_, ids[2], err = createTestNamedType("comment", "str")
	assert.Nil(t, err, "Failed to create named type(s)!")
	var w *httptest.ResponseRecorder
	var id string
	w, id, err = createTestCollection("body", ids)
	assert.Nil(t, err, "Failed to create collection")
	assert.Equal(t, 201, w.Code)
	// check if collection got created
	w = makeRequest("GET", fmt.Sprintf("/api/v1/collections/%s", id))
	assert.Equal(t, 200, w.Code)
	b, _ := ioutil.ReadAll(w.Body)
	assert.True(t, strings.Contains(string(b), "body"))
	// delete collection
	w, err = deleteTestCollection(id)
	assert.Equal(t, 200, w.Code)
	// check if collection got deleted
	w = makeRequest("GET", "/api/v1/collections/")
	b, _ = ioutil.ReadAll(w.Body)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", string(b))
}
