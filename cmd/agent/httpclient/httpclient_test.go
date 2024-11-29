package httpclient_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/agent/httpclient"
	"evgen3000/go-musthave-metrics-tpl.git/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSendMetrics_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "gzip", r.Header.Get("Content-Encoding"))

		body, _ := decompressBody(t, r.Body)
		expectedHash := crypto.GenerateHash(body, "test_key")
		assert.Equal(t, expectedHash, r.Header.Get("HashSHA256"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := httpclient.NewHTTPClient(server.Listener.Addr().String(), "test_key")

	data := []byte(`{"metric":"value"}`)
	client.SendMetrics(data)
}

func TestSendMetrics_RequestSendError(t *testing.T) {
	client := httpclient.NewHTTPClient("http://localhost:12345", "test_key") // Невалидный порт

	data := []byte(`{"metric":"value"}`)

	// Захватываем вывод
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	client.SendMetrics(data)

	// Считываем и возвращаем вывод
	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldStdout
	_, _ = io.Copy(&buf, r)

	// Проверяем, что вывод содержит сообщение об ошибке
	assert.Contains(t, buf.String(), "Error sending request")
}

func decompressBody(t *testing.T, body io.ReadCloser) ([]byte, error) {
	defer func() {
		err := body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	gzipReader, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := gzipReader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	return io.ReadAll(gzipReader)
}
