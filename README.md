# gomockhttp [![Go Report Card](https://goreportcard.com/badge/github.com/le-yams/gomockhttp)](https://goreportcard.com/report/github.com/le-yams/gomockhttp) [![GoDoc](https://godoc.org/github.com/le-yams/gomockhttp?status.svg)](https://godoc.org/github.com/le-yams/gomockhttp) [![Version](https://img.shields.io/github/tag/le-yams/gomockhttp.svg)](https://github.com/le-yams/gomockhttp/releases)

A fluent testing library for mocking http apis.
* Mock the external http apis your code is calling.
* Verify/assert the calls your code performed

## Getting started

Create the API mock 
```go
func TestApiCall(t *testing.T) {
  api := mockhttp.Api(t)
  defer func() { api.Close() }()
  //...
}
```

Stub endpoints
```go
api.
  Stub(http.MethodGet, "/endpoint").
  WithJson(http.StatusOK, jsonStub).

  Stub(http.MethodPost, "/endpoint").
  WithStatusCode(http.StatusCreated)
```

Call the mocked API
```go
resultBasedOnMockedResponses, err := codeCallingTheApi(api.GetUrl())
```

Verify the API invocations
```go
calls := api.
  Verify(http.MethodPost, "/endpoint").
  HasBeenCalled(2)

expectCall1 := calls[0]
expectCall1.WithPayload(expectedPayload1)

expectCall2 := calls[1]
expectCall2.WithPayload(expectedPayload2)
```



## Example

```go
package main

import (
  "github.com/le-yams/gomockhttp"
  "fmt"
  "net/http"
  "testing"
)

type FooDto struct {
  Foo string `json:"foo"`
}

func TestApiCall(t *testing.T) {
  // Arrange
  api := mockhttp.Api(t)
  defer func() { api.Close() }()

  api.
    Stub(http.MethodGet, "/foo").
    WithJson(http.StatusOK, &FooDto{Foo: "bar"})
  token := "testToken"

  //Act
  fooService := NewFooService(api.GetUrl(), token)
  foo := fooService.GetFoo() // the code actually making the http call to the api endpoint

  // Assert
  if foo != "bar" {
    t.Errorf("unexpected value: %s\n", foo)
  }

  api.
    Verify(http.MethodGet, "/foo").
    HasBeenCalledOnce().
    WithHeader("Authorization", fmt.Sprintf("Bearer %s", token))
}

```

