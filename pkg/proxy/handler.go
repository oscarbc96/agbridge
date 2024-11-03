package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/oscarbc96/agbridge/pkg/log"
	"github.com/samber/lo"
)

type Handler struct {
	ResourceID string
	RestAPIID  string
	Methods    []string
	Config     aws.Config
}

func defaultHandleRequest(w http.ResponseWriter, r *http.Request, handlerMapping map[string]Handler) {
	start := time.Now()

	pathWithoutQuery := getURLWithoutQuery(r.URL)
	handler, ok := handlerMapping[pathWithoutQuery]
	if !ok {
		handleError(w, r, nil, "Handler not found")
		return
	}

	if !lo.Contains(handler.Methods, r.Method) {
		handleError(w, r, nil, "Method not supported")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, r, err, "Error reading request body")
		return
	}

	client := apigateway.NewFromConfig(handler.Config)
	resp, err := client.TestInvokeMethod(
		context.Background(),
		&apigateway.TestInvokeMethodInput{
			ResourceId:          &handler.ResourceID,
			RestApiId:           &handler.RestAPIID,
			HttpMethod:          &r.Method,
			PathWithQueryString: aws.String(r.URL.String()),
			Body:                aws.String(string(body)),
			MultiValueHeaders:   r.Header,
		},
	)
	if err != nil {
		handleError(w, r, err, "Error calling API Gateway")
		return
	}

	// Copy HTTP status code from test-invoke response to the proxy response
	w.WriteHeader(int(resp.Status))
	// Copy the headers from test-invoke response to the proxy response
	for key, value := range resp.Headers {
		w.Header().Set(key, value)
	}
	for key, values := range resp.MultiValueHeaders {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	// Copy the body from test-invoke response to the proxy response
	_, err = io.Copy(w, strings.NewReader(*resp.Body))
	if err != nil {
		handleError(w, r, err, "Error copying response body")
		return
	}

	log.Info(
		r.URL.String(),
		log.String("method", r.Method),
		log.Int("status_code", int(resp.Status)),
		log.Duration("elapsed_ms", time.Since(start)),
	)
}

func getURLWithoutQuery(u *url.URL) string {
	uCopy := *u
	uCopy.RawQuery = ""
	return uCopy.String()
}

func handleError(w http.ResponseWriter, r *http.Request, err error, message string) {
	logger := log.With(
		log.String("path", r.URL.String()),
		log.String("method", r.Method),
	)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Raised from AGBridge: %s Error: %s", message, err.Error()),
			http.StatusInternalServerError,
		)
		logger.Error(message, log.Err(err))
	} else {
		http.Error(
			w,
			fmt.Sprintf("Raised from AGBridge: %s", message),
			http.StatusInternalServerError,
		)
		logger.Error(message)
	}
}
