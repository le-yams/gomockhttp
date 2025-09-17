# gomockhttp 
[![Version](https://img.shields.io/github/tag/le-yams/gomockhttp.svg)](https://github.com/le-yams/gomockhttp/releases)
[![Coverage Status](https://coveralls.io/repos/github/le-yams/gomockhttp/badge.svg?branch=master)](https://coveralls.io/github/le-yams/gomockhttp?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/le-yams/gomockhttp)](https://goreportcard.com/report/github.com/le-yams/gomockhttp)
[![GoDoc](https://godoc.org/github.com/le-yams/gomockhttp?status.svg)](https://godoc.org/github.com/le-yams/gomockhttp)
[![License](https://img.shields.io/github/license/le-yams/gomockhttp)](LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fle-yams%2Fgomockhttp.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fle-yams%2Fgomockhttp?ref=badge_shield&issueType=license)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fle-yams%2Fgomockhttp.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Fle-yams%2Fgomockhttp?ref=badge_shield&issueType=security)


A fluent testing library for mocking http apis.
* Mock the external http apis your code is calling.
* Verify/assert the calls your code performed

## Getting started

### 1.  Create the API mock 
```go
func TestApiCall(t *testing.T) {
  api := mockhttp.Api(t)
  defer func() { api.Close() }()
  //...
}
```

### 2. Stub endpoints
```go
api.
  Stub(http.MethodGet, "/endpoint").
  WithJSON(http.StatusOK, jsonStub).

  Stub(http.MethodPost, "/endpoint").
  WithStatusCode(http.StatusCreated)
```

See the [StubBuilder documentation](https://pkg.go.dev/github.com/le-yams/gomockhttp#StubBuilder) for full list of stubbing methods.

### 3. Call the mocked API
```go
resultBasedOnMockedResponses, err := codeCallingTheApi(api.GetURL())
```

### 4. Verify the API invocations
```go
calls := api.
  Verify(http.MethodPost, "/endpoint").
  HasBeenCalled(3)

expectCall1 := calls[0]
expectCall1.WithPayload(expectedPayload1)

expectCall2 := calls[1]
expectCall2.WithPayload(expectedPayload2)

expectCall2 := calls[2]
expectCall2.WithJSONPayload(map[string]any{"foo": "bar"})
```
See [CallVerifier documentation](https://pkg.go.dev/github.com/le-yams/gomockhttp#CallVerifier) for full list of verification methods.


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
    WithHeader("Content-Type", "application/json").
    WithBearerAuthHeader(token)
}

```

