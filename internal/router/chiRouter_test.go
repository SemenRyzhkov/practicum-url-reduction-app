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
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
)

const (
	expectedURL       = "https://dzen.ru/?yredirect=true"
	expectedReduceURL = "http://localhost:8080/1f67218b4bfbc6af9e52d502c3e5ef3d"
)

var (
	reduceSeveralURLRequest = []entity.URLWithIDRequest{
		{CorrelationID: "test1", OriginalURL: "yandex1.ru"},
		{CorrelationID: "test2", OriginalURL: "yandex2.ru"},
	}
	reduceSeveralURLResponse = []entity.URLWithIDResponse{
		{
			CorrelationID: "test1",
			ShortURL:      "http://localhost:8080/b6ad61b613c33a6d62e6d14198e465b8",
		},
		{
			CorrelationID: "test2",
			ShortURL:      "http://localhost:8080/50754651b2f907807de0b789248f1f1b",
		},
	}
)

func setupTestServer() *httptest.Server {
	utils.LoadEnvironments("../../.env")
	repo := utils.CreateMemoryOrFileRepository(utils.GetFilePath())
	us := urlservice.New(repo)
	cs := cookieservice.New(utils.GetKey())
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

func testSeveralJSONRequest(t *testing.T, ts *httptest.Server) *http.Request {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reduceSeveralURLRequest)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten/batch", &buf)
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

func TestNewRouterGetAll(t *testing.T) {
	expectedURLsList := []entity.FullURL{{ShortURL: expectedReduceURL, OriginalURL: expectedURL}}
	ts := setupTestServer()
	defer ts.Close()

	req := testJSONRequest(t, ts)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	req = testRequest(t, ts, "GET", "/api/user/urls", "")
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actualURLsList []entity.FullURL
	err = json.NewDecoder(resp.Body).Decode(&actualURLsList)
	require.NoError(t, err)
	assert.Equal(t, expectedURLsList, actualURLsList)
	defer resp.Body.Close()
	testutils.AfterTest()
}

func TestNewRouterReducingSeveralURLToJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	req := testSeveralJSONRequest(t, ts)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	var actualResponse []entity.URLWithIDResponse
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, reduceSeveralURLResponse, actualResponse)
	defer resp.Body.Close()

	testutils.AfterTest()
}
