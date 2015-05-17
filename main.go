package main

import (
	"blob-server/httpserver"
	"blob-server/mongostorage"
	"flag"
	"log"
)

func main() {
	httpAddr := flag.String("httpaddr", ":8887", "HTTP server address")
	mongoUrl := flag.String("mongourl", "localhost", "MongoDB URL")
	mongoPrefix := flag.String("mongoprefix", "fs", "MongoDB prefix")
	flag.Parse()

	log.Println("Starting blob-server, listening on", *httpAddr)
	log.Println("MongoDB URL", *mongoUrl, "MongoDB prefix", *mongoPrefix)
	storage, err := mongostorage.Start(*mongoUrl, *mongoPrefix)
	if err != nil {
		log.Fatal("Error starting mongo storage: ", err)
	}
	defer storage.Stop()
	err = httpserver.Serve(*httpAddr, storage)
	if err != nil {
		log.Fatal("Error starting HTTP server ", err)
	}
}
