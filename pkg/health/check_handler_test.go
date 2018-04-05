package health_test

import (
	"adv-caja-x-api/pkg/health"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

const jsonContentType = "application/json"
const talapiBasePath = "/v1/api"

type CheckHandlerSuite struct {
	suite.Suite
	talapiRequest *http.Request
	talapi        *httptest.Server
	talapiBaseUrl *url.URL
	talapiHandler *http.Handler
	recorder      *httptest.ResponseRecorder
}

func (s *CheckHandlerSuite) SetupTest() {
	s.talapiRequest = nil
	s.talapi = nil
	s.talapiBaseUrl = nil
	s.talapiHandler = nil
	s.recorder = httptest.NewRecorder()
}

func (s *CheckHandlerSuite) AfterTest(suiteName, testName string) {
	if s.talapi != nil {
		s.talapi.Close()
	}
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(CheckHandlerSuite))
}

func (s *CheckHandlerSuite) StartTalapiServer(handler http.Handler) {
	s.talapi = httptest.NewServer(handler)
	s.talapiBaseUrl, _ = url.Parse(s.talapi.URL)
	s.talapiBaseUrl.Path = path.Join(s.talapiBaseUrl.Path, talapiBasePath)
}

func (s *CheckHandlerSuite) BuildTalapiHandler(contentType, body string, status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		fmt.Fprintln(w, body)
		s.talapiRequest = r
	})
}
func (s *CheckHandlerSuite) CommonResponseAssertions(expectedBodyContent string) {
	assert := s.Assert()
	assert.Contains(s.recorder.Body.String(), expectedBodyContent)
	assert.Equal(http.StatusOK, s.recorder.Result().StatusCode)
	assert.Equal("application/json; charset=utf-8", s.recorder.HeaderMap.Get("Content-Type"))
	if s.talapiRequest != nil {
		assert.Equal(talapiBasePath+"/health", s.talapiRequest.URL.RequestURI())
	}
}

func (s *CheckHandlerSuite) TestCheckHandlerSuccess() {
	s.StartTalapiServer(s.BuildTalapiHandler(jsonContentType, `{"message":"Cool!"}`, http.StatusOK))
	handler := health.CheckHandler(*s.talapiBaseUrl, "123")
	req := httptest.NewRequest("GET", "http://base/health", nil)
	handler.ServeHTTP(s.recorder, req)
	s.CommonResponseAssertions(`{"version":"123","talapi":"Cool!"}`)
}

func (s *CheckHandlerSuite) TestCheckHandlerTalapiConnectError() {
	fakeUrl, _ := url.Parse("http://this.domain.is.bogus")
	handler := health.CheckHandler(*fakeUrl, "123")
	req := httptest.NewRequest("GET", "http://base/health", nil)
	handler.ServeHTTP(s.recorder, req)
	s.CommonResponseAssertions("An error occurred calling Talapi")
}

func (s *CheckHandlerSuite) TestCheckHandlerTalapiBogusUrl() {
	handler := health.CheckHandler(url.URL{Scheme: "bogus"}, "123")
	req := httptest.NewRequest("GET", "http://base/health", nil)
	handler.ServeHTTP(s.recorder, req)
	s.CommonResponseAssertions("An error occurred calling Talapi")
}

func (s *CheckHandlerSuite) TestCheckHandlerTalapiNonSuccessResponse() {
	s.StartTalapiServer(s.BuildTalapiHandler(jsonContentType, `{"message":"Bummer"}`, http.StatusInternalServerError))
	handler := health.CheckHandler(*s.talapiBaseUrl, "123")
	req := httptest.NewRequest("GET", "http://base/health", nil)
	handler.ServeHTTP(s.recorder, req)
	s.CommonResponseAssertions("Talapi responded with status: 500 Internal Server Error")
}

func (s *CheckHandlerSuite) TestCheckHandlerTalapiBadJson() {
	s.StartTalapiServer(s.BuildTalapiHandler(jsonContentType, `<blah>`, http.StatusOK))
	handler := health.CheckHandler(*s.talapiBaseUrl, "123")
	req := httptest.NewRequest("GET", "http://base/health", nil)
	handler.ServeHTTP(s.recorder, req)
	s.CommonResponseAssertions("Failed to decode talapi healthcheck response")
}
