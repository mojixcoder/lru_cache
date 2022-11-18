package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	r := newRouter()

	testcases := []struct {
		name   string
		reqUrl string
		method string
	}{
		{name: "get", reqUrl: "/get/10", method: http.MethodGet},
		{name: "set", reqUrl: "/set", method: http.MethodPost},
		{name: "flush", reqUrl: "/flush", method: http.MethodGet},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			body := bytes.NewReader([]byte(``))

			req, err := http.NewRequest(tc.method, tc.reqUrl, body)
			assert.NoError(t, err)

			r.ServeHTTP(rr, req)

			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			_, err = ioutil.ReadAll(rr.Body)

			assert.NoError(t, err)
		})
	}

}

func TestRunServer(t *testing.T) {
	go RunServer()
	time.Sleep(10 * time.Millisecond)

	testcases := []struct {
		name   string
		reqUrl string
		method string
	}{
		{name: "get", reqUrl: "http://localhost:2376/get/10", method: http.MethodGet},
		{name: "set", reqUrl: "http://localhost:2376/set", method: http.MethodPost},
		{name: "get", reqUrl: "http://localhost:2376/flush", method: http.MethodGet},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			body := bytes.NewReader([]byte(``))

			req, err := http.NewRequest(tc.method, tc.reqUrl, body)
			assert.NoError(t, err)

			res, err := http.DefaultClient.Do(req)

			assert.NoError(t, err)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
		})
	}

	os.Setenv("SERVER_WRITE_TIMEOUT", "invalid_value")

	assert.Panics(t, func() {
		RunServer()
	})

	os.Setenv("SERVER_WRITE_TIMEOUT", "1s")
}
