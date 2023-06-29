package mockhttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type ApiMock struct {
	testServer  *httptest.Server
	calls       map[HttpCall]http.HandlerFunc
	t           *testing.T
	invocations map[HttpCall]*Invocation
}

type HttpCall struct {
	Method string
	Path   string
}

type Invocation struct {
	request        *http.Request
	requestContent []byte
	t              *testing.T
}

func Api(t *testing.T) *ApiMock {
	mockedApi := &ApiMock{
		calls:       map[HttpCall]http.HandlerFunc{},
		t:           t,
		invocations: map[HttpCall]*Invocation{},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, request *http.Request) {

		call := HttpCall{
			Method: strings.ToLower(request.Method),
			Path:   request.RequestURI,
		}

		bytes, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatal(err)
		}

		invocations := mockedApi.invocations
		invocations[call] = &Invocation{
			request:        request,
			requestContent: bytes,
			t:              t,
		}
		mockedApi.invocations = invocations

		handler := mockedApi.calls[call]
		if handler != nil {
			handler(res, request)
		} else {
			res.WriteHeader(404)
			t.Fatalf("unmocked invocation %s %s\n", call.Method, call.Path)
		}
	}))
	mockedApi.testServer = testServer

	return mockedApi
}

func (mockedApi *ApiMock) Close() {
	mockedApi.testServer.Close()
}

func (mockedApi *ApiMock) GetUrl() *url.URL {
	testServerUrl, err := url.Parse(mockedApi.testServer.URL)
	if err != nil {
		mockedApi.t.Fatal(err)
	}
	return testServerUrl
}

func (mockedApi *ApiMock) GetHost() string {
	return mockedApi.GetUrl().Host
}

func (mockedApi *ApiMock) Stub(method string, path string) *StubBuilder {
	return &StubBuilder{
		api: mockedApi,
		call: &HttpCall{
			Method: strings.ToLower(method),
			Path:   path,
		},
	}
}
