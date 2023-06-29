package mockhttp

import (
	assertions "github.com/stretchr/testify/assert"
	"testing"
)

func TestApiUrl(t *testing.T) {
	// Arrange
	mockedApi := Api(t)
	defer func() { mockedApi.Close() }()

	// Assert
	assert := assertions.New(t)

	assert.Equal(mockedApi.testServer.URL, mockedApi.GetUrl().String())
	assert.Equal(mockedApi.GetUrl().Host, mockedApi.GetHost())
}
