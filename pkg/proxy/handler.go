package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/oscarbc96/agbridge/pkg/log"
	"github.com/samber/lo"
)

type Handler struct {
	Path       string
	ResourceID string
	RestAPIID  string
	Methods    []string
	Config     aws.Config
}

func defaultHandleRequest(w http.ResponseWriter, r *http.Request, handlerMapping map[*regexp.Regexp]Handler) {
	start := time.Now()

	path := getPath(r.URL)

	var handler *Handler
	for pattern, h := range handlerMapping {
		if pattern.MatchString(path) {
			handler = &h
			break
		}
	}

	if handler == nil {
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

	log.Debug("Sending request to API Gateway",
		log.String("url", r.URL.String()),
		log.String("method", r.Method),
		log.Any("headers", r.Header),
		log.String("body", string(body)),
	)

	client := apigateway.NewFromConfig(handler.Config)
	resp, err := client.TestInvokeMethod(
		r.Context(),
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

	log.Debug("Received response from API Gateway",
		log.Int("status_code", int(resp.Status)),
		log.Any("headers", resp.Headers),
		log.Any("multi_value_headers", resp.MultiValueHeaders),
		log.String("body", aws.ToString(resp.Body)),
	)

	// Copy the headers from test-invoke response to the proxy response
	//for key, value := range resp.Headers {
	//	w.Header().Set(key, value)
	//}
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

	// Copy HTTP status code from test-invoke response to the proxy response
	w.WriteHeader(int(resp.Status))

	log.Info(
		r.URL.String(),
		log.String("method", r.Method),
		log.Int("status_code", int(resp.Status)),
		log.Duration("elapsed_ms", time.Since(start)),
	)
}

func getPath(u *url.URL) string {
	uCopy := *u
	uCopy.RawQuery = ""
	return strings.TrimRight(uCopy.String(), "/")
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
