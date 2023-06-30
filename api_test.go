package mockhttp

import (
	assertions "github.com/stretchr/testify/assert"
	"testing"
)

type MockT struct {
	t             *testing.T
	errorOccurred bool
	fatalOccurred bool
}

func NewTestingMock(t *testing.T) *MockT {
	return &MockT{
		t:             t,
		errorOccurred: false,
		fatalOccurred: false,
	}
}

func (testState *MockT) Error(args ...interface{}) {
	_ = args
	testState.errorOccurred = true
}

func (testState *MockT) Errorf(format string, args ...interface{}) {
	_, _ = format, args
	testState.errorOccurred = true
}

func (testState *MockT) Fatalf(format string, args ...interface{}) {
	_, _ = format, args
	testState.fatalOccurred = true
}

func (testState *MockT) Fatal(args ...interface{}) {
	_ = args
	testState.fatalOccurred = true
}

func (testState *MockT) assertDidNotFailed() {
	if testState.errorOccurred {
		testState.t.Error("unexpected failure with error")
	}
	if testState.fatalOccurred {
		testState.t.Error("unexpected failure with fatal")
	}
}

func (testState *MockT) assertFailedWithError() {
	if !testState.errorOccurred {
		testState.t.Error("an error was expected to occur but did not")
	}
}

func (testState *MockT) assertFailedWithFatal() {
	if !testState.fatalOccurred {
		testState.t.Error("a fatal was expected to occur but did not")
	}
}

func TestApiUrl(t *testing.T) {
	// Arrange
	mockedApi := Api(NewTestingMock(t))
	defer func() { mockedApi.Close() }()

	// Assert
	assert := assertions.New(t)
	assert.Equal(mockedApi.testServer.URL, mockedApi.GetUrl().String())
	assert.Equal(mockedApi.GetUrl().Host, mockedApi.GetHost())
}
