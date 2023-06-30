package mockhttp

import (
	assertions "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestVerifyingInvocationsCountPasses(t *testing.T) {
	// Arrange
	testState := NewTestingMock()
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	content := &TestDto{Value: "Hello"}
	mockedApi.
		Stub("GET", "/endpoint").
		WithJson(http.StatusOK, content)

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	mockedApi.Verify("GET", "/endpoint").HasBeenCalled(2)

	// Assert
	assert := assertions.New(t)
	assert.Equal(false, testState.DidFatalOccurred())
}

func TestVerifyingInvocationsCountFails(t *testing.T) {
	// Arrange
	testState := NewTestingMock()
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	content := &TestDto{Value: "Hello"}
	mockedApi.
		Stub("GET", "/endpoint").
		WithJson(http.StatusOK, content)

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	mockedApi.Verify("GET", "/endpoint").HasBeenCalled(3)

	// Assert
	assert := assertions.New(t)
	assert.Equal(true, testState.DidFatalOccurred())
}

func TestVerifyingSingleInvocationPasses(t *testing.T) {
	// Arrange
	testState := NewTestingMock()
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	content := &TestDto{Value: "Hello"}
	mockedApi.
		Stub("GET", "/endpoint").
		WithJson(http.StatusOK, content)

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	mockedApi.Verify("GET", "/endpoint").HasBeenCalledOnce()

	// Assert
	assert := assertions.New(t)
	assert.Equal(false, testState.DidFatalOccurred())

}

func TestVerifyingSingleInvocationFails(t *testing.T) {
	// Arrange
	testState := NewTestingMock()
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	content := &TestDto{Value: "Hello"}
	mockedApi.
		Stub("GET", "/endpoint").
		WithJson(http.StatusOK, content)

	// Act
	client := http.Client{}
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	_, _ = client.Get(mockedApi.GetUrl().String() + "/endpoint")
	mockedApi.Verify("GET", "/endpoint").HasBeenCalledOnce()

	// Assert
	assert := assertions.New(t)
	assert.Equal(true, testState.DidFatalOccurred())
}
