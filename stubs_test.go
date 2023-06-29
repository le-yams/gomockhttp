package mockhttp

import (
	"github.com/gavv/httpexpect/v2"
	assertions "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestApiNotStubbedEndpoint(t *testing.T) {
	// Arrange
	testState := &testing.T{}
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	// Act
	client := http.Client{}
	_, err := client.Get(mockedApi.GetUrl().String() + "/endpoint")

	// Assert
	assert := assertions.New(t)
	assert.Error(err)
	assert.Equal(true, testState.Failed())
}

func TestApiStubbedEndpoint(t *testing.T) {
	// Arrange
	testState := &testing.T{}
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	mockedApi.
		Stub("get", "/endpoint").
		With(func(writer http.ResponseWriter, request *http.Request) {

			writer.Header().Add("Content-Type", "text/plain")
			writer.WriteHeader(201)
			_, err := writer.Write([]byte("Hello"))
			if err != nil {
				t.Fatal(err)
			}
		})

	// Act
	e := httpexpect.Default(t, mockedApi.GetUrl().String())

	// Assert
	assert := assertions.New(t)
	e.GET("/endpoint").
		Expect().
		Status(http.StatusCreated).
		Body().IsEqual("Hello")

	assert.Equal(false, testState.Failed())
}

func TestApiStubbedEndpointWithJson(t *testing.T) {
	// Arrange
	testState := &testing.T{}
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	content := &TestDto{Value: "Hello"}
	mockedApi.
		Stub("GET", "/endpoint").
		WithJson(http.StatusOK, content)

	// Act
	e := httpexpect.Default(t, mockedApi.GetUrl().String())

	// Assert
	assert := assertions.New(t)
	assert.Equal(false, testState.Failed())

	responseObject := e.GET("/endpoint").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	responseObject.Value("value").IsEqual("Hello")
}

type TestDto struct {
	Value string `json:"value"`
}
