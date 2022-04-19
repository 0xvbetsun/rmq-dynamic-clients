// Entry point for serving AsyncApi documentation https://www.asyncapi.com/
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vbetsun/rmq-dynamic-clients/configs"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Printf("start listening on port :%d", cfg.Docs.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Docs.Port), nil))
}
