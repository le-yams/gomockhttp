package mockhttp

import (
	assertions "github.com/stretchr/testify/assert"
	"testing"
)

type MockT struct {
	errorOccurred bool
	fatalOccurred bool
}

func NewTestingMock() *MockT {
	return &MockT{
		errorOccurred: false,
		fatalOccurred: false,
	}
}

func (t *MockT) Fatalf(format string, args ...interface{}) {
	_, _ = format, args
	t.fatalOccurred = true
}

func (t *MockT) Fatal(args ...interface{}) {
	_ = args
	t.fatalOccurred = true
}

func (t *MockT) DidErrorOccurred() bool {
	return t.errorOccurred
}

func (t *MockT) DidFatalOccurred() bool {
	return t.fatalOccurred
}

func TestApiUrl(t *testing.T) {
	// Arrange
	mockedApi := Api(NewTestingMock())
	defer func() { mockedApi.Close() }()

	// Assert
	assert := assertions.New(t)
	assert.Equal(mockedApi.testServer.URL, mockedApi.GetUrl().String())
	assert.Equal(mockedApi.GetUrl().Host, mockedApi.GetHost())
}
