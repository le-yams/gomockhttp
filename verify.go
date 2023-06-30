package mockhttp

type CallVerifier struct {
	api  *ApiMock
	call *HttpCall
}

func (verifier *CallVerifier) HasBeenCalled(expectedCallsCount int) {
	actualCallsCount := len(verifier.api.invocations[*verifier.call])
	if actualCallsCount != expectedCallsCount {
		verifier.api.testState.Fatalf("got %d http calls but was expecting %d\n", actualCallsCount, expectedCallsCount)
	}
}

func (verifier *CallVerifier) HasBeenCalledOnce() {
	verifier.HasBeenCalled(1)
}
