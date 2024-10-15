package mockhttp

import (
	"bytes"
	"net/http"
	"testing"

	assertions "github.com/stretchr/testify/assert"
)

func TestVerifyingInvocationsCountPasses(t *testing.T) {
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
	client := http.Client{}
	_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
	_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
	mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalled(2)

	// Assert
	testState.assertDidNotFailed()
}

func TestVerifyingInvocationsCountFails(t *testing.T) {
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
	client := http.Client{}
	_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
	_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
	mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalled(3)

	// Assert
	testState.assertFailedWithFatal()
}

func TestVerifyingInvocationsCountReturnsThePerformedCalls(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

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
	testState.assertDidNotFailed()
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
}

func TestVerifyingSingleInvocationPasses(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	mockedAPI.
		Stub(http.MethodGet, "/endpoint").
		WithStatusCode(http.StatusOK)

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
	mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalledOnce()

	// Assert
	testState.assertDidNotFailed()
}

func TestVerifyingSingleInvocationFails(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	mockedAPI.
		Stub(http.MethodGet, "/endpoint").
		WithStatusCode(http.StatusOK)

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
	_, _ = client.Get(mockedAPI.GetURL().String() + "/endpoint")
	mockedAPI.Verify(http.MethodGet, "/endpoint").HasBeenCalledOnce()

	// Assert
	testState.assertFailedWithFatal()
}

func TestVerifyingSingleInvocationReturnsThePerformedCall(t *testing.T) {
	t.Parallel()
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

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
}
