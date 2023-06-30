package mockhttp

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type Invocation struct {
	request        *http.Request
	requestContent []byte
	testState      TestingT
}

func (call *Invocation) GetRequest() *http.Request {
	return call.request
}

func (call *Invocation) GetRequestContent() []byte {
	return call.requestContent
}

func (call *Invocation) ReadRequestContentAsString() string {
	return string(call.requestContent)
}

func (call *Invocation) ReadRequestContentAsJson(obj any) error {
	return json.Unmarshal(call.requestContent, obj)
}

func (call *Invocation) WithPayload(content []byte) {
	if !reflect.DeepEqual(call.requestContent, content) {
		call.testState.Error()
	}
}
