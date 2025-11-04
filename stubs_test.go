package mockhttp

import (
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/le-yams/gotestingmock"
)

func (mockedAPI *APIMock) testCall(method, path string, t *testing.T) *httpexpect.Response {
	e := httpexpect.Default(t, mockedAPI.GetURL().String())
	return e.Request(method, path).Expect()
}

func (mockedAPI *APIMock) testCallWithQuery(method, path string, queryObject any, t *testing.T) *httpexpect.Response {
	e := httpexpect.Default(t, mockedAPI.GetURL().String())
	return e.Request(method, path).WithQueryObject(queryObject).Expect()
}

func Test_StubbedEndpoint(t *testing.T) {
	t.Parallel()

	t.Run("fails when endpoint is not stubbed", func(t *testing.T) {
		t.Parallel()
		// Arrange
		testState := testingmock.New(t)
		mockedAPI := API(testState)
		t.Cleanup(mockedAPI.Close)

		// Act
		call := mockedAPI.testCall(http.MethodGet, "/endpoint", t)

		// Assert
		testState.AssertFailedWithFatal()
		call.Status(http.StatusNotFound)
	})

	t.Run("returns stubbed response", func(t *testing.T) {
		t.Parallel()
		// Arrange
		testState := testingmock.New(t)
		mockedAPI := API(testState)
		t.Cleanup(mockedAPI.Close)

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

		testState.AssertDidNotFailed()
	})

	t.Run("can access query parameters", func(t *testing.T) {
		t.Parallel()
		// Arrange
		testState := testingmock.New(t)
		mockedAPI := API(testState)
		t.Cleanup(mockedAPI.Close)

		mockedAPI.
			Stub(http.MethodGet, "/endpoint").
			With(func(writer http.ResponseWriter, request *http.Request) {
				name := request.FormValue("name")
				writer.Header().Add("Content-Type", "text/plain")
				writer.WriteHeader(http.StatusCreated)
				_, err := writer.Write([]byte("Hello " + name))
				if err != nil {
					t.Fatal(err)
				}
			})

		query := map[string]any{"name": "John"}

		// Act
		call := mockedAPI.testCallWithQuery(http.MethodGet, "/endpoint", query, t)

		// Assert
		call.
			Status(http.StatusCreated).
			Body().IsEqual("Hello John")

		testState.AssertDidNotFailed()
	})

	t.Run("can return json response", func(t *testing.T) {
		t.Parallel()
		// Arrange
		testState := testingmock.New(t)
		mockedAPI := API(testState)
		t.Cleanup(mockedAPI.Close)

		mockedAPI.
			Stub(http.MethodGet, "/endpoint").
			WithJSON(http.StatusOK, struct {
				Value string `json:"value"`
			}{Value: "Hello"})

		// Act
		call := mockedAPI.testCall(http.MethodGet, "/endpoint", t)

		// Assert
		testState.AssertDidNotFailed()

		call.Header("Content-Type").IsEqual("application/json")

		responseObject := call.
			Status(http.StatusOK).
			JSON().Object()

		responseObject.Value("value").IsEqual("Hello")
	})

	t.Run("can return specified body", func(t *testing.T) {
		t.Parallel()
		// Arrange
		testState := testingmock.New(t)
		mockedAPI := API(testState)
		t.Cleanup(mockedAPI.Close)
		body := []byte("Hello!")

		mockedAPI.
			Stub(http.MethodGet, "/endpoint").
			WithBody(http.StatusOK, body, "text/plain")

		// Act
		call := mockedAPI.testCall(http.MethodGet, "/endpoint", t)

		// Assert
		testState.AssertDidNotFailed()

		call.Header("Content-Type").IsEqual("text/plain")

		call.
			Status(http.StatusOK).
			Body().IsEqual("Hello!")
	})

	t.Run("can timeout with specified delay", func(t *testing.T) {
		t.Parallel()
		// Arrange
		testState := testingmock.New(t)
		mockedAPI := API(testState)
		t.Cleanup(mockedAPI.Close)

		stubbedDelay := 500 * time.Millisecond

		mockedAPI.
			Stub(http.MethodPost, "/delayed").
			WithDelay(stubbedDelay).
			WithStatusCode(http.StatusOK)

		// Act
		call := mockedAPI.testCall(http.MethodPost, "/delayed", t)

		// Assert
		testState.AssertDidNotFailed()
		call.RoundTripTime().Ge(stubbedDelay)
	})
}
