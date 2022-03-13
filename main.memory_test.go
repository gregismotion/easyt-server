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

func createTestCollection(name string) (*httptest.ResponseRecorder, string, error) {
	b, err := json.Marshal(body.CollectionRequestBody{Name: name})
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

func createTestDataPoints(colId string, dataPoints []body.DataRequestBody) (*httptest.ResponseRecorder, *storage.ReferenceGroups, error) {
	b, err := json.Marshal(dataPoints)
	if err == nil {
		w := makeRequest("POST", fmt.Sprintf("/api/v1/collections/data/%s", colId), bytes.NewReader(b))
		var resp storage.ReferenceGroups
		err = json.Unmarshal([]byte(w.Body.String()), &resp)
		return w, &resp, err
	} else {
		return nil, nil, err
	}
}
func deleteTestData(colId, groupId, dataId string) (*httptest.ResponseRecorder, error) {
	w := makeRequest("DELETE", fmt.Sprintf("/api/v1/collections/data/%s/%s/%s", colId, groupId, dataId))
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
	var w *httptest.ResponseRecorder
	var id string
	// create collection
	w, id, err = createTestCollection("body")
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

func TestData(t *testing.T) {
	storageBackend = getStorageBackend()
	// create collection
	namedIds := make([]string, 2)
	_, namedIds[0], _ = createTestNamedType("weight", "num")
	_, namedIds[1], _ = createTestNamedType("comment", "str")
	w, colId, err := createTestCollection("body")
	var dataGroup *storage.ReferenceGroups 
	group0 := []body.DataRequestBody { 
		body.DataRequestBody { NamedType: namedIds[0], Value: "5" },
		body.DataRequestBody { NamedType: namedIds[1], Value: "lean" },
	}
	group1 := []body.DataRequestBody { 
		body.DataRequestBody { NamedType: namedIds[0], Value: "10" },
		body.DataRequestBody { NamedType: namedIds[1], Value: "fat" },
	}
	w, dataGroup, err = createTestDataPoints(colId, group0)
	w, dataGroup, err = createTestDataPoints(colId, group1)
	assert.Equal(t, 201, w.Code)
	assert.Nil(t, err)
	assert.Greater(t, len(*dataGroup), 0)
	// check for a data
	var groupId string
	for k := range *dataGroup {
		groupId = k
		break
	}
	w = makeRequest("GET", fmt.Sprintf("/api/v1/collections/data/%s/%s/%s", colId, groupId, (*dataGroup)[groupId][0].Id))
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	b, _ := ioutil.ReadAll(w.Body)
	assert.True(t, strings.Contains(string(b), (*dataGroup)[groupId][0].Id))
	// delete a data
	w, err = deleteTestData(colId, groupId, (*dataGroup)[groupId][1].Id)
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	// check if deletion worked
	w = makeRequest("GET", fmt.Sprintf("/api/v1/collections/%s", colId))
	assert.Equal(t, 200, w.Code)
	b, _ = ioutil.ReadAll(w.Body)
	assert.False(t, strings.Contains(string(b), (*dataGroup)[groupId][1].Id))
}
