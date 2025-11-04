package mockhttp

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"

	"github.com/le-yams/gotestingmock"
)

type T testingmock.IT

// APIMock is a representation of a mocked API. It allows to stub HTTP calls and verify invocations.
type APIMock struct {
	testServer  *httptest.Server
	calls       map[HTTPCall]http.HandlerFunc
	testState   T
	invocations map[HTTPCall][]*Invocation
	mu          sync.Mutex
}

// HTTPCall is a simple representation of an endpoint call.
type HTTPCall struct {
	Method string
	Path   string
}

// API creates a new APIMock instance and starts a server exposing it. The server is automatically stopped during test cleanup.
func API(testState T) *APIMock {
	mockedAPI := &APIMock{
		calls:       map[HTTPCall]http.HandlerFunc{},
		testState:   testState,
		invocations: map[HTTPCall][]*Invocation{},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, request *http.Request) {
		call := HTTPCall{
			Method: strings.ToLower(request.Method),
			Path:   request.URL.Path,
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
			res.WriteHeader(http.StatusNotFound)
			testState.Fatalf("unmocked invocation %s %s\n", call.Method, call.Path)
		}
	}))
	mockedAPI.testServer = testServer
	testState.Cleanup(mockedAPI.Close)

	return mockedAPI
}

// Close stops the underlying server. This method is automatically called during test cleanup.
func (mockedAPI *APIMock) Close() {
	mockedAPI.testServer.Close()
}

// GetURL returns the URL of the API underlying server.
func (mockedAPI *APIMock) GetURL() *url.URL {
	testServerURL, err := url.Parse(mockedAPI.testServer.URL)
	if err != nil {
		mockedAPI.testState.Fatal(err)
	}
	return testServerURL
}

// GetHost returns the host of the API underlying server.
func (mockedAPI *APIMock) GetHost() string {
	return mockedAPI.GetURL().Host
}

// Stub creates a new StubBuilder instance for the given method and path.
func (mockedAPI *APIMock) Stub(method string, path string) *StubBuilder {
	return &StubBuilder{
		api: mockedAPI,
		call: &HTTPCall{
			Method: strings.ToLower(method),
			Path:   path,
		},
	}
}

// Verify creates a new CallVerifier instance for the given method and path.
func (mockedAPI *APIMock) Verify(method string, path string) *CallVerifier {
	return &CallVerifier{
		api: mockedAPI,
		call: &HTTPCall{
			Method: strings.ToLower(method),
			Path:   path,
		},
	}
}
