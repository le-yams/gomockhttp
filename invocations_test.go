package mockhttp

import (
	assertions "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestInvocation_GetRequest(t *testing.T) {
	expectedRequest, _ := http.NewRequest(http.MethodGet, "/endpoint", nil)

	invocation := &Invocation{
		request: expectedRequest,
	}

	assert := assertions.New(t)
	assert.Equal(expectedRequest, invocation.GetRequest())
}

func TestInvocation_WithPayload_Pass(t *testing.T) {
	testState := NewTestingMock(t)
	invocation := &Invocation{
		requestContent: []byte("foo"),
		testState:      testState,
	}

	invocation.WithPayload([]byte("foo"))

	testState.assertDidNotFailed()

}

func TestInvocation_WithPayload_Fail(t *testing.T) {
	testState := NewTestingMock(t)
	invocation := &Invocation{
		requestContent: []byte("foo"),
		testState:      testState,
	}

	invocation.WithPayload([]byte("bar"))

	testState.assertFailedWithError()
}

func TestInvocation_GetRequestContent(t *testing.T) {
	expectedRequestContent := []byte{42}
	invocation := &Invocation{
		requestContent: expectedRequestContent,
	}

	assert := assertions.New(t)
	assert.Equal(expectedRequestContent, invocation.GetRequestContent())
}

func TestInvocation_ReadRequestContentAsString(t *testing.T) {
	invocation := &Invocation{
		requestContent: []byte("foo"),
	}

	assert := assertions.New(t)
	assert.Equal("foo", invocation.ReadRequestContentAsString())
}

func TestInvocation_ReadRequestContentAsJson(t *testing.T) {
	invocation := &Invocation{
		requestContent: []byte(`{"foo":"bar"}`),
	}
	json := struct {
		Foo string `json:"foo"`
	}{}

	_ = invocation.ReadRequestContentAsJson(&json)

	assert := assertions.New(t)
	assert.Equal("bar", json.Foo)
}
