package httpcaller

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCaller(t *testing.T) {
	mockClient := &http.Client{
		Transport: &mockTransport{},
	}

	t.Run("Successful GET request with default headers", func(t *testing.T) {
		caller := NewGetCaller[map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
			CallerOptions{
				DefaultHeaders: map[string]string{
					"Authorization": "Bearer default_token_here",
				},
			},
		)

		ctx := context.Background()
		res, err := caller.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "data", res["test"])
	})

	t.Run("Successful GET request with additional headers", func(t *testing.T) {
		caller := NewGetCaller[map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
			CallerOptions{
				DefaultHeaders: map[string]string{
					"Authorization": "Bearer default_token_here",
				},
			},
		)

		optional := CallOption{
			Header: map[string]string{
				"Custom-Header": "custom_value",
			},
		}

		ctx := context.Background()
		res, err := caller.Get(ctx, optional)
		assert.NoError(t, err)
		assert.Equal(t, "data", res["test"])
	})

	t.Run("Successful GET request with path parameters", func(t *testing.T) {
		caller := NewGetCaller[map[string]interface{}](
			mockClient,
			"https://example.com",
			"test/:id",
			CallerOptions{
				DefaultHeaders: map[string]string{
					"Authorization": "Bearer default_token_here",
				},
			},
		)

		optional := CallOption{
			PathParam: map[string]string{
				"id": "123",
			},
		}

		ctx := context.Background()
		res, err := caller.Get(ctx, optional)
		assert.NoError(t, err)
		assert.Equal(t, "data", res["test"])
	})

	t.Run("Successful GET request with base success response validation", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				mockResponseBody: `{"test": "data", "status": "success"}`,
			},
		}

		caller := NewGetCaller[map[string]interface{}](
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

		ctx := context.Background()
		res, err := caller.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "data", res["test"])
	})

	t.Run("Failed GET request due to unmatched base success response", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				mockResponseBody: `{"test": "data", "status": "error"}`,
			},
		}

		caller := NewGetCaller[map[string]interface{}](
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

		ctx := context.Background()
		_, err := caller.Get(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsuccessful response for key status")
	})

	t.Run("Failed GET request due to network error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				networkError: true,
			},
		}

		caller := NewGetCaller[map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		ctx := context.Background()
		_, err := caller.Get(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get request error")
	})

	t.Run("Failed GET request due to HTTP request creation error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				requestCreationError: true,
			},
		}

		caller := NewGetCaller[map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		ctx := context.Background()
		_, err := caller.Get(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create request error")
	})

	t.Run("Failed GET request due to response read error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				readError: true,
			},
		}

		caller := NewGetCaller[map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		ctx := context.Background()
		_, err := caller.Get(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "read response error")
	})

	t.Run("Failed GET request due to response unmarshal error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockTransport{
				unmarshalError: true,
			},
		}

		caller := NewGetCaller[map[string]interface{}](
			mockClient,
			"https://example.com",
			"test",
		)

		ctx := context.Background()
		_, err := caller.Get(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal response error")
	})
}
