package mockhttp

import (
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
)

func (mockedAPI *APIMock) testCall(method, path string, t *testing.T) *httpexpect.Response {
	e := httpexpect.Default(t, mockedAPI.GetURL().String())
	return e.Request(method, path).Expect()
}

func TestApiNotStubbedEndpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	// Act
	call := mockedAPI.testCall(http.MethodGet, "/endpoint", t)

	// Assert
	testState.assertFailedWithFatal()
	call.Status(http.StatusNotFound)
}

func TestApiStubbedEndpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	mockedAPI.
		Stub(http.MethodGet, "/endpoint").
		With(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-Type", "text/plain")
			writer.WriteHeader(http.StatusCreated)
			_, err := writer.Write([]byte("Hello"))
			if err != nil {
				t.Fatal(err)
			}
		})

	// Act
	call := mockedAPI.testCall(http.MethodGet, "/endpoint", t)

	// Assert
	call.
		Status(http.StatusCreated).
		Body().IsEqual("Hello")

	testState.assertDidNotFailed()
}

func TestApiStubbedEndpointWithJson(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	mockedAPI.
		Stub(http.MethodGet, "/endpoint").
		WithJSON(http.StatusOK, struct {
			Value string `json:"value"`
		}{Value: "Hello"})

	// Act
	call := mockedAPI.testCall(http.MethodGet, "/endpoint", t)

	// Assert
	testState.assertDidNotFailed()

	call.Header("Content-Type").IsEqual("application/json")

	responseObject := call.
		Status(http.StatusOK).
		JSON().Object()

	responseObject.Value("value").IsEqual("Hello")
}

func TestApiStubbedEndpointWithBody(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()
	body := []byte("Hello!")

	mockedAPI.
		Stub(http.MethodGet, "/endpoint").
		WithBody(http.StatusOK, body, "text/plain")

	// Act
	call := mockedAPI.testCall(http.MethodGet, "/endpoint", t)

	// Assert
	testState.assertDidNotFailed()

	call.Header("Content-Type").IsEqual("text/plain")

	call.
		Status(http.StatusOK).
		Body().IsEqual("Hello!")
}

func TestApiStubbedEndpointWithDelay(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	stubbedDelay := 500 * time.Millisecond

	mockedAPI.
		Stub(http.MethodPost, "/delayed").
		WithDelay(stubbedDelay).
		WithStatusCode(http.StatusOK)

	// Act
	call := mockedAPI.testCall(http.MethodPost, "/delayed", t)

	// Assert
	testState.assertDidNotFailed()
	call.RoundTripTime().Ge(stubbedDelay)
}
