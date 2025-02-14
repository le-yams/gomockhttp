package mockhttp

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	assertions "github.com/stretchr/testify/assert"
)

func TestInvocation_GetRequest(t *testing.T) {
	t.Parallel()
	expectedRequest := buildRequest(t, http.MethodGet, "/foo", nil)

	invocation := newInvocation(expectedRequest, t)

	assert := assertions.New(t)
	assert.Equal(expectedRequest, invocation.GetRequest())
}

func TestInvocation_Chaining_Verifications(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "bar")

	testState.assertDidNotFailed()
}

func TestInvocation_WithHeader_MultipleValues_Pass(t *testing.T) {
	t.Parallel()
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar1")
	request.Header.Add("foo", "bar2")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "bar1", "bar2")

	testState.assertDidNotFailed()
}

func TestInvocation_WithHeader_Fail_WhenMissingHeader(t *testing.T) {
	t.Parallel()
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "bar")

	testState.assertFailedWithError()
}

func TestInvocation_WithHeader_Fail_WhenHeaderWrongValue(t *testing.T) {
	t.Parallel()
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "notbar")

	testState.assertFailedWithError()
}

func TestInvocation_WithHeader_Fail_WhenHeaderWrongValues(t *testing.T) {
	t.Parallel()
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar1")
	request.Header.Add("foo", "bar2")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithHeader("foo", "bar")

	testState.assertFailedWithError()
}

func TestInvocation_WithoutHeader_Pass(t *testing.T) {
	t.Parallel()
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithoutHeader("foo")

	testState.assertDidNotFailed()
}

func TestInvocation_WithoutHeader_Fail(t *testing.T) {
	t.Parallel()
	request := buildRequest(t, http.MethodGet, "/endpoint", nil)
	request.Header.Add("foo", "bar")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithoutHeader("foo")

	testState.assertFailedWithError()
}

func TestInvocation_WithPayload_Pass(t *testing.T) {
	t.Parallel()
	testState := NewTestingMock(t)

	request := buildRequestWithBody(t, []byte("foo"))
	invocation := newInvocation(request, testState)

	invocation.WithPayload([]byte("foo"))

	testState.assertDidNotFailed()
}

func TestInvocation_WithPayload_Fail(t *testing.T) {
	t.Parallel()
	testState := NewTestingMock(t)
	request := buildRequestWithBody(t, []byte{42})
	invocation := newInvocation(request, testState)

	invocation.WithPayload([]byte{43})

	testState.assertFailedWithError()
}

func TestInvocation_GetRequestContent(t *testing.T) {
	t.Parallel()
	expectedRequestContent := []byte{42}
	request := buildRequestWithBody(t, expectedRequestContent)
	invocation := newInvocation(request, t)

	assert := assertions.New(t)
	assert.Equal(expectedRequestContent, invocation.GetPayload())
}

func TestInvocation_WithStringPayload_Pass(t *testing.T) {
	t.Parallel()
	request := buildRequestWithBody(t, []byte("foo"))
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithStringPayload("foo")

	testState.assertDidNotFailed()
}

func TestInvocation_WithStringPayload_Fail(t *testing.T) {
	t.Parallel()
	request := buildRequestWithBody(t, []byte("foo"))
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithStringPayload("notfoo")

	testState.assertFailedWithError()
}

func TestInvocation_ReadJsonPayload(t *testing.T) {
	t.Parallel()
	request := buildRequestWithBody(t, []byte(`{"foo":"bar"}`))
	invocation := newInvocation(request, t)

	json := struct {
		Foo string `json:"foo"`
	}{}

	invocation.ReadJSONPayload(&json)

	assert := assertions.New(t)
	assert.Equal("bar", json.Foo)
}

func TestInvocation_ReadJsonPayload_ErrorHandling(t *testing.T) {
	t.Parallel()
	request := buildRequestWithBody(t, []byte(`{"invalid json"}`))
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	json := struct {
		Foo string `json:"foo"`
	}{}
	invocation.ReadJSONPayload(&json)

	testState.assertFailedWithFatal()
}

func TestInvocation_WithUrlEncodedForm_Fail(t *testing.T) {
	t.Parallel()
	request := buildRequestWithBody(t, []byte("key1=value1&key2=value+2%21"))

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithUrlEncodedFormPayload()
	testState.assertFailedWithError()
}

func TestInvocation_WithUrlEncodedForm_Get(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedValue2 := "value 2!"
	request := buildRequestWithBody(t, []byte("key1=value1&key2="+url.QueryEscape(expectedValue2)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	// Act
	requestFormPayload := invocation.WithUrlEncodedFormPayload()

	// Assert
	assert := assertions.New(t)
	assert.Equal("value1", requestFormPayload.Get("key1"))
	assert.Equal(expectedValue2, requestFormPayload.Get("key2"))
}

func TestInvocation_WithUrlEncodedForm_Values_Pass(t *testing.T) {
	t.Parallel()

	testCases := []map[string]string{
		{"key1": "value1"},
		{"key1": "value1", "key2": "value 2!"},
	}

	for i := range testCases {
		values := testCases[i]
		t.Run("", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte("key1=value1&key2=value+2%21"))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithUrlEncodedFormPayload().WithValues(values)

			testState.assertDidNotFailed()
		})
	}
}

func TestInvocation_WithUrlEncodedForm_Values_Fail(t *testing.T) {
	t.Parallel()

	testCases := []map[string]string{
		{"key1": "not value1"},
		{"not_key1": "value1"},
	}

	for i := range testCases {
		values := testCases[i]
		t.Run("", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte("key1=value1&key2=value+2%21"))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithUrlEncodedFormPayload().WithValues(values)

			testState.assertFailedWithError()
		})
	}
}

func TestInvocation_WithUrlEncodedForm_ValuesExactly_Pass(t *testing.T) {
	t.Parallel()
	request := buildRequestWithBody(t, []byte("key1=value1&key2=value+2%21"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)
	invocation.
		WithUrlEncodedFormPayload().
		WithValuesExactly(map[string]string{"key1": "value1", "key2": "value 2!"})

	testState.assertDidNotFailed()
}

func TestInvocation_WithUrlEncodedForm_ValuesExactly_Fail(t *testing.T) {
	t.Parallel()

	testCases := []map[string]string{
		{"key1": "value1"},
		{"key1": "not value1"},
		{"not_key1": "value1"},
	}

	for i := range testCases {
		values := testCases[i]
		t.Run("", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte("key1=value1&key2=value+2%21"))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithUrlEncodedFormPayload().WithValuesExactly(values)

			testState.assertFailedWithError()
		})
	}
}

func TestInvocation_WithUrlEncodedForm_Value_Pass(t *testing.T) {
	t.Parallel()
	expectedValue2 := "value 2!"
	request := buildRequestWithBody(t, []byte("key1=value1&key2="+url.QueryEscape(expectedValue2)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	invocation.WithUrlEncodedFormPayload().
		WithValue("key1", "value1").
		WithValue("key2", expectedValue2)

	testState.assertDidNotFailed()
}

func TestInvocation_WithUrlEncodedForm_Value_Fail(t *testing.T) {
	t.Parallel()

	key := "k"
	value := "v"
	testCases := []map[string]string{
		{key: "not" + value},
		{"not" + key: value},
	}

	for i := range testCases {
		values := testCases[i]
		body := ""
		for k, v := range values {
			if body != "" {
				body += "&"
			}
			body += k + "=" + v
		}

		t.Run("", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte(body))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithUrlEncodedFormPayload().WithValue(key, value)

			testState.assertFailedWithError()
		})
	}
}

func TestInvocation_WithUrlEncodedForm_Chainable_Methods(t *testing.T) {
	t.Parallel()
	expectedValue2 := "value 2!"
	request := buildRequestWithBody(t, []byte("key1=value1&key2="+url.QueryEscape(expectedValue2)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	testState := NewTestingMock(t)
	invocation := newInvocation(request, testState)

	_ = invocation.WithUrlEncodedFormPayload().
		WithValues(map[string]string{"key1": "value1"}).
		WithValue("key2", expectedValue2).
		Get("key1")

	testState.assertDidNotFailed()
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
