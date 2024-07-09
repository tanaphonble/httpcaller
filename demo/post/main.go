package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tanaphonble/httpcaller"
)

type Post struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

type PostResponse struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func createSuccessHTTPCaller() *httpcaller.PostCaller[Post, PostResponse] {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	options := httpcaller.CallerOptions{
		DefaultHeaders: map[string]string{
			"Authorization": "Bearer default_token_here",
			"Custom-Header": "default_value",
		},
		BaseSuccessResponse: map[string]interface{}{
			"title": "foo",
		},
	}

	return httpcaller.NewPostCaller[Post, PostResponse](
		httpClient,
		"https://jsonplaceholder.typicode.com",
		"posts",
		options,
	)
}

func createFailedHTTPCaller() *httpcaller.PostCaller[Post, PostResponse] {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	options := httpcaller.CallerOptions{
		DefaultHeaders: map[string]string{
			"Authorization": "Bearer default_token_here",
			"Custom-Header": "default_value",
		},
		BaseSuccessResponse: map[string]interface{}{
			"title": "not match",
		},
	}

	return httpcaller.NewPostCaller[Post, PostResponse](
		httpClient,
		"https://jsonplaceholder.typicode.com",
		"posts",
		options,
	)
}

func successWithBaseResponse() {
	caller := createSuccessHTTPCaller()

	req := Post{
		Title:  "foo",
		Body:   "bar",
		UserID: 1,
	}

	optional := httpcaller.CallOption{
		Header: map[string]string{"Another-Header": "header_value"},
	}

	ctx := context.Background()
	res, err := caller.Post(ctx, req, optional)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Success With Base Response: %+v\n", res)
}

func failedWithBaseResponseNotMatch() {
	caller := createFailedHTTPCaller()

	req := Post{
		Title:  "foo",
		Body:   "bar",
		UserID: 1,
	}

	// The baseSuccessResponse expects the title to be "foo", this will fail
	optional := httpcaller.CallOption{
		Header: map[string]string{"Another-Header": "header_value"},
	}

	ctx := context.Background()
	res, err := caller.Post(ctx, req, optional)
	if err != nil {
		fmt.Println("Error (Base Response Not Match):", err)
		return
	}

	fmt.Printf("Failed With Base Response Not Match: %+v\n", res)
}

func main() {
	fmt.Println("Running successWithBaseResponse...")
	successWithBaseResponse()

	fmt.Println("\nRunning failedWithBaseResponseNotMatch...")
	failedWithBaseResponseNotMatch()
}
