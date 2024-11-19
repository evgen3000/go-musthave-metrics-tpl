package httpclient_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/agent/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClient_SendMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/updates/", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "gzip", r.Header.Get("Content-Encoding"))

		// Check that the body is correctly gzipped
		gzipReader, err := gzip.NewReader(r.Body)
		require.NoError(t, err)
		defer func(gzipReader *gzip.Reader) {
			err := gzipReader.Close()
			if err != nil {
				fmt.Println("Cant close gzip reader")
			}
		}(gzipReader)

		var buf bytes.Buffer
		_, err = buf.ReadFrom(gzipReader)
		require.NoError(t, err)

		expected := `[{"id":"testMetric","type":"gauge","value":42.42}]`
		assert.JSONEq(t, expected, buf.String())

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := httpclient.NewHTTPClient(server.Listener.Addr().String(), "123")

	data := `[{"id":"testMetric","type":"gauge","value":42.42}]`
	client.SendMetrics([]byte(data))
}
