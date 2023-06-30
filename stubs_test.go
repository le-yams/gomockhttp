package mockhttp

import (
	"github.com/gavv/httpexpect/v2"
	assertions "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestApiNotStubbedEndpoint(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	// Act
	client := http.Client{}
	response, err := client.Get(mockedApi.GetUrl().String() + "/endpoint")

	// Assert
	assert := assertions.New(t)
	assert.NoError(err)
	testState.assertFailedWithFatal()
	assert.Equal(404, response.StatusCode)
}

func TestApiStubbedEndpoint(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedApi := Api(testState)
	defer func() { mockedApi.Close() }()

	mockedApi.
		Stub(http.MethodGet, "/endpoint").
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
	e.GET("/endpoint").
		Expect().
		Status(http.StatusCreated).
		Body().IsEqual("Hello")

	testState.assertDidNotFailed()
}

func TestApiStubbedEndpointWithJson(t *testing.T) {
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
	e := httpexpect.Default(t, mockedApi.GetUrl().String())

	// Assert
	testState.assertDidNotFailed()

	responseObject := e.GET("/endpoint").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	responseObject.Value("value").IsEqual("Hello")
}

type TestDto struct {
	Value string `json:"value"`
}
