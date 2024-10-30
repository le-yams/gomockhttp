package mockhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	assertions "github.com/stretchr/testify/assert"
)

type Invocation struct {
	request   *http.Request
	payload   []byte
	testState TestingT
}

type InvocationRequestForm struct {
	invocation *Invocation
	formValues url.Values
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

func (call *Invocation) WithUrlEncodedFormPayload() *InvocationRequestForm {
	formValues, err := url.ParseQuery(string(call.GetPayload()))
	if err != nil {
		call.testState.Fatal(err)
		return nil
	}
	return &InvocationRequestForm{
		invocation: call.WithHeader("Content-Type", "application/x-www-form-urlencoded"),
		formValues: formValues,
	}
}

func (form InvocationRequestForm) WithValues(expectedValues map[string]string) {
	for key, value := range expectedValues {
		assertions.Equal(form.invocation.testState, value, form.formValues.Get(key))
	}
}

func (form InvocationRequestForm) WithValuesExactly(expectedValues map[string]string) {
	assertions.Equal(form.invocation.testState, len(expectedValues), len(form.formValues))
	for key, value := range expectedValues {
		assertions.Equal(form.invocation.testState, value, form.formValues.Get(key))
	}
}
