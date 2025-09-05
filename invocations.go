package mockhttp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	assertions "github.com/stretchr/testify/assert"
)

// Invocation represents a single HTTP request made to the mock server.
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

// GetRequest returns the invocation request
func (call *Invocation) GetRequest() *http.Request {
	return call.request
}

// GetPayload returns the invocation request payload
func (call *Invocation) GetPayload() []byte {
	return call.payload
}

// WithHeader asserts that the invocation request contains the specified header
func (call *Invocation) WithHeader(name string, expectedValues ...string) *Invocation {
	values := call.request.Header.Values(name)
	assertions.Equal(call.testState, expectedValues, values)
	return call
}

// WithoutHeader asserts that the invocation request does not contain the specified header
func (call *Invocation) WithoutHeader(name string) *Invocation {
	if call.request.Header.Values(name) != nil {
		call.testState.Errorf("header '%s' found where it was expected not to")
	}
	return call
}

// WithAuthHeader asserts that the invocation request contains the specified auth header
func (call *Invocation) WithAuthHeader(scheme string, value string) {
	call.WithHeader("Authorization", scheme+" "+value)
}

// WithBasicAuthHeader asserts that the invocation request contains the specified basic auth header
func (call *Invocation) WithBasicAuthHeader(username string, password string) {
	call.WithAuthHeader("Basic", base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}

// WithBearerAuthHeader asserts that the invocation request contains the specified bearer auth header
func (call *Invocation) WithBearerAuthHeader(token string) {
	call.WithAuthHeader("Bearer", token)
}

// WithPayload asserts that the invocation request contains the specified payload
func (call *Invocation) WithPayload(expected []byte) *Invocation {
	assertions.Equal(call.testState, expected, call.GetPayload())
	return call
}

// WithStringPayload asserts that the invocation request contains the specified string payload
func (call *Invocation) WithStringPayload(expected string) *Invocation {
	assertions.Equal(call.testState, expected, string(call.GetPayload()))
	return call
}

// ReadJSONPayload reads the invocation request payload and unmarshals it into the specified object
func (call *Invocation) ReadJSONPayload(obj any) {
	err := json.Unmarshal(call.GetPayload(), obj)
	if err != nil {
		call.testState.Fatal(err)
	}
}

// WithUrlEncodedFormPayload assert that the invocation content type is application/x-www-form-urlencoded then returns
// the form to be asserted.
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

func (call *Invocation) WithQueryValue(name string, value string) *Invocation {
	query := call.request.URL.Query()
	if query.Has(name) {
		assertions.Equal(call.testState, value, query.Get(name))
	} else {
		call.testState.Errorf("query parameter '%s' not found", name)
	}
	return call
}

func (call *Invocation) WithQueryValues(values map[string]string) *Invocation {
	query := call.request.URL.Query()
	for key, value := range values {
		if query.Has(key) {
			assertions.Equal(call.testState, value, query.Get(key))
		} else {
			call.testState.Errorf("query parameter '%s' not found", key)
		}
	}
	return call
}

func (call *Invocation) WithQueryValuesExactly(values map[string]string) *Invocation {
	query := call.request.URL.Query()
	for key, value := range values {
		if query.Has(key) {
			assertions.Equal(call.testState, value, query.Get(key))
		} else {
			call.testState.Errorf("query parameter '%s' not found", key)
		}
	}

	if len(query) > len(values) {
		for key := range query {
			if _, ok := values[key]; !ok {
				call.testState.Errorf("query parameter '%s' not expected", key)
			}
		}
	}

	return call
}

// InvocationRequestForm represents a form payload of an HTTP request made to the mock server.
type InvocationRequestForm struct {
	invocation *Invocation
	formValues url.Values
}

// WithValues asserts that the form contains at least the specified values
func (form InvocationRequestForm) WithValues(expectedValues map[string]string) InvocationRequestForm {
	for key, value := range expectedValues {
		assertions.Equal(form.invocation.testState, value, form.formValues.Get(key))
	}
	return form
}

// WithValuesExactly asserts that the form contains exactly the specified values
func (form InvocationRequestForm) WithValuesExactly(expectedValues map[string]string) InvocationRequestForm {
	assertions.Equal(form.invocation.testState, len(expectedValues), len(form.formValues))
	for key, value := range expectedValues {
		assertions.Equal(form.invocation.testState, value, form.formValues.Get(key))
	}
	return form
}

func (form InvocationRequestForm) Get(s string) string {
	return form.formValues.Get(s)
}

func (form InvocationRequestForm) WithValue(key string, value string) InvocationRequestForm {
	assertions.Equal(form.invocation.testState, value, form.formValues.Get(key))
	return form
}
