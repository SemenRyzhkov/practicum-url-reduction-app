package router

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
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
)

func setupTestServerForExample() *httptest.Server {
	utils.LoadEnvironments("../../.env")
	repo := utils.CreateMemoryOrFileRepository(utils.GetFilePath())
	us := urlservice.New(repo)
	cs, _ := cookieservice.New(utils.GetKey())
	h := handlers.NewHandler(us, cs)
	router := NewRouter(h)
	return httptest.NewServer(router)
}

func testRequestForExample(ts *httptest.Server, method, path, body string) *http.Request {
	var req *http.Request

	if method == http.MethodGet {
		req, _ = http.NewRequest(method, ts.URL+path, nil)
	} else {
		req, _ = http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	}
	return req
}

func testJSONRequestForExample(ts *httptest.Server) *http.Request {
	request := entity.URLRequest{URL: expectedURL}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", &buf)
	return req
}

func testSeveralJSONRequestForExample(ts *httptest.Server) *http.Request {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reduceSeveralURLRequest)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten/batch", &buf)
	return req
}

func exampleReduceURLAndGetURLByID() {
	//настраиваем сервер
	ts := setupTestServerForExample()
	defer ts.Close()

	//выполняем запрос для сокращения url
	req := testRequestForExample(ts, "POST", "/", expectedURL)
	resp, _ := http.DefaultClient.Do(req)

	//получаем значение ответа
	actualReduceURL, _ := io.ReadAll(resp.Body)
	fmt.Printf("Reduce URL is %s", actualReduceURL)
	defer resp.Body.Close()

	//выполняем запрос для получения оригинального URL
	req = testRequestForExample(ts, "GET", "/1f67218b4bfbc6af9e52d502c3e5ef3d", "")
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	transport := http.Transport{}
	resp, _ = transport.RoundTrip(req)

	//получаем значение ответа из заголовка
	actualURL := resp.Header.Get("Location")
	fmt.Printf("Original URL is %s", actualURL)
	defer resp.Body.Close()
	testutils.AfterTest()

	// Output:
	// Reduce URL is http://localhost:8080/1f67218b4bfbc6af9e52d502c3e5ef3dOriginal URL is https://dzen.ru/?yredirect=true
}

func exampleReduceURLTOJSON() {
	//настраиваем сервер
	ts := setupTestServerForExample()
	defer ts.Close()

	//выполняем запрос для сокращения url
	req := testJSONRequestForExample(ts)
	resp, _ := http.DefaultClient.Do(req)

	var actualResponse entity.URLResponse
	json.NewDecoder(resp.Body).Decode(&actualResponse)
	defer resp.Body.Close()

	//для проверки полученного результата выполняем Get запрос
	req = testRequestForExample(ts, "GET", "/1f67218b4bfbc6af9e52d502c3e5ef3d", "")
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	transport := http.Transport{}
	resp, _ = transport.RoundTrip(req)

	//получаем значение из заголовка
	actualURL := resp.Header.Get("Location")
	fmt.Printf("Original URL is %s", actualURL)
	defer resp.Body.Close()
	testutils.AfterTest()

	// Output:
	// Original URL is https://dzen.ru/?yredirect=true
}

func exampleGetAll() {
	//настраиваем сервер
	ts := setupTestServerForExample()
	defer ts.Close()

	//выполняем запрос для сокращения url
	req := testJSONRequestForExample(ts)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//выполняем запрос для получения всех URL по userID
	req = testRequestForExample(ts, "GET", "/api/user/urls", "")
	req.Header["Cookie"] = append(req.Header["Cookie"], resp.Header.Get("Set-Cookie"))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//получаем значение из ответа
	var actualURLsList []entity.FullURL
	json.NewDecoder(resp.Body).Decode(&actualURLsList)
	fmt.Printf("Original URL is %s", actualURLsList[0].OriginalURL)
	fmt.Printf("Short URL is %s", actualURLsList[0].ShortURL)

	defer resp.Body.Close()
	testutils.AfterTest()

	// Output:
	// Original URL is https://dzen.ru/?yredirect=trueShort URL is http://localhost:8080/1f67218b4bfbc6af9e52d502c3e5ef3d
}

func exampleReduceSeveralURL() {
	//настраиваем сервер
	ts := setupTestServerForExample()
	defer ts.Close()

	//выполняем запрос для сокращения нескольких URL
	req := testSeveralJSONRequestForExample(ts)
	resp, _ := http.DefaultClient.Do(req)

	//получаем значение из ответа
	var actualResponse []entity.URLWithIDResponse
	json.NewDecoder(resp.Body).Decode(&actualResponse)
	fmt.Printf("First short URL is %s", actualResponse[0].ShortURL)
	fmt.Printf("Second short URL is %s", actualResponse[0].ShortURL)

	defer resp.Body.Close()
	testutils.AfterTest()

	// Output:
	// First short URL is http://localhost:8080/b6ad61b613c33a6d62e6d14198e465b8Second short URL is http://localhost:8080/b6ad61b613c33a6d62e6d14198e465b8
}
