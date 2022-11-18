### LRU Cache implementation in Golang (Least Recently Used Items)
___

This is a `LRU Cache` HTTP server which has two endpoints to get and set stuff to cache.  
All of the `SET`, `GET` and `FLUSH` operations have **O(1)** time complexity.

#### Setup

 Just run server using this command `go run go run cmd/lrucache/main.go`.  
Now you can store key-value pairs in the cache like this:  
```
curl --request POST --data '{"key":"first_key","value":[1, "val"]}' http://127.0.0.1:2376/set

// Response
{"message": "ok"}
```
to get key from cache:
```
curl http://127.0.0.1:2376/get/first_key

// Response
{"key":"first_key","value":[1,"val"]}
```

to flush the whole cache:
```
curl http://127.0.0.1:2376/flush

// Response 
{"message": "ok"}
```

#### Endpoints

 1. GET `/get/{key}`
 2. POST `/set`
    - request body `{"key": "string", "value": any}`
 3. GET `/flush`
 
#### Config Environment Variables
 1. **CACHE_CAPACITY:** maximum stored key-value pairs. defaults to `2048`.
 2. **SERVER_ADDRESS:** address which server will be served on, defaults to `127.0.0.1:2376`
 3. **SERVER_WRITE_TIMEOUT:** The maximum duration before timing out writes of the response.  defaults to `1s`.
 4. **SERVER_READ_TIMEOUT:** the maximum duration for reading the entire request, including the body. A zero or negative value means there will be no timeout. defaults to `1s`.
