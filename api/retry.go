package api

import "net/http"

// IsRetryHeader marks a request that retries a previously failed request. GitHub's load
// balancers use it to identify retry traffic so that it can be shed first during an incident.
// It must only be sent on retry attempts, never on the first attempt of a request.
const IsRetryHeader = "X-Is-Retry"

// WithIsRetry sets IsRetryHeader on the request when isRetry is true and does nothing otherwise,
// so a retry loop can pass it on every attempt.
func WithIsRetry(isRetry bool) RequestOption {
	return func(cfg *requestConfig) {
		if isRetry {
			cfg.headers[IsRetryHeader] = "true"
		}
	}
}

// SetIsRetry is the equivalent of WithIsRetry for a request the caller built itself, such as one
// sent with DoRequest.
func SetIsRetry(req *http.Request, isRetry bool) {
	if isRetry {
		req.Header.Set(IsRetryHeader, "true")
	}
}
