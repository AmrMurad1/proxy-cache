package internal

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ProxyObject struct {
	Origin string
	Cache  map[string]*CacheObject
	Mutex  sync.RWMutex
}

func NewProxy(origin string) *ProxyObject {
	return &ProxyObject{
		Origin: origin,
		Cache:  make(map[string]*CacheObject),
		Mutex:  sync.RWMutex{},
	}
}
func (p *ProxyObject) ClearCache() {
	p.Mutex.Lock()
	p.Cache = make(map[string]*CacheObject)
	p.Mutex.Unlock()
	fmt.Println("Cache Cleared Successfully")
}

// Flow of the Request
// Create cache key -> Check if the Key is already exist
// 	Using Simple key for now for example GET method on dummyjson.com
// 	CACHE_KEY = GET:https://dummyjson.com
//
// If Yes -> Respond with the cached http.Response and Body
//    Set the HEaders
//    X-Cache ; HIT
//
// If No -> Forward the request to Origin
//    Cache the Origin Response
//    Set the Headers
//    X-Cache : MISS
//    Write the reponse

func (p *ProxyObject) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	CACHE_KEY := r.Method + ":" + r.URL.String()

	p.Mutex.RLock()
	// if cache already exixt
	if c, ok := p.Cache[CACHE_KEY]; ok {
		p.Mutex.RUnlock()
		RespondWithHeaders(w, *c.Response, c.ResponseBody, "HIT", CACHE_KEY)
		return
	}
	p.Mutex.RUnlock()

	// if cache not already exist
	fmt.Printf("Cache Not Exist for key : %s \n", CACHE_KEY)
	orginURL := p.Origin + r.URL.String()
	resp, err := http.Get(orginURL)
	if err != nil {
		http.Error(w, "Error Forwarding Request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error Forwarding Request Body", http.StatusInternalServerError)
		return
	}

	// Actually Caching the response
	p.Mutex.Lock()
	p.Cache[CACHE_KEY] = &CacheObject{
		Response:     resp,
		ResponseBody: body,
		Created:      time.Now(),
	}
	p.Mutex.Unlock()
	RespondWithHeaders(w, *resp, body, "MISS", CACHE_KEY)
}

// final response for the client
func RespondWithHeaders(w http.ResponseWriter, response http.Response, body []byte, cacheHeader string, KEY string) {
	fmt.Printf("Cache : %s %s \n", cacheHeader, KEY)

	for k, v := range response.Header {
		w.Header()[k] = v
	}
	w.Header().Set("X-Cache", cacheHeader)
	w.WriteHeader(response.StatusCode)

	w.Write(body)
}
