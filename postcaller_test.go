package httpcaller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostCaller(t *testing.T) {
	mockClient := &http.Client{
		Transport: &mockTransport{},
	}

	t.Run("Successful POST request with default headers", func(t *testing.T) {
		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
			CallerOptions{
				DefaultHeaders: map[string]string{
					"Authorization": "Bearer default_token_here",
				},
			},
		)

		req := map[string]interface{}{
			"test": "data",
		}

		ctx := context.Background()
		res, err := caller.Post(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "data", res["test"])
	})

	t.Run("Successful POST request with additional headers", func(t *testing.T) {
		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
			CallerOptions{
				DefaultHeaders: map[string]string{
					"Authorization": "Bearer default_token_here",
				},
			},
		)

		req := map[string]interface{}{
			"test": "data",
		}

		optional := CallOption{
			Header: map[string]string{
				"Custom-Header": "custom_value",
			},
		}

		ctx := context.Background()
		res, err := caller.Post(ctx, req, optional)
		assert.NoError(t, err)
		assert.Equal(t, "data", res["test"])
	})

	t.Run("Successful POST request with path parameters", func(t *testing.T) {
		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test/:id",
			CallerOptions{
				DefaultHeaders: map[string]string{
					"Authorization": "Bearer default_token_here",
				},
			},
		)

		req := map[string]interface{}{
			"test": "data",
		}

		optional := CallOption{
			PathParam: map[string]string{
				"id": "123",
			},
		}

		ctx := context.Background()
		res, err := caller.Post(ctx, req, optional)
		assert.NoError(t, err)
		assert.Equal(t, "data", res["test"])
	})

	t.Run("Successful POST request with base success response validation", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				mockResponseBody: `{"test": "data", "status": "success"}`,
			},
		}

		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
			CallerOptions{
				DefaultHeaders: map[string]string{
					"Authorization": "Bearer default_token_here",
				},
				BaseSuccessResponse: map[string]interface{}{
					"status": "success",
				},
			},
		)

		req := map[string]interface{}{
			"test": "data",
		}

		ctx := context.Background()
		res, err := caller.Post(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "data", res["test"])
	})

	t.Run("Failed POST request due to unmatched base success response", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				mockResponseBody: `{"test": "data", "status": "error"}`,
			},
		}

		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
			CallerOptions{
				DefaultHeaders: map[string]string{
					"Authorization": "Bearer default_token_here",
				},
				BaseSuccessResponse: map[string]interface{}{
					"status": "success",
				},
			},
		)

		req := map[string]interface{}{
			"test": "data",
		}

		ctx := context.Background()
		_, err := caller.Post(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsuccessful response for key status")
	})

	t.Run("Failed POST request due to network error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				networkError: true,
			},
		}

		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		req := map[string]interface{}{
			"test": "data",
		}

		ctx := context.Background()
		_, err := caller.Post(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post request error")
	})

	t.Run("Failed POST request due to JSON marshal error", func(t *testing.T) {
		caller := NewPostCaller[chan int, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		req := make(chan int)

		ctx := context.Background()
		_, err := caller.Post(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "marshal request error")
	})

	t.Run("Failed POST request due to HTTP request creation error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				requestCreationError: true,
			},
		}

		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		req := map[string]interface{}{
			"test": "data",
		}

		ctx := context.Background()
		_, err := caller.Post(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create request error")
	})

	t.Run("Failed POST request due to response read error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				readError: true,
			},
		}

		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		req := map[string]interface{}{
			"test": "data",
		}

		ctx := context.Background()
		_, err := caller.Post(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "read response error")
	})

	t.Run("Failed POST request due to response unmarshal error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				unmarshalError: true,
			},
		}

		caller := NewPostCaller[map[string]interface{}, map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		req := map[string]interface{}{
			"test": "data",
		}

		ctx := context.Background()
		_, err := caller.Post(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal response error")
	})
}

// Mock transport for simulating different scenarios
type mockTransport struct {
	networkError         bool
	requestCreationError bool
	readError            bool
	unmarshalError       bool
	mockResponseBody     string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.networkError {
		return nil, fmt.Errorf("network error")
	}

	if m.requestCreationError {
		return nil, fmt.Errorf("create request error")
	}

	if m.readError {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(&errorReader{}),
		}, nil
	}

	if m.unmarshalError {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("{invalid json}")),
		}, nil
	}

	if m.mockResponseBody != "" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(m.mockResponseBody)),
		}, nil
	}

	body := `{"test": "data", "status": "error"}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}, nil
}

// Error reader to simulate read errors
type errorReader struct{}

func (*errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}

func (*errorReader) Close() error {
	return nil
}
