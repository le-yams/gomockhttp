package mockhttp

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/le-yams/gotestingmock"
	assertions "github.com/stretchr/testify/assert"
)

func Test_Verify_Endpoint(t *testing.T) {
	t.Parallel()

	t.Run("HasBeenCalled()", func(t *testing.T) {
		t.Parallel()

		t.Run("passes when the endpoint was called the expected number of times", func(t *testing.T) {
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
			client := http.Client{}
			_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
			_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
			mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalled(2)

			// Assert
			testState.AssertDidNotFailed()
		})

		t.Run("fails when the endpoint was not called the expected number of times", func(t *testing.T) {
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
			client := http.Client{}
			_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
			_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
			mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalled(3)

			// Assert
			testState.AssertFailedWithFatal()
		})

		t.Run("returns the performed calls", func(t *testing.T) {
			t.Parallel()
			// Arrange
			testState := testingmock.New(t)
			mockedAPI := API(testState)
			t.Cleanup(mockedAPI.Close)

			endpoint := "/endpoint"
			endpointURL := mockedAPI.GetURL().String() + endpoint
			mockedAPI.
				Stub(http.MethodPost, endpoint).
				WithStatusCode(http.StatusOK)
			client := http.Client{}
			_, _ = client.Post(endpointURL, "application/json", bytes.NewBuffer([]byte(`{"foo": "bar"}`)))
			_, _ = client.Post(endpointURL, "text/plain", bytes.NewBuffer([]byte("Hello")))

			// Act
			calls := mockedAPI.Verify(http.MethodPost, endpoint).HasBeenCalled(2)

			// Assert
			testState.AssertDidNotFailed()
			assert := assertions.New(t)
			assert.Len(calls, 2)

			call1 := calls[0]
			assert.Equal("application/json", call1.GetRequest().Header.Get("Content-Type"))
			call1Content := struct {
				Foo string `json:"foo"`
			}{}
			call1.ReadJSONPayload(&call1Content)
			assert.Equal("bar", call1Content.Foo)

			call2 := calls[1]
			call2.
				WithHeader("Content-Type", "text/pain").
				WithStringPayload("Hello")
		})
	})

	t.Run("HasBeenCalledOnce()", func(t *testing.T) {
		t.Parallel()

		t.Run("passes when the endpoint was called exactly once", func(t *testing.T) {
			t.Parallel()
			// Arrange
			testState := testingmock.New(t)
			mockedAPI := API(testState)
			t.Cleanup(mockedAPI.Close)

			mockedAPI.
				Stub(http.MethodGet, "/endpoint").
				WithStatusCode(http.StatusOK)

			// Act
			client := http.Client{}
			_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
			mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalledOnce()

			// Assert
			testState.AssertDidNotFailed()
		})

		t.Run("fails when the endpoint was called more than once", func(t *testing.T) {
			t.Parallel()
			// Arrange
			testState := testingmock.New(t)
			mockedAPI := API(testState)
			t.Cleanup(mockedAPI.Close)

			mockedAPI.
				Stub(http.MethodGet, "/endpoint").
				WithStatusCode(http.StatusOK)

			// Act
			client := http.Client{}
			_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
			_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
			mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalledOnce()

			// Assert
			testState.AssertFailedWithFatal()
		})

		t.Run("fails when the endpoint was not called", func(t *testing.T) {
			t.Parallel()
			// Arrange
			testState := testingmock.New(t)
			mockedAPI := API(testState)
			t.Cleanup(mockedAPI.Close)

			mockedAPI.
				Stub(http.MethodGet, "/endpoint").
				WithStatusCode(http.StatusOK)

			// Act
			mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalledOnce()

			// Assert
			testState.AssertFailedWithFatal()
		})

		t.Run("returns the performed call", func(t *testing.T) {
			t.Parallel()
			// Arrange
			testState := testingmock.New(t)
			mockedAPI := API(testState)
			t.Cleanup(mockedAPI.Close)

			endpoint := "/endpoint"
			endpointURL := mockedAPI.GetURL().String() + endpoint
			mockedAPI.
				Stub(http.MethodPost, endpoint).
				WithStatusCode(http.StatusOK)
			client := http.Client{}
			_, _ = client.Post(endpointURL, "text/plain", bytes.NewBuffer([]byte("Hello")))

			// Act
			_ = mockedAPI.
				Verify(http.MethodPost, endpoint).
				HasBeenCalledOnce().
				WithStringPayload("Hello").
				WithHeader("Content-Type", "text/plain")
		})
	})

	t.Run("HasNotBeenCalled()", func(t *testing.T) {
		t.Parallel()

		t.Run("passes when the endpoint was not called", func(t *testing.T) {
			t.Parallel()
			// Arrange
			testState := testingmock.New(t)
			mockedAPI := API(testState)
			t.Cleanup(mockedAPI.Close)

			mockedAPI.
				Stub(http.MethodGet, "/endpoint").
				WithStatusCode(http.StatusOK)

			// Act
			mockedAPI.Verify(http.MethodGet, "/endpoint").HasNotBeenCalled()

			// Assert
			testState.AssertDidNotFailed()
		})
		t.Run("fails when the endpoint actually was called", func(t *testing.T) {
			t.Parallel()
			// Arrange
			testState := testingmock.New(t)
			mockedAPI := API(testState)
			t.Cleanup(mockedAPI.Close)

			mockedAPI.
				Stub(http.MethodGet, "/endpoint").
				WithStatusCode(http.StatusOK)

			// Act
			client := http.Client{}
			_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
			mockedAPI.Verify(http.MethodGet, "/endpoint").HasNotBeenCalled()

			// Assert
			testState.AssertFailedWithFatal()
		})
	})
}
