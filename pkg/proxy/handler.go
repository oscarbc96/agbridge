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
	StagePath      string
	Path           string
	ResourceID     string
	RestAPIID      string
	Methods        []string
	Config         aws.Config
	StageVariables map[string]string
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

	pathWithQuery := handler.Path
	if rawQuery := r.URL.RawQuery; rawQuery != "" {
		pathWithQuery += "?" + rawQuery
	}

	log.Debug("Sending request to API Gateway" +
		"\nProxy URL: " + r.URL.String() +
		"\nResource ID: " + handler.ResourceID +
		"\nREST API ID: " + handler.RestAPIID +
		"\nMethod: " + r.Method +
		"\nURL: " + pathWithQuery +
		"\nBody: " + string(body) +
		"\nHeaders: " + fmt.Sprint(r.Header) +
		"\nStage Variables: " + fmt.Sprint(handler.StageVariables),
	)

	client := apigateway.NewFromConfig(handler.Config)
	resp, err := client.TestInvokeMethod(
		r.Context(),
		&apigateway.TestInvokeMethodInput{
			ResourceId:          &handler.ResourceID,
			RestApiId:           &handler.RestAPIID,
			HttpMethod:          &r.Method,
			PathWithQueryString: aws.String(pathWithQuery),
			Body:                aws.String(string(body)),
			MultiValueHeaders:   r.Header,
			StageVariables:      handler.StageVariables,
		},
	)
	if err != nil {
		handleError(w, r, err, "Error calling API Gateway")
		return
	}

	log.Debug("Received response from API Gateway:\n" + *resp.Log)

	// Copy the headers from test-invoke response to the proxy response
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
