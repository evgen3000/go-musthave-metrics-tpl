package logger_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"evgen3000/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		logger.InitLogger()
	}, "InitLogger should not panic")

	assert.NotNil(t, logger.GetLogger(), "Logger instance should not be nil after InitLogger is called")
}

func TestLoggingMiddleware(t *testing.T) {
	logger.InitLogger()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	wrappedHandler := logger.LoggingMiddleware(handler)
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("request body"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code to be 200")
	assert.Equal(t, "OK", string(body), "Expected response body to be 'OK'")
}

//func TestLoggingResponseWriter(t *testing.T) {
//	response := httptest.NewRecorder()
//	responseData := &logger.responseData{}
//	lw := logger.loggingResponseWriter{
//		ResponseWriter: response,
//		responseData:   responseData,
//	}
//
//	body := []byte("response body")
//	size, err := lw.Write(body)
//
//	assert.NoError(t, err, "Expected no error writing response body")
//	assert.Equal(t, len(body), size, "Expected written size to match body length")
//	assert.Equal(t, len(body), lw.responseData.size, "Expected responseData size to match body length")
//}

func TestLoggingMiddleware_RequestBodyLogging(t *testing.T) {
	logger.InitLogger()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	wrappedHandler := logger.LoggingMiddleware(handler)
	reqBody := "sample request body"
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Expected status code to be 200")
}
