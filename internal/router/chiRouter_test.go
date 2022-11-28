package router

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/testutils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookieService"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlService"
)

const (
	expectedURL       = "https://dzen.ru/?yredirect=true"
	expectedReduceURL = "http://localhost:8080/1f67218b4bfbc6af9e52d502c3e5ef3d"
)

func setupTestServer() *httptest.Server {
	utils.LoadEnvironments("../../.env")
	repo := utils.CreateRepository(utils.GetFilePath())
	us := urlService.NewURLService(repo)
	cs := cookieService.New(utils.GetKey())
	h := handlers.NewHandler(us, cs)
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

func testJSONRequest(t *testing.T, ts *httptest.Server) *http.Request {
	request := entity.URLRequest{URL: expectedURL}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", &buf)
	require.NoError(t, err)
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
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	transport := http.Transport{}
	resp, err = transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	actualURL := resp.Header.Get("Location")
	assert.Equal(t, expectedURL, actualURL)
	defer resp.Body.Close()
	testutils.AfterTest()
}

func TestNewRouterReducingJSON(t *testing.T) {
	expectedResponse := entity.URLResponse{Result: "http://localhost:8080/1f67218b4bfbc6af9e52d502c3e5ef3d"}

	ts := setupTestServer()
	defer ts.Close()

	req := testJSONRequest(t, ts)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	var actualResponse entity.URLResponse
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
	defer resp.Body.Close()

	req = testRequest(t, ts, "GET", "/1f67218b4bfbc6af9e52d502c3e5ef3d", "")
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	transport := http.Transport{}
	resp, err = transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	actualURL := resp.Header.Get("Location")
	assert.Equal(t, expectedURL, actualURL)
	defer resp.Body.Close()
	testutils.AfterTest()
}
