package integration_test_test

import (
	"fmt"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/config"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestSuite struct {
	suite.Suite
	app    *app.App
	server *httptest.Server
}

func (s *TestSuite) SetupSuite() {
	var err error
	s.app, err = app.New(config.Config{})
	s.Require().NoError(err)

	s.server = httptest.NewServer(s.app.HTTPServer.Handler)
}

func (s *TestSuite) TearDownSuite() {
	s.server.Close()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) Test_Handlers() {
	expectedUrl := "https://dzen.ru/?yredirect=true"
	postRequest, postErr := http.NewRequest(http.MethodPost, s.server.URL+"/", strings.NewReader(expectedUrl))
	s.Require().NoError(postErr)

	postResponse, postErr := http.DefaultClient.Do(postRequest)
	s.Require().NoError(postErr)

	s.Require().Equal(http.StatusCreated, postResponse.StatusCode)

	b, err := io.ReadAll(postResponse.Body)
	s.Require().NoError(err)
	reduceUrl := string(b[:])
	fmt.Println(reduceUrl)

	getRequest, err := http.NewRequest("GET", s.server.URL+"/XVlBz", nil)

	transport := http.Transport{}
	getResponse, err := transport.RoundTrip(getRequest)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusTemporaryRedirect, getResponse.StatusCode)

	actualUrl := getResponse.Header.Get("Location")
	s.Require().Equal(expectedUrl, actualUrl)

}
