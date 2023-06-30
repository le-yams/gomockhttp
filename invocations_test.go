package mockhttp

import (
	"bytes"
	assertions "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestInvocation_GetRequest(t *testing.T) {
	expectedRequest := buildRequest(t, http.MethodGet, "/endpoint", nil)

	invocation := newInvocation(expectedRequest, t)

	assert := assertions.New(t)
	assert.Equal(expectedRequest, invocation.GetRequest())
}

func TestInvocation_Chaining_Verifications(t *testing.T) {
	payload := []byte{42}
	request := buildRequest(t, http.MethodPost, "/endpoint", payload)
	request.Header.Add("foo", "bar")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.
		WithHeader("foo", "bar").
		WithPayload(payload).
		WithoutHeader("dummy").
		WithHeader("foo", "bar")

	testState.assertDidNotFailed()
}

func TestInvocation_WithHeader_Pass(t *testing.T) {
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "bar")

	testState.assertDidNotFailed()
}

func TestInvocation_WithHeader_MultipleValues_Pass(t *testing.T) {
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar1")
	request.Header.Add("foo", "bar2")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "bar1", "bar2")

	testState.assertDidNotFailed()
}

func TestInvocation_WithHeader_Fail_WhenMissingHeader(t *testing.T) {
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "bar")

	testState.assertFailedWithError()
}

func TestInvocation_WithHeader_Fail_WhenHeaderWrongValue(t *testing.T) {
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "notbar")

	testState.assertFailedWithError()
}

func TestInvocation_WithHeader_Fail_WhenHeaderWrongValues(t *testing.T) {
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar1")
	request.Header.Add("foo", "bar2")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "bar")

	testState.assertFailedWithError()
}

func TestInvocation_WithoutHeader_Pass(t *testing.T) {
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithoutHeader("foo")

	testState.assertDidNotFailed()
}

func TestInvocation_WithoutHeader_Fail(t *testing.T) {
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithoutHeader("foo")

	testState.assertFailedWithError()
}

func TestInvocation_WithPayload_Pass(t *testing.T) {
	testState := NewTestingMock(t)

	request := buildRequestWithBody(t, []byte("foo"))
	invocation := newInvocation(request, testState)

	invocation.WithPayload([]byte("foo"))

	testState.assertDidNotFailed()

}

func TestInvocation_WithPayload_Fail(t *testing.T) {
	testState := NewTestingMock(t)
	request := buildRequestWithBody(t, []byte{42})
	invocation := newInvocation(request, testState)

	invocation.WithPayload([]byte{43})

	testState.assertFailedWithError()
}

func TestInvocation_GetRequestContent(t *testing.T) {
	expectedRequestContent := []byte{42}
	request := buildRequestWithBody(t, expectedRequestContent)
	invocation := newInvocation(request, t)

	assert := assertions.New(t)
	assert.Equal(expectedRequestContent, invocation.GetPayload())
}

func TestInvocation_WithStringPayload_Pass(t *testing.T) {
	request := buildRequestWithBody(t, []byte("foo"))
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithStringPayload("foo")

	testState.assertDidNotFailed()
}

func TestInvocation_WithStringPayload_Fail(t *testing.T) {
	request := buildRequestWithBody(t, []byte("foo"))
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithStringPayload("notfoo")

	testState.assertFailedWithError()
}

func TestInvocation_ReadJsonPayload(t *testing.T) {
	request := buildRequestWithBody(t, []byte(`{"foo":"bar"}`))
	invocation := newInvocation(request, t)

	json := struct {
		Foo string `json:"foo"`
	}{}

	invocation.ReadJsonPayload(&json)

	assert := assertions.New(t)
	assert.Equal("bar", json.Foo)
}

func TestInvocation_ReadJsonPayload_ErrorHandling(t *testing.T) {
	request := buildRequestWithBody(t, []byte(`{"invalid json"}`))
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	json := struct {
		Foo string `json:"foo"`
	}{}
	invocation.ReadJsonPayload(&json)

	testState.assertFailedWithFatal()
}

func buildRequest(t *testing.T, method string, url string, data []byte) *http.Request {
	request, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	return request

}

func buildRequestWithBody(t *testing.T, data []byte) *http.Request {
	request, err := http.NewRequest("dummy", "dummy", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	return request

}
