package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"

	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestRouteSimple(t *testing.T) {
	mux := http.NewServeMux()

	globalMiddle := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Executing globalmiddle")
			next.ServeHTTP(w, r)
		})
	}

	var app = New(mux, globalMiddle)

	testHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Executing handler")
			Respond(next, w, r, "hi", http.StatusOK)
			//next.ServeHTTP(w, r)
		})
	}

	testMiddlware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Executing middlewareOne")
			next.ServeHTTP(w, r)
			//log.Println("m1 after")
		})
	}

	testMiddlware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Executing middlewareTwo")
			next.ServeHTTP(w, r)
			//log.Println("m2 after")
		})
	}

	app.Handle(GET, "/test", testHandler, testMiddlware, testMiddlware2)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	app.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")
	b, _ := ioutil.ReadAll(w.Body)
	fmt.Println("BODY", string(b))
}

//
//func TestRouteSimpleWithResponse(t *testing.T) {
//	l := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
//	var app = New(l, nil)
//
//	testHandler := func(log *log.Logger, w http.ResponseWriter, r *http.Request) error {
//		a := make(map[string]string)
//		a["key"] = "value"
//		Respond(log, w, a, http.StatusOK)
//		return nil
//	}
//
//	app.Handle(http.MethodGet, "/test", testHandler)
//
//	r := httptest.NewRequest(http.MethodGet, "/test", nil)
//	w := httptest.NewRecorder()
//
//	app.ServeHTTP(w, r)
//	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")
//	b, err := ioutil.ReadAll(w.Body)
//	assert.NoError(t, err)
//	assert.JSONEq(t, `{"key": "value"}`, string(b))
//}

