package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"radishcache"
)

// 模拟数据库存储
var db = map[string]string{
	"Tom":   "93",
	"Jack":  "89",
	"Linda": "78",
}

func createGroup() *radishcache.Group {
	return radishcache.NewGroup("scores", 2<<10, radishcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, radish *radishcache.Group) {
	peers := radishcache.NewHTTPPool(addr)
	peers.Set(addrs...)
	radish.RegisterPeers(peers)
	log.Println("radishcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, radish *radishcache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := radish.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	radish := createGroup()
	if api {
		go startAPIServer(apiAddr, radish)
	}
	startCacheServer(addrMap[port], []string(addrs), radish)
}
