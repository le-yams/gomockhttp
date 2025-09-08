package mockhttp

// CallVerifier is a helper to verify invocations of a specific HTTP call
type CallVerifier struct {
	api  *APIMock
	call *HTTPCall
}

// HasBeenCalled asserts that the HTTP call has been made the expected number of times.
// It returns all invocations of the call.
func (verifier *CallVerifier) HasBeenCalled(expectedCallsCount int) []*Invocation {
	invocations := verifier.api.invocations[*verifier.call]
	actualCallsCount := len(invocations)
	if actualCallsCount != expectedCallsCount {
		verifier.api.testState.Fatalf("got %d http calls but was expecting %d\n", actualCallsCount, expectedCallsCount)
	}
	return invocations
}

// HasBeenCalledOnce asserts that the HTTP call has been made exactly once then returns the invocation.
func (verifier *CallVerifier) HasBeenCalledOnce() *Invocation {
	invocations := verifier.HasBeenCalled(1)
	if len(invocations) > 0 {
		return invocations[0]
	}
	return nil
}

// HasNotBeenCalled asserts that no HTTP call has been made.
func (verifier *CallVerifier) HasNotBeenCalled() {
	_ = verifier.HasBeenCalled(0)
}
