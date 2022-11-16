package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockReader struct {
}

func (r mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read body error")
}

func setToCache(t *testing.T, r *mux.Router, key string, value any) {
	rr := httptest.NewRecorder()

	jsonBody := []byte(fmt.Sprintf(`{"key": "%s", "value": %v}`, key, value))
	reader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, "/set", reader)
	assert.NoError(t, err)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGet(t *testing.T) {
	r := newRouter()

	setToCache(t, r, "10", 10)

	testcases := []struct {
		name       string
		statusCode int
		reqUrl     string
		resp       []byte
	}{
		{name: "ok", statusCode: http.StatusOK, reqUrl: "/get/10", resp: []byte(`{"key":"10","value":10}`)},
		{name: "not_found", statusCode: http.StatusNotFound, reqUrl: "/get/not_found", resp: []byte(`{"detail": "not found"}`)},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tc.reqUrl, nil)
			assert.NoError(t, err)

			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.statusCode, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			resp, err := ioutil.ReadAll(rr.Body)

			assert.NoError(t, err)
			assert.Equal(t, tc.resp, resp)
		})
	}

}

func TestSet(t *testing.T) {
	r := newRouter()

	testcases := []struct {
		name       string
		statusCode int
		resp       []byte
		body       []byte
	}{
		{name: "ok", statusCode: http.StatusOK, resp: []byte(`{"message": "ok"}`), body: []byte(`{"key":"true","value":true}`)},
		{name: "unmarshal_error", statusCode: http.StatusInternalServerError, resp: []byte(`{"detail": "internal server error"}`), body: []byte(`{"key":not_str,"value":true}`)},
		{name: "body_error", statusCode: http.StatusInternalServerError, resp: []byte(`{"detail": "internal server error"}`), body: nil},
		{name: "empty_key", statusCode: http.StatusBadRequest, resp: []byte(`{"detail": "key is required"}`), body: []byte(`{"key":"","value":true}`)},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			reader := bytes.NewReader(tc.body)

			req, err := http.NewRequest(http.MethodPost, "/set", reader)
			assert.NoError(t, err)

			if tc.body == nil {
				mockReader := mockReader{}
				req, err = http.NewRequest(http.MethodPost, "/set", mockReader)
				assert.NoError(t, err)
			}

			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.statusCode, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			resp, err := ioutil.ReadAll(rr.Body)

			assert.NoError(t, err)
			assert.Equal(t, tc.resp, resp)
		})
	}
}

func TestNewApp(t *testing.T) {
	app := newApp()

	assert.NotNil(t, app)
	assert.NotNil(t, app.cache)
	assert.Equal(t, []byte(`{"detail": "internal server error"}`), app.InternalServerError)
	assert.Equal(t, []byte(`{"detail": "key is required"}`), app.KeyEmptyResp)
	assert.Equal(t, []byte(`{"detail": "not found"}`), app.NotFoundResp)
	assert.Equal(t, []byte(`{"message": "ok"}`), app.SetResp)
	assert.Equal(t, []byte(`{"detail": "timeout"}`), app.TimeoutResp)
}

func TestSetTimeout(t *testing.T) {
	r := newRouter()

	testcases := []struct {
		name   string
		reqUrl string
		method string
	}{
		{name: "get", reqUrl: "/get/10", method: http.MethodGet},
		{name: "set", reqUrl: "/set", method: http.MethodPost},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			ctx, cancel := context.WithCancel(context.Background())

			var req *http.Request
			if tc.method == http.MethodGet {
				req1, err := http.NewRequest(tc.method, tc.reqUrl, nil)
				assert.NoError(t, err)
				req = req1
			} else {
				req1, err := http.NewRequest(tc.method, tc.reqUrl, bytes.NewReader([]byte(`{"key":"key","value":10}`)))
				assert.NoError(t, err)
				req = req1
			}

			req = req.WithContext(ctx)
			cancel()

			r.ServeHTTP(rr, req)

			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			_, err := ioutil.ReadAll(rr.Body)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusGatewayTimeout, rr.Code)
		})
	}
}
