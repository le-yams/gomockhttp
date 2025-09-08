package mockhttp

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	assertions "github.com/stretchr/testify/assert"
)

func Test_invocations(t *testing.T) {
	t.Parallel()

	t.Run("should returns call request", func(t *testing.T) {
		t.Parallel()
		expectedRequest := buildRequest(t, http.MethodPost, "/foo")

		invocation := newInvocation(expectedRequest, t)

		assert := assertions.New(t)
		assert.Equal(expectedRequest, invocation.GetRequest())
	})

	t.Run("WithHeader() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass with expected header", func(t *testing.T) {
			t.Parallel()
			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("foo", "bar")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithHeader("foo", "bar")

			testState.assertDidNotFailed()
		})

		t.Run("pass with expected header having multiple values", func(t *testing.T) {
			t.Parallel()
			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("foo", "bar1")
			request.Header.Add("foo", "bar2")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithHeader("foo", "bar1", "bar2")

			testState.assertDidNotFailed()
		})

		t.Run("fails when header is missing", func(t *testing.T) {
			t.Parallel()
			request := buildRequest(t, http.MethodGet, "/endpoint")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithHeader("foo", "bar")

			testState.assertFailedWithError()
		})

		t.Run("fails when header has wrong value", func(t *testing.T) {
			t.Parallel()
			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("foo", "bar")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithHeader("foo", "notbar")

			testState.assertFailedWithError()
		})

		t.Run("fails when header has wrong values", func(t *testing.T) {
			t.Parallel()
			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("foo", "bar1")
			request.Header.Add("foo", "bar2")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithHeader("foo", "bar")

			testState.assertFailedWithError()
		})
	})

	t.Run("WithoutHeader() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass when header is missing", func(t *testing.T) {
			t.Parallel()

			request := buildRequest(t, http.MethodGet, "/endpoint")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithoutHeader("foo")

			testState.assertDidNotFailed()
		})

		t.Run("fail when header is present", func(t *testing.T) {
			t.Parallel()
			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("foo", "bar")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithoutHeader("foo")

			testState.assertFailedWithError()
		})
	})

	t.Run("WithAuthHeader() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass with expected auth header", func(t *testing.T) {
			t.Parallel()
			scheme := "foo"
			value := "bar"

			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("Authorization", scheme+" "+value)

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithAuthHeader(scheme, value)

			testState.assertDidNotFailed()
		})

		t.Run("fail when header is missing", func(t *testing.T) {
			t.Parallel()
			scheme := "foo"
			value := "bar"

			request := buildRequest(t, http.MethodGet, "/endpoint")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithAuthHeader(scheme, value)

			testState.assertFailedWithError()
		})

		t.Run("fail with wrong scheme", func(t *testing.T) {
			t.Parallel()
			scheme := "foo"
			value := "bar"

			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("Authorization", "not"+scheme+" "+value)

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithAuthHeader(scheme, value)

			testState.assertFailedWithError()
		})

		t.Run("fail with wrong value", func(t *testing.T) {
			t.Parallel()
			scheme := "foo"
			value := "bar"

			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("Authorization", scheme+" not"+value)

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithAuthHeader(scheme, value)

			testState.assertFailedWithError()
		})
	})

	t.Run("WithBasicAuthHeader() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass with expected basic auth header", func(t *testing.T) {
			t.Parallel()
			username := "foo"
			password := "bar"

			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.SetBasicAuth(username, password)

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithBasicAuthHeader(username, password)

			testState.assertDidNotFailed()
		})

		t.Run("fail when header is missing", func(t *testing.T) {
			t.Parallel()
			username := "foo"
			password := "bar"

			request := buildRequest(t, http.MethodGet, "/endpoint")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithBasicAuthHeader(username, password)

			testState.assertFailedWithError()
		})

		t.Run("fail with wrong username", func(t *testing.T) {
			t.Parallel()
			username := "foo"
			password := "bar"

			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.SetBasicAuth("not"+username, password)

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithBasicAuthHeader(username, password)

			testState.assertFailedWithError()
		})

		t.Run("fail with wrong password", func(t *testing.T) {
			t.Parallel()
			username := "foo"
			password := "bar"

			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.SetBasicAuth(username, "not"+password)

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithBasicAuthHeader(username, password)

			testState.assertFailedWithError()
		})
	})

	t.Run("WithBearerAuthHeader() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass with expected bearer auth header", func(t *testing.T) {
			t.Parallel()

			token := "foo"

			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("Authorization", "Bearer "+token)

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithBearerAuthHeader(token)

			testState.assertDidNotFailed()
		})

		t.Run("fail when header is missing", func(t *testing.T) {
			t.Parallel()
			token := "foo"

			request := buildRequest(t, http.MethodGet, "/endpoint")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithBearerAuthHeader(token)

			testState.assertFailedWithError()
		})

		t.Run("fail with wrong token", func(t *testing.T) {
			t.Parallel()
			token := "foo"

			request := buildRequest(t, http.MethodGet, "/endpoint")
			request.Header.Add("Authorization", "Bearer wrong"+token)

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithBearerAuthHeader(token)

			testState.assertFailedWithError()
		})
	})

	t.Run("WithPayload() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass with expected payload", func(t *testing.T) {
			t.Parallel()
			testState := NewTestingMock(t)

			request := buildRequestWithBody(t, []byte("foo"))
			invocation := newInvocation(request, testState)

			invocation.WithPayload([]byte("foo"))

			testState.assertDidNotFailed()
		})

		t.Run("fail with wrong payload", func(t *testing.T) {
			t.Parallel()
			testState := NewTestingMock(t)
			request := buildRequestWithBody(t, []byte{42})
			invocation := newInvocation(request, testState)

			invocation.WithPayload([]byte{43})

			testState.assertFailedWithError()
		})

		t.Run("fail without payload", func(t *testing.T) {
			t.Parallel()
			testState := NewTestingMock(t)
			request := buildRequestWithoutBody(t)
			invocation := newInvocation(request, testState)

			invocation.WithPayload([]byte{43})

			testState.assertFailedWithError()
		})
	})

	t.Run("GetRequestContent() should return the request content", func(t *testing.T) {
		t.Parallel()
		expectedRequestContent := []byte{42}
		request := buildRequestWithBody(t, expectedRequestContent)
		invocation := newInvocation(request, t)

		assert := assertions.New(t)
		assert.Equal(expectedRequestContent, invocation.GetPayload())
	})

	t.Run("WithStringPayload() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass with expected string payload", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte("foo"))
			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithStringPayload("foo")

			testState.assertDidNotFailed()
		})

		t.Run("fail with wrong string payload", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte("foo"))
			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithStringPayload("notfoo")

			testState.assertFailedWithError()
		})

		t.Run("fail without payload", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithoutBody(t)
			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithStringPayload("notfoo")

			testState.assertFailedWithError()
		})
	})

	t.Run("ReadJSONPayload() should", func(t *testing.T) {
		t.Run("unmarshal json payload", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte(`{"foo":"bar"}`))
			invocation := newInvocation(request, t)

			json := struct {
				Foo string `json:"foo"`
			}{}

			invocation.ReadJSONPayload(&json)

			assert := assertions.New(t)
			assert.Equal("bar", json.Foo)
		})

		t.Run("fail the test on unmarshal error", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte(`{"invalid json"}`))
			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			json := struct {
				Foo string `json:"foo"`
			}{}
			invocation.ReadJSONPayload(&json)

			testState.assertFailedWithFatal()
		})
	})

	t.Run("WithUrlEncodedForm()", func(t *testing.T) {
		t.Parallel()

		t.Run("should fail when no Content-Type header", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte("key1=value1&key2=value+2%21"))

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithUrlEncodedFormPayload()
			testState.assertFailedWithError()
		})

		t.Run("should fail when Content-Type header is not  application/x-www-form-urlencoded", func(t *testing.T) {
			t.Parallel()
			request := buildRequestWithBody(t, []byte("key1=value1&key2=value+2%21"))
			request.Header.Add("Content-Type", "not a form")

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.WithUrlEncodedFormPayload()
			testState.assertFailedWithError()
		})

		t.Run("Get() should return parameter value", func(t *testing.T) {
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
		})

		t.Run("WithValues() should", func(t *testing.T) {
			t.Parallel()

			t.Run("pass when form contains the expected values", func(t *testing.T) {
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
			})

			t.Run("fail when form does not contains the expected values", func(t *testing.T) {
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
			})
		})

		t.Run("WithValuesExactly() should", func(t *testing.T) {
			t.Parallel()

			t.Run("pass when form contains exactly the expected values", func(t *testing.T) {
				t.Parallel()
				request := buildRequestWithBody(t, []byte("key1=value1&key2=value+2%21"))
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				testState := NewTestingMock(t)
				invocation := newInvocation(request, testState)
				invocation.
					WithUrlEncodedFormPayload().
					WithValuesExactly(map[string]string{"key1": "value1", "key2": "value 2!"})

				testState.assertDidNotFailed()
			})

			t.Run("fail when form does not contains exactly the expected values", func(t *testing.T) {
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
			})
		})

		t.Run("WithValue() should", func(t *testing.T) {
			t.Parallel()

			t.Run("pass when form contains the expected value", func(t *testing.T) {
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
			})

			t.Run("fail when form does not contains the expected value", func(t *testing.T) {
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
			})
		})
	})

	t.Run("WithQueryValue() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass when the URL query contains the expected parameter", func(t *testing.T) {
			t.Parallel()

			value := "value !"
			request := buildRequestWithURL(t, "/endpoint?foo=bar&bar=foo&value="+url.QueryEscape(value))

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.
				WithQueryValue("foo", "bar").
				WithQueryValue("bar", "foo").
				WithQueryValue("value", value)

			testState.assertDidNotFailed()
		})

		t.Run("fail when the URL query does not contains the expected parameter", func(t *testing.T) {
			t.Parallel()

			key := "foo"
			value := "bar"
			useCases := []struct {
				name  string
				key   string
				value string
			}{
				{"not the expected value", key, "not" + value},
				{"key not present", "not" + key, value},
			}
			request := buildRequestWithURL(t, "/endpoint?foo=bar")

			for i := range useCases {
				useCase := useCases[i]
				t.Run(useCase.name, func(t *testing.T) {
					t.Parallel()
					testState := NewTestingMock(t)
					invocation := newInvocation(request, testState)

					invocation.WithQueryValue(useCase.key, useCase.value)

					testState.assertFailedWithError()
				})
			}
		})
	})

	t.Run("WithQueryValues() should", func(t *testing.T) {
		t.Parallel()

		t.Run("pass when the URL query contains the expected parameters", func(t *testing.T) {
			t.Parallel()
			value := "value !"
			request := buildRequestWithURL(t, "/endpoint?foo=bar&bar=foo&value="+url.QueryEscape(value))

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.
				WithQueryValues(map[string]string{
					"foo":   "bar",
					"value": value,
				})

			testState.assertDidNotFailed()
		})

		t.Run("fail when the URL query does not contains the expected parameters", func(t *testing.T) {
			t.Parallel()

			values := map[string]string{
				"foo": "bar",
			}
			useCases := []struct {
				name string
				url  string
			}{
				{"not the expected value", "/endpoint?foo=notbar"},
				{"key not present", "/endpoint?notfoo=bar"},
			}

			for i := range useCases {
				useCase := useCases[i]
				t.Run(useCase.name, func(t *testing.T) {
					t.Parallel()
					request := buildRequestWithURL(t, useCase.url)
					testState := NewTestingMock(t)
					invocation := newInvocation(request, testState)

					invocation.WithQueryValues(values)

					testState.assertFailedWithError()
				})
			}
		})
	})

	t.Run("WithQueryValuesExactly() should", func(t *testing.T) {
		t.Parallel()
		t.Run("pass when the URL query contains exactly the expected parameters", func(t *testing.T) {
			t.Parallel()
			value := "value !"
			request := buildRequestWithURL(t, "/endpoint?foo=bar&value="+url.QueryEscape(value))

			testState := NewTestingMock(t)
			invocation := newInvocation(request, testState)

			invocation.
				WithQueryValuesExactly(map[string]string{
					"foo":   "bar",
					"value": value,
				})

			testState.assertDidNotFailed()
		})
		t.Run("fail when the URL query does not contains exactly the expected parameters", func(t *testing.T) {
			t.Parallel()

			values := map[string]string{
				"foo": "bar",
			}
			useCases := []struct {
				name string
				url  string
			}{
				{"not the expected value", "/endpoint?foo=notbar"},
				{"key not present", "/endpoint?notfoo=bar"},
				{"key present but not expected", "/endpoint?foo=bar&bar=foo"},
			}

			for i := range useCases {
				useCase := useCases[i]
				t.Run(useCase.name, func(t *testing.T) {
					t.Parallel()
					request := buildRequestWithURL(t, useCase.url)
					testState := NewTestingMock(t)
					invocation := newInvocation(request, testState)

					invocation.WithQueryValuesExactly(values)

					testState.assertFailedWithError()
				})
			}
		})
	})
}

func buildRequest(t *testing.T, method string, url string) *http.Request {
	request, err := http.NewRequest(method, url, nil)
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

func buildRequestWithoutBody(t *testing.T) *http.Request {
	request, err := http.NewRequest("dummy", "dummy", nil)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func buildRequestWithURL(t *testing.T, url string) *http.Request {
	request, err := http.NewRequest("dummy", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return request
}
