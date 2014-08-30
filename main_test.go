package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseRecorder(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("GET", "http://foo/bar", nil)
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, w.Code, 500)

	expected := fmt.Sprint("something failed\n")
	assert.Equal(t, w.Body.String(), expected)

	assert.Equal(t, w.HeaderMap["Content-Type"][0], "text/plain; charset=utf-8")
}

func TestNotFound(t *testing.T) {
	notFound := func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}

	req, err := http.NewRequest("GET", "bla", nil)
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	notFound(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestHeaderLocation(t *testing.T) {
	createdHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/users/123")
		w.WriteHeader(http.StatusCreated)
	}

	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	createdHandler(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, w.HeaderMap["Location"][0], "/users/123")
}

func TestGoHandlerSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/hello?message=gorules", nil)
	w := httptest.NewRecorder()

	GoHandler(w, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Body.String(), "RIGHT")
}

func TestGoHandlerFailure(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/hello?message=gosucks", nil)
	w := httptest.NewRecorder()

	GoHandler(w, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Equal(t, w.Body.String(), "Error\n")
}

func TestServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "HelloWorld")
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, string(greeting), "HelloWorld")
}
