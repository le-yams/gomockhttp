package mockhttp

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	assertions "github.com/stretchr/testify/assert"
)

func TestApiNotStubbedEndpoint(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	// Act
	client := http.Client{}
	response, err := client.Get(mockedAPI.GetURL().String() + "/endpoint")

	// Assert
	assert := assertions.New(t)
	assert.NoError(err)
	testState.assertFailedWithFatal()
	assert.Equal(404, response.StatusCode)
}

func TestApiStubbedEndpoint(t *testing.T) {
	// Arrange
	testState := NewTestingMock(t)
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	mockedAPI.
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
	e := httpexpect.Default(t, mockedAPI.GetURL().String())

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
	mockedAPI := API(testState)
	defer func() { mockedAPI.Close() }()

	mockedAPI.
		Stub(http.MethodGet, "/endpoint").
		WithJSON(http.StatusOK, struct {
			Value string `json:"value"`
		}{Value: "Hello"})

	// Act
	e := httpexpect.Default(t, mockedAPI.GetURL().String())

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
