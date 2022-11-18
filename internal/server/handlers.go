package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/MojtabaArezoomand/lru_cache/internal/cache"
	"github.com/gorilla/mux"
)

type (
	// App is the type for handling handlers.
	App struct {
		cache               *cache.Cache
		NotFoundResp        []byte
		TimeoutResp         []byte
		InternalServerError []byte
		KeyEmptyResp        []byte
		OKResp              []byte
	}

	// GetResponse is the response of get handler.
	GetResponse struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}

	// SetRequest is the request of set handler.
	SetRequest struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}
)

// newApp returns a new app.
func newApp() *App {
	app := App{
		cache:               cache.NewCache(),
		NotFoundResp:        []byte(`{"detail": "not found"}`),
		TimeoutResp:         []byte(`{"detail": "timeout"}`),
		InternalServerError: []byte(`{"detail": "internal server error"}`),
		KeyEmptyResp:        []byte(`{"detail": "key is required"}`),
		OKResp:              []byte(`{"message": "ok"}`),
	}
	return &app
}

// Get fetches a key from cache.
func (app *App) Get(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	w.Header().Set("Content-Type", "application/json")

	if v, err := app.cache.Get(r.Context(), key); err == cache.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write(app.NotFoundResp)
		log.Println("not found")
	} else if err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write(app.TimeoutResp)
		log.Println("error in fetching key, reason:", err)
	} else {
		resp := GetResponse{Key: key, Value: v}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(app.InternalServerError)
			log.Println("error in marshaling response, reason:", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		log.Println("GET: ok")
	}
}

// Set sets key to cache.
func (app *App) Set(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(app.InternalServerError)
		log.Println("error in reading request body, reason:", err)
		return
	}
	defer r.Body.Close()

	var req SetRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(app.InternalServerError)
		log.Println("error in unmarshaling request, reason:", err)
		return
	}

	if req.Key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(app.KeyEmptyResp)
		log.Println("key was empty string")
		return
	}

	err = app.cache.Set(r.Context(), req.Key, req.Value)
	if err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write(app.TimeoutResp)
		log.Println("error in setting key, reason:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(app.OKResp)
	log.Println("SET: ok")
}

// Flush flushes the whole cache.
func (app *App) Flush(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := app.cache.Flush(r.Context()); err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write(app.TimeoutResp)
		log.Println("error in flushing cache, reason:", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(app.OKResp)
		log.Println("FLUSH: ok")
	}
}
