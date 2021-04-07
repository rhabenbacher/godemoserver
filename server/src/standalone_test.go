package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var noLog = serverLogs{
	infoLog:  log.New(ioutil.Discard, "", 0),
	errorLog: log.New(ioutil.Discard, "", 0),
}

func executeHttpGetTest(handler func(http.ResponseWriter, *http.Request), target string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, target, nil)
	responseRecorder := httptest.NewRecorder()
	handler(responseRecorder, request)
	return responseRecorder
}

func TestHomeHandler(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		want       string
		statusCode int
	}{
		{
			name:       "test / and expect 200",
			target:     "/",
			want:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "test /unknown and expect not found",
			target:     "/unknown",
			want:       "",
			statusCode: http.StatusNotFound,
		},
	}

	for _, singletest := range tests {
		t.Run(singletest.name, func(t *testing.T) {
			responseRecorder := executeHttpGetTest(noLog.home, singletest.target)
			if responseRecorder.Code != singletest.statusCode {
				t.Errorf("Want status %d, got %d", singletest.statusCode, responseRecorder.Code)
			}

		})
	}

}

func TestFibonacciHandler(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		want       string
		statusCode int
	}{

		{
			name:       "test /fibonacci and expect number",
			target:     "/fibonacci?n=100",
			want:       "<p>354224848179261915075</p>",
			statusCode: http.StatusOK,
		},
		{
			name:       "test /fibonacci without parameter",
			target:     "/fibonacci",
			want:       "n is required",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "test /fibonacci with wrong parameter",
			target:     "/fibonacci?qu=100",
			want:       "n is required",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, singletest := range tests {
		t.Run(singletest.name, func(t *testing.T) {

			responseRecorder := executeHttpGetTest(noLog.fibonacci, singletest.target)
			if responseRecorder.Code != singletest.statusCode {
				t.Errorf("Want status %d, got %d", singletest.statusCode, responseRecorder.Code)
			}
			if !strings.Contains(responseRecorder.Body.String(), singletest.want) {
				t.Errorf("%v not in response body", singletest.want)
				t.Errorf("Body:\n%v", responseRecorder.Body.String())
			}

		})
	}

}
