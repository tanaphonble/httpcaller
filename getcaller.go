package httpcaller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GetCaller[response any] struct {
	httpClient          *http.Client
	baseURL             string
	endpoint            string
	defaultHeaders      map[string]string
	baseSuccessResponse map[string]interface{}
}

func NewGetCaller[response any](
	httpClient *http.Client,
	baseURL string,
	endpoint string,
	options ...CallerOptions,
) *GetCaller[response] {
	defaultHeaders := make(map[string]string)
	baseSuccessResponse := make(map[string]interface{})

	if len(options) > 0 {
		opt := options[0]
		if opt.DefaultHeaders != nil {
			defaultHeaders = opt.DefaultHeaders
		}
		if opt.BaseSuccessResponse != nil {
			baseSuccessResponse = opt.BaseSuccessResponse
		}
	}

	return &GetCaller[response]{
		httpClient:          httpClient,
		baseURL:             baseURL,
		endpoint:            endpoint,
		defaultHeaders:      defaultHeaders,
		baseSuccessResponse: baseSuccessResponse,
	}
}

func (h *GetCaller[response]) Get(ctx context.Context, optional ...CallOption) (response, error) {
	var res response

	headers := make(map[string]string)
	for key, value := range h.defaultHeaders {
		headers[key] = value
	}

	pathParams := make(map[string]string)
	if len(optional) > 0 {
		opt := optional[0]
		if opt.Header != nil {
			for key, value := range opt.Header {
				headers[key] = value
			}
		}
		if opt.PathParam != nil {
			pathParams = opt.PathParam
		}
	}

	url := h.baseURL + "/" + h.endpoint
	for key, value := range pathParams {
		placeholder := fmt.Sprintf(":%s", key)
		url = strings.Replace(url, placeholder, value, -1)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return res, fmt.Errorf("create request error: %s", err)
	}

	for key, value := range headers {
		httpReq.Header.Set(key, value)
	}

	serverResponse, err := h.httpClient.Do(httpReq)
	if err != nil {
		return res, fmt.Errorf("get request error: %s", err)
	}
	defer serverResponse.Body.Close()

	bytesResponse, err := io.ReadAll(serverResponse.Body)
	if err != nil {
		return res, fmt.Errorf("read response error: %s", err)
	}

	err = json.Unmarshal(bytesResponse, &res)
	if err != nil {
		return res, fmt.Errorf("unmarshal response error: %s", err)
	}

	if len(h.baseSuccessResponse) > 0 {
		responseMap := make(map[string]interface{})
		if err := json.Unmarshal(bytesResponse, &responseMap); err != nil {
			return res, fmt.Errorf("unmarshal response map error: %s", err)
		}
		for key, expectedValue := range h.baseSuccessResponse {
			actualValue, exists := responseMap[key]
			if !exists || fmt.Sprintf("%v", actualValue) != fmt.Sprintf("%v", expectedValue) {
				return res, fmt.Errorf("unsuccessful response for key %s: expected %v, got %v", key, expectedValue, actualValue)
			}
		}
	}

	return res, nil
}
