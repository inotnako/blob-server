package httpserver

import (
	"blob-server/storage"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

type RequestHandler struct {
	storage storage.Storage
	handlerFunc func(storage.Storage, http.ResponseWriter, *http.Request)
}

func (handler RequestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.handlerFunc(handler.storage, writer, request)
}

func FilePostHandler(storage storage.Storage, writer http.ResponseWriter, request *http.Request) {
	id, err := storage.Post(request.Body)
	if (err != nil) {
		log.Println("File post error: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.Write([]byte(id))
	}
}

func writeResponseToIdRequestError(writer http.ResponseWriter, err storage.IdRequestError) {
	if (err.NotFound()) {
		writer.WriteHeader(http.StatusNotFound)
	} else if (err.IllFormed()) {
		writer.WriteHeader(http.StatusBadRequest)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func FileGetHandler(storage storage.Storage, writer http.ResponseWriter, request *http.Request) {
	err := storage.Get(mux.Vars(request)["id"], writer)
	if (err != nil) {
		log.Println("File get error: ", err)
		writeResponseToIdRequestError(writer, err)
	}
}

func FileGetListHandler(storage storage.Storage, writer http.ResponseWriter, request *http.Request) {
	ids, err := storage.GetList()
	if (err != nil) {
		log.Println("File get list error: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		json, err := json.Marshal(ids)
		if (err != nil) {
			log.Println("JSON marshalling error: ", err)
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			writer.Header().Set("Content-Type", "application/json")
			writer.Write(json)
		}
	}
}

func FileDeleteHandler(storage storage.Storage, writer http.ResponseWriter, request *http.Request) {
	err := storage.Delete(mux.Vars(request)["id"])
	if (err != nil) {
		log.Println("File delete error: ", err)
		writeResponseToIdRequestError(writer, err)
	}
}

func Serve(addr string, storage storage.Storage) (error) {
	router := mux.NewRouter()

	router.Handle("/api/v1/file", RequestHandler{storage, FilePostHandler}).Methods("POST")
	router.Handle("/api/v1/file/:{id}", RequestHandler{storage, FileGetHandler}).Methods("GET")
	router.Handle("/api/v1/file", RequestHandler{storage, FileGetListHandler}).Methods("GET")
	router.Handle("/api/v1/file/:{id}", RequestHandler{storage, FileDeleteHandler}).Methods("DELETE")

	router.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	})

	http.Handle("/", router)

	return http.ListenAndServe(addr, nil)
}

