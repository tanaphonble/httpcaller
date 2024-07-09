package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tanaphonble/httpcaller"
)

type GetResponse struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func createSuccessGetCaller() *httpcaller.GetCaller[GetResponse] {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	options := httpcaller.CallerOptions{
		DefaultHeaders: map[string]string{
			"Authorization": "Bearer default_token_here",
			"Custom-Header": "default_value",
		},
		BaseSuccessResponse: map[string]interface{}{
			"userId": 1,
		},
	}

	return httpcaller.NewGetCaller[GetResponse](
		httpClient,
		"https://jsonplaceholder.typicode.com",
		"posts/1",
		options,
	)
}

func createFailedGetCaller() *httpcaller.GetCaller[GetResponse] {
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

	return httpcaller.NewGetCaller[GetResponse](
		httpClient,
		"https://jsonplaceholder.typicode.com",
		"posts/1",
		options,
	)
}

func successWithBaseResponse() {
	caller := createSuccessGetCaller()

	ctx := context.Background()
	res, err := caller.Get(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Success With Base Response: %+v\n", res)
}

func failedWithBaseResponseNotMatch() {
	caller := createFailedGetCaller()

	ctx := context.Background()
	res, err := caller.Get(ctx)
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
