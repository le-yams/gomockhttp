package mockhttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

type ApiMock struct {
	testServer  *httptest.Server
	calls       map[HttpCall]http.HandlerFunc
	testState   TestingT
	invocations map[HttpCall][]*Invocation
}

type HttpCall struct {
	Method string
	Path   string
}

type Invocation struct {
	request        *http.Request
	requestContent []byte
	testState      TestingT
}

func Api(testCase TestingT) *ApiMock {
	mockedApi := &ApiMock{
		calls:       map[HttpCall]http.HandlerFunc{},
		testState:   testCase,
		invocations: map[HttpCall][]*Invocation{},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, request *http.Request) {

		call := HttpCall{
			Method: strings.ToLower(request.Method),
			Path:   request.RequestURI,
		}

		bytes, err := io.ReadAll(request.Body)
		if err != nil {
			testCase.Fatal(err)
		}

		invocations := mockedApi.invocations[call]
		invocations = append(invocations, &Invocation{
			request:        request,
			requestContent: bytes,
			testState:      testCase,
		})
		mockedApi.invocations[call] = invocations

		handler := mockedApi.calls[call]
		if handler != nil {
			handler(res, request)
		} else {
			res.WriteHeader(404)
			testCase.Fatalf("unmocked invocation %s %s\n", call.Method, call.Path)
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
		mockedApi.testState.Fatal(err)
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

func (mockedApi *ApiMock) Verify(method string, path string) *CallVerifier {
	return &CallVerifier{
		api: mockedApi,
		call: &HttpCall{
			Method: strings.ToLower(method),
			Path:   path,
		},
	}
}
