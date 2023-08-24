package mockhttp

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
)

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

type APIMock struct {
	testServer  *httptest.Server
	calls       map[HTTPCall]http.HandlerFunc
	testState   TestingT
	invocations map[HTTPCall][]*Invocation
	mu          sync.Mutex
}

type HTTPCall struct {
	Method string
	Path   string
}

func API(testState TestingT) *APIMock {
	mockedAPI := &APIMock{
		calls:       map[HTTPCall]http.HandlerFunc{},
		testState:   testState,
		invocations: map[HTTPCall][]*Invocation{},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, request *http.Request) {
		call := HTTPCall{
			Method: strings.ToLower(request.Method),
			Path:   request.RequestURI,
		}

		mockedAPI.mu.Lock()
		invocations := mockedAPI.invocations[call]
		invocations = append(invocations, newInvocation(request, testState))
		mockedAPI.invocations[call] = invocations
		mockedAPI.mu.Unlock()

		handler := mockedAPI.calls[call]
		if handler != nil {
			handler(res, request)
		} else {
			res.WriteHeader(404)
			testState.Fatalf("unmocked invocation %s %s\n", call.Method, call.Path)
		}
	}))
	mockedAPI.testServer = testServer

	return mockedAPI
}

func (mockedAPI *APIMock) Close() {
	mockedAPI.testServer.Close()
}

func (mockedAPI *APIMock) GetURL() *url.URL {
	testServerURL, err := url.Parse(mockedAPI.testServer.URL)
	if err != nil {
		mockedAPI.testState.Fatal(err)
	}
	return testServerURL
}

func (mockedAPI *APIMock) GetHost() string {
	return mockedAPI.GetURL().Host
}

func (mockedAPI *APIMock) Stub(method string, path string) *StubBuilder {
	return &StubBuilder{
		api: mockedAPI,
		call: &HTTPCall{
			Method: strings.ToLower(method),
			Path:   path,
		},
	}
}

func (mockedAPI *APIMock) Verify(method string, path string) *CallVerifier {
	return &CallVerifier{
		api: mockedAPI,
		call: &HTTPCall{
			Method: strings.ToLower(method),
			Path:   path,
		},
	}
}
