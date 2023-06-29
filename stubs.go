package mockhttp

import (
	"encoding/json"
	"log"
	"net/http"
)

type StubBuilder struct {
	api  *ApiMock
	call *HttpCall
}

func (stub *StubBuilder) With(handler http.HandlerFunc) *ApiMock {
	stub.api.calls[*stub.call] = handler
	return stub.api
}

func (stub *StubBuilder) WithJson(statusCode int, content interface{}) *ApiMock {
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
