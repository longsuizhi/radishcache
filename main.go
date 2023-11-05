package main

import (
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

func main() {
	radishcache.NewGroup("scores", 2<<10, radishcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:10001"
	peers := radishcache.NewHTTPPool(addr)
	log.Println("radish is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
