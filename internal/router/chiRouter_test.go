package router

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	expectedUrl       = "https://dzen.ru/?yredirect=true"
	expectedReduceUrl = "http://localhost:8080/XVlBz"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path, body string) *http.Request {
	var req *http.Request
	var err error
	if method == "GET" {
		req, err = http.NewRequest(method, ts.URL+path, nil)
		require.NoError(t, err)
	} else {
		req, err = http.NewRequest(method, ts.URL+path, strings.NewReader(body))
		require.NoError(t, err)
	}
	return req
}

func TestNewRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := testRequest(t, ts, "POST", "/", expectedUrl)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	actualReduceUrl, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedReduceUrl, string(actualReduceUrl))

	req = testRequest(t, ts, "GET", "/XVlBz", "")
	transport := http.Transport{}
	resp, err = transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	actualUrl := resp.Header.Get("Location")
	assert.Equal(t, expectedUrl, actualUrl)
	defer resp.Body.Close()

}
