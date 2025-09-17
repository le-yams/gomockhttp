package mockhttp

import (
	"fmt"
	"testing"

	assertions "github.com/stretchr/testify/assert"
)

type MockT struct {
	t             *testing.T
	errorOccurred bool
	fatalOccurred bool
	errors        []error
	fatals        []error
}

func NewTestingMock(t *testing.T) *MockT {
	return &MockT{
		t:             t,
		errorOccurred: false,
		fatalOccurred: false,
		errors:        []error{},
		fatals:        []error{},
	}
}

type errorLevel int

const (
	levelError errorLevel = 0
	levelFatal errorLevel = 1
)

func (testState *MockT) register(level errorLevel, args ...any) {
	if len(args) == 0 {
		return
	}
	var errors *[]error
	switch level {
	case levelError:
		testState.errorOccurred = true
		errors = &testState.errors
	case levelFatal:
		testState.fatalOccurred = true
		errors = &testState.fatals
	default:
		panic("unknown error level")
	}

	var errs []error
	arg1 := args[0]
	if err, ok := arg1.(error); ok {
		errs = append(errs, err)
		*errors = errs
		return
	}
	
	switch t := arg1.(type) {
	case string:
		err := fmt.Errorf(t, args[1:]...)
		errs = append(errs, err)
	case error:
		errs = append(errs, t)
	default:
		_ = t
	}

	*errors = errs
}

func (testState *MockT) registerError(args ...any) {
	testState.register(levelError, args...)
}

func (testState *MockT) registerFatal(args ...any) {
	testState.register(levelFatal, args...)
}

func (testState *MockT) Error(args ...any) {
	_ = args
	testState.registerError(args...)
}

func (testState *MockT) Errorf(format string, args ...any) {
	_, _ = format, args
	testState.registerError(args...)
}

func (testState *MockT) Fatalf(format string, args ...any) {
	_, _ = format, args
	testState.registerFatal(args...)
}

func (testState *MockT) Fatal(args ...any) {
	_ = args
	testState.registerFatal(args...)
}

func formatErrors(errors []error) string {
	msg := ""
	for _, err := range errors {
		msg += "  - " + err.Error() + "\n"
	}
	return msg
}

func (testState *MockT) assertDidNotFailed() {
	if testState.errorOccurred {
		testState.t.Error("unexpected error\n" + formatErrors(testState.errors))
	}
	if testState.fatalOccurred {
		testState.t.Error("unexpected fatal\n" + formatErrors(testState.fatals))
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

func Test_api_mock_should_return_underlying_server_url(t *testing.T) {
	t.Parallel()
	// Arrange
	mockedAPI := API(NewTestingMock(t))
	defer func() { mockedAPI.Close() }()

	// Assert
	assert := assertions.New(t)
	assert.Equal(mockedAPI.testServer.URL, mockedAPI.GetURL().String())
	assert.Equal(mockedAPI.GetURL().Host, mockedAPI.GetHost())
}
