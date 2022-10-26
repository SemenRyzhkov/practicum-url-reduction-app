package integration_test

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/config"
	"github.com/stretchr/testify/suite"
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

func (s *TestSuite) TestSuite_GetUrlById() {

	request, err := http.NewRequest(http.MethodPost, s.server.URL+"/", strings.NewReader("example"))
	s.Require().NoError(err)

	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusCreated, response.StatusCode)

}
