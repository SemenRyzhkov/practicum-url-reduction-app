package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/testutils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
)

const (
	expectedURL = "https://dzen.ru/?yredirect=true"
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
	cs, _ := cookieservice.New(utils.GetKey())
	h := NewHandler(us, cs)
	router := router.NewRouter(h)
	return httptest.NewServer(router)
}

func testRequest(ts *httptest.Server, method, path, body string) *http.Request {
	var req *http.Request

	if method == http.MethodGet {
		req, _ = http.NewRequest(method, ts.URL+path, nil)
	} else {
		req, _ = http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	}
	return req
}

func testJSONRequest(ts *httptest.Server) *http.Request {
	request := entity.URLRequest{URL: expectedURL}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", &buf)
	return req
}

func testSeveralJSONRequest(ts *httptest.Server) *http.Request {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reduceSeveralURLRequest)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten/batch", &buf)
	return req
}

func ExampleReduceURLAndGetURLByID() {
	//настраиваем сервер
	ts := setupTestServer()
	defer ts.Close()

	//выполняем запрос для сокращения url
	req := testRequest(ts, "POST", "/", expectedURL)
	resp, _ := http.DefaultClient.Do(req)

	//получаем значение ответа
	actualReduceURL, _ := io.ReadAll(resp.Body)
	fmt.Printf("Reduce URL is %s", actualReduceURL)
	defer resp.Body.Close()

	//выполняем запрос для получения оригинального URL
	req = testRequest(ts, "GET", "/1f67218b4bfbc6af9e52d502c3e5ef3d", "")
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	transport := http.Transport{}
	resp, _ = transport.RoundTrip(req)

	//получаем значение ответа из заголовка
	actualURL := resp.Header.Get("Location")
	fmt.Printf("Original URL is %s", actualURL)
	defer resp.Body.Close()
	testutils.AfterTest()

	// Output:
	// http://localhost:8080/1f67218b4bfbc6af9e52d502c3e5ef3d
	// https://dzen.ru/?yredirect=true
}

func ExampleReduceURLTOJSON() {
	//настраиваем сервер
	ts := setupTestServer()
	defer ts.Close()

	//выполняем запрос для сокращения url
	req := testJSONRequest(ts)
	resp, _ := http.DefaultClient.Do(req)

	var actualResponse entity.URLResponse
	json.NewDecoder(resp.Body).Decode(&actualResponse)
	defer resp.Body.Close()

	//для проверки полученного результата выполняем Get запрос
	req = testRequest(ts, "GET", "/1f67218b4bfbc6af9e52d502c3e5ef3d", "")
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	transport := http.Transport{}
	resp, _ = transport.RoundTrip(req)

	//получаем значение из заголовка
	actualURL := resp.Header.Get("Location")
	fmt.Printf("Original URL is %s", actualURL)
	defer resp.Body.Close()
	testutils.AfterTest()

	// Output:
	// https://dzen.ru/?yredirect=true
}

func ExampleGetAll() {
	//настраиваем сервер
	ts := setupTestServer()
	defer ts.Close()

	//выполняем запрос для сокращения url
	req := testJSONRequest(ts)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	//выполняем запрос для получения всех URL по userID
	req = testRequest(ts, "GET", "/api/user/urls", "")
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	resp, _ = http.DefaultClient.Do(req)

	//получаем значение из ответа
	var actualURLsList []entity.FullURL
	json.NewDecoder(resp.Body).Decode(&actualURLsList)
	fmt.Printf("Original URL is %s", actualURLsList[0].OriginalURL)
	fmt.Printf("Short URL is %s", actualURLsList[0].ShortURL)

	defer resp.Body.Close()
	testutils.AfterTest()

	// Output:
	// https://dzen.ru/?yredirect=true
	// http://localhost:8080/1f67218b4bfbc6af9e52d502c3e5ef3d
}

func ExampleReduceSeveralURL() {
	//настраиваем сервер
	ts := setupTestServer()
	defer ts.Close()

	//выполняем запрос для сокращения нескольких URL
	req := testSeveralJSONRequest(ts)
	resp, _ := http.DefaultClient.Do(req)

	//получаем значение из ответа
	var actualResponse []entity.URLWithIDResponse
	json.NewDecoder(resp.Body).Decode(&actualResponse)
	fmt.Printf("First short URL is %s", actualResponse[0].ShortURL)
	fmt.Printf("Second short URL is %s", actualResponse[0].ShortURL)

	defer resp.Body.Close()
	testutils.AfterTest()

	// Output:
	// http://localhost:8080/b6ad61b613c33a6d62e6d14198e465b8
	// http://localhost:8080/50754651b2f907807de0b789248f1f1b
}
