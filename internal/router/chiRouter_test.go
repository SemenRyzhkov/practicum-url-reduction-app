package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service"
)

const (
	expectedURL       = "https://dzen.ru/?yredirect=true"
	expectedReduceURL = "http://localhost:8080/1f67218b4bfbc6af9e52d502c3e5ef3d"
)

func setupTestServer() *httptest.Server {
	repo := repositories.NewURLRepository()
	s := service.NewURLService(repo)
	h := handlers.NewHandler(s)
	router := NewRouter(h)
	return httptest.NewServer(router)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, body string) *http.Request {
	var req *http.Request
	var err error
	if method == http.MethodGet {
		req, err = http.NewRequest(method, ts.URL+path, nil)
		require.NoError(t, err)
	} else {
		req, err = http.NewRequest(method, ts.URL+path, strings.NewReader(body))
		require.NoError(t, err)
	}
	return req
}

func TestNewRouter(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	req := testRequest(t, ts, "POST", "/", expectedURL)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	actualReduceURL, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedReduceURL, string(actualReduceURL))
	defer resp.Body.Close()

	req = testRequest(t, ts, "GET", "/1f67218b4bfbc6af9e52d502c3e5ef3d", "")
	transport := http.Transport{}
	resp, err = transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	actualURL := resp.Header.Get("Location")
	assert.Equal(t, expectedURL, actualURL)
	defer resp.Body.Close()
}
