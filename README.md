# gomockhttp [![Go Report Card](https://goreportcard.com/badge/github.com/le-yams/gomockhttp)](https://goreportcard.com/report/github.com/le-yams/gomockhttp)


* Mock the external http APIs your code is calling.
* Verify/assert the calls your code performed

## Getting started

Create the API mock 
```go
func TestApiCall(t *testing.T) {
  api := httpmock.Api(t)
  //...
}
```

Stub endpoints
```go
api.
  Stub(http.MethodGet, "/endpoint").
  WithOkJson(http.StatusOK, jsonStub).

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



## Examples

```go
package main

import (
	"github.com/le-yams/gomockhttp"
	"net/http"
	"testing"
)

type Foo struct {
	Bar string `json:"foo"`
}

func TestApiCall(t *testing.T) {
	api := httpmock.Api(t)
	foo := NewFoo(api.GetUrl())

	api.Stub(http.MethodGet, "/foo").
		WithJson(http.StatusOK, &Foo{Bar: "bar"})

	bar := foo.GetBar() // the code to be tested calling the api endpoint to rerieve the value

	if bar != "bar" {
		t.Errorf("unexpected bar value: %s\n", bar)
	}
}
```

