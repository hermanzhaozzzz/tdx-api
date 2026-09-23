package tdx

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

type bjCodeTestTransport func(*http.Request) (*http.Response, error)

func (transport bjCodeTestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

func TestGetBjCodesHasRequestDeadlineAndReturnsError(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	calls := 0
	http.DefaultTransport = bjCodeTestTransport(func(request *http.Request) (*http.Response, error) {
		calls++
		deadline, ok := request.Context().Deadline()
		if !ok {
			t.Fatal("BSE request has no deadline")
		}
		if remaining := time.Until(deadline); remaining <= 0 || remaining > 10*time.Second {
			t.Fatalf("unexpected deadline: %s", remaining)
		}
		return nil, context.DeadlineExceeded
	})
	_, err := GetBjCodes()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline error was not returned: %v", err)
	}
	if calls != 1 {
		t.Fatalf("got %d requests, want one without retry", calls)
	}
}
