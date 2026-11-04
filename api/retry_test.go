package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithIsRetry(t *testing.T) {
	tests := []struct {
		name    string
		isRetry bool
		want    []string
	}{
		{
			name:    "first attempt",
			isRetry: false,
			want:    nil,
		},
		{
			name:    "retry attempt",
			isRetry: true,
			want:    []string{"true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &httpmock.Registry{}
			defer reg.Verify(t)
			client := newTestClient(reg)

			reg.Register(httpmock.REST(http.MethodGet, "user"), httpmock.StringResponse(`{}`))
			reg.Register(httpmock.REST(http.MethodGet, "user"), httpmock.StringResponse(`{}`))

			err := client.REST("github.com", http.MethodGet, "user", nil, nil, WithIsRetry(tt.isRetry))
			require.NoError(t, err)
			_, err = client.RESTWithNext("github.com", http.MethodGet, "user", nil, nil, WithIsRetry(tt.isRetry))
			require.NoError(t, err)

			require.Len(t, reg.Requests, 2)
			for _, req := range reg.Requests {
				if tt.want == nil {
					assert.Empty(t, req.Header.Values(IsRetryHeader))
				} else {
					assert.Equal(t, tt.want, req.Header.Values(IsRetryHeader))
				}
			}
		})
	}
}

func TestSetIsRetry(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	SetIsRetry(req, false)
	assert.Empty(t, req.Header.Values(IsRetryHeader))

	SetIsRetry(req, true)
	assert.Equal(t, []string{"true"}, req.Header.Values(IsRetryHeader))
}
