package health_test

import (
	"damo/pkg/health"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

const jsonContentType = "application/json"

type CheckHandlerSuite struct {
	suite.Suite
	recorder *httptest.ResponseRecorder
}

func (s *CheckHandlerSuite) SetupTest() {
	s.recorder = httptest.NewRecorder()
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(CheckHandlerSuite))
}
func (s *CheckHandlerSuite) CommonResponseAssertions(expectedBodyContent string) {
	assert := s.Assert()
	assert.Contains(s.recorder.Body.String(), expectedBodyContent)
	assert.Equal(http.StatusOK, s.recorder.Result().StatusCode)
	assert.Equal("application/json; charset=utf-8", s.recorder.HeaderMap.Get("Content-Type"))
}

func (s *CheckHandlerSuite) TestCheckHandlerSuccess() {
	handler := health.CheckHandler("123")
	req := httptest.NewRequest("GET", "http://base/health", nil)
	handler.ServeHTTP(s.recorder, req)
	s.CommonResponseAssertions(`{"version":"123"}`)
}
