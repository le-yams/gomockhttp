package mockhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type StubBuilder struct {
	api   *APIMock
	call  *HTTPCall
	delay time.Duration
}

func (stub *StubBuilder) With(handler http.HandlerFunc) *APIMock {
	if stub.delay > 0 {
		stub.api.calls[*stub.call] = func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(stub.delay)
			handler(writer, request)
		}
	} else {
		stub.api.calls[*stub.call] = handler
	}
	return stub.api
}

func (stub *StubBuilder) WithDelay(delay time.Duration) *StubBuilder {
	stub.delay = delay
	return stub
}

func (stub *StubBuilder) WithStatusCode(statusCode int) *APIMock {
	return stub.With(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(statusCode)
	})
}

func (stub *StubBuilder) WithJSON(statusCode int, content interface{}) *APIMock {
	body, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}

	return stub.WithBody(statusCode, body, "application/json")
}

func (stub *StubBuilder) WithBody(statusCode int, body []byte, contentType string) *APIMock {
	return stub.With(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", contentType)
		writer.WriteHeader(statusCode)
		_, err := writer.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	})
}
