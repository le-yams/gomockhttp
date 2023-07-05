package mockhttp

type CallVerifier struct {
	api  *APIMock
	call *HTTPCall
}

func (verifier *CallVerifier) HasBeenCalled(expectedCallsCount int) []*Invocation {
	invocations := verifier.api.invocations[*verifier.call]
	actualCallsCount := len(invocations)
	if actualCallsCount != expectedCallsCount {
		verifier.api.testState.Fatalf("got %d http calls but was expecting %d\n", actualCallsCount, expectedCallsCount)
	}
	return invocations
}

func (verifier *CallVerifier) HasBeenCalledOnce() *Invocation {
	return verifier.HasBeenCalled(1)[0]
}
