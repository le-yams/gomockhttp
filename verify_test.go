package mockhttp

import (
	"bytes"
	assertions "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestVerifyingInvocationsCountPasses(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	mockedApi.
		Stub(http.MethodGet, "/endpoint").
		WithJson(http.StatusOK, struct {
			Value string `json:"value"`
		}{Value: "Hello"})

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	mockedApi.Verify(http.MethodGet, "/endpoint").HasBeenCalled(2)

	// Assert
	testState.assertDidNotFailed()
}

func TestVerifyingInvocationsCountFails(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	mockedApi.
		Stub(http.MethodGet, "/endpoint").
		WithJson(http.StatusOK, struct {
			Value string `json:"value"`
		}{Value: "Hello"})

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	mockedApi.Verify(http.MethodGet, "/endpoint").HasBeenCalled(3)

	// Assert
	testState.assertFailedWithFatal()
}

func TestVerifyingInvocationsCountReturnsThePerformedCalls(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	endpoint := "/endpoint"
	endpointUrl := mockedApi.GetUrl().String() + endpoint
	mockedApi.
		Stub(http.MethodPost, endpoint).
		WithStatusCode(http.StatusOK)
	client := http.Client{}
	_, _ = client.Post(endpointUrl, "application/json", bytes.NewBuffer([]byte(`{"foo": "bar"}`)))
	_, _ = client.Post(endpointUrl, "text/plain", bytes.NewBuffer([]byte("Hello")))

	// Act
	calls := mockedApi.Verify(http.MethodPost, endpoint).HasBeenCalled(2)

	// Assert
	testState.assertDidNotFailed()
	assert := assertions.New(t)
	assert.Len(calls, 2)

	call1 := calls[0]
	assert.Equal("application/json", call1.GetRequest().Header.Get("Content-Type"))
	call1Content := struct {
		Foo string `json:"foo"`
	}{}
	err := call1.ReadRequestContentAsJson(&call1Content)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("bar", call1Content.Foo)

	call2 := calls[1]
	assert.Equal("text/plain", call2.GetRequest().Header.Get("Content-Type"))
	assert.Equal("Hello", call2.ReadRequestContentAsString())
}

func TestVerifyingSingleInvocationPasses(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	mockedApi.
		Stub(http.MethodGet, "/endpoint").
		WithStatusCode(http.StatusOK)

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	mockedApi.Verify(http.MethodGet, "/endpoint").HasBeenCalledOnce()

	// Assert
	testState.assertDidNotFailed()

}

func TestVerifyingSingleInvocationFails(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	mockedApi.
		Stub(http.MethodGet, "/endpoint").
		WithStatusCode(http.StatusOK)

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	mockedApi.Verify(http.MethodGet, "/endpoint").HasBeenCalledOnce()

	// Assert
	testState.assertFailedWithFatal()
}

func TestVerifyingSingleInvocationReturnsThePerformedCall(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	endpoint := "/endpoint"
	endpointUrl := mockedApi.GetUrl().String() + endpoint
	mockedApi.
		Stub(http.MethodPost, endpoint).
		WithStatusCode(http.StatusOK)
	client := http.Client{}
	_, _ = client.Post(endpointUrl, "text/plain", bytes.NewBuffer([]byte("Hello")))

	// Act
	call := mockedApi.Verify(http.MethodPost, endpoint).HasBeenCalledOnce()

	// Assert
	assert := assertions.New(t)
	assert.Equal("text/plain", call.GetRequest().Header.Get("Content-Type"))
	assert.Equal("Hello", call.ReadRequestContentAsString())
}
