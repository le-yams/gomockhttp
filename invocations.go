package mockhttp

import (
	"bytes"
	"encoding/json"
	assertions "github.com/stretchr/testify/assert"
	"io"
	"net/http"
)

type Invocation struct {
	request   *http.Request
	payload   []byte
	testState TestingT
}

func newInvocation(request *http.Request, testState TestingT) *Invocation {
	var data []byte
	var err error

	if request.Body != nil {
		data, err = io.ReadAll(request.Body)
		if err != nil {
			testState.Fatal(err)
		}
		request.Body = io.NopCloser(bytes.NewReader(data))
	}

	return &Invocation{
		request:   request,
		payload:   data,
		testState: testState,
	}
}

func (call *Invocation) GetRequest() *http.Request {
	return call.request
}

func (call *Invocation) GetPayload() []byte {
	return call.payload
}

func (call *Invocation) WithHeader(name string, expectedValues ...string) *Invocation {
	values := call.request.Header.Values(name)
	assertions.Equal(call.testState, expectedValues, values)
	return call
}

func (call *Invocation) WithoutHeader(name string) *Invocation {
	if call.request.Header.Values(name) != nil {
		call.testState.Errorf("header '%s' found where it was expected not to")
	}
	return call
}

func (call *Invocation) WithPayload(expected []byte) *Invocation {
	assertions.Equal(call.testState, expected, call.GetPayload())
	return call
}

func (call *Invocation) WithStringPayload(expected string) *Invocation {
	assertions.Equal(call.testState, expected, string(call.GetPayload()))
	return call
}

func (call *Invocation) ReadJSONPayload(obj any) {
	err := json.Unmarshal(call.GetPayload(), obj)
	if err != nil {
		call.testState.Fatal(err)
	}
}
