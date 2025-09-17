package mockhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// StubBuilder is a helper to build stubs for a specific HTTP call
type StubBuilder struct {
	api   *APIMock
	call  *HTTPCall
	delay time.Duration
}

// With creates a new stub for the HTTP call with the specified handler
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

// WithDelay adds a delay to the stub when called before it executes the corresponding handler
func (stub *StubBuilder) WithDelay(delay time.Duration) *StubBuilder {
	stub.delay = delay
	return stub
}

// WithStatusCode creates a new stub handler returning the specified status code
func (stub *StubBuilder) WithStatusCode(statusCode int) *APIMock {
	return stub.With(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(statusCode)
	})
}

// WithJSON creates a new stub handler returning the specified status code and JSON content.
// The response header "Content-Type" is set to "application/json".
func (stub *StubBuilder) WithJSON(statusCode int, content any) *APIMock {
	body, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}

	return stub.WithBody(statusCode, body, "application/json")
}

// WithBody creates a new stub handler returning the specified status code and body content.
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
