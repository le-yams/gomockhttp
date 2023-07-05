package mockhttp

import (
	"encoding/json"
	"log"
	"net/http"
)

type StubBuilder struct {
	api  *APIMock
	call *HTTPCall
}

func (stub *StubBuilder) With(handler http.HandlerFunc) *APIMock {
	stub.api.calls[*stub.call] = handler
	return stub.api
}

func (stub *StubBuilder) WithStatusCode(statusCode int) *APIMock {
	return stub.With(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(statusCode)
	})
}

func (stub *StubBuilder) WithJSON(statusCode int, content interface{}) *APIMock {
	return stub.With(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(statusCode)

		bytes, err := json.Marshal(content)
		if err != nil {
			log.Fatal(err)
		}
		_, err = writer.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
	})
}
