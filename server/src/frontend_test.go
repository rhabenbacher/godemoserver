package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestFrontendHandler(t *testing.T) {

	tests := []struct {
		name       string
		setenv     func(apiHost string, api_Port string)
		handler    func(w http.ResponseWriter, r *http.Request)
		want       string
		statusCode int
	}{

		{
			name: "expect time output and status OK",
			setenv: func(apiHost string, apiPort string) {
				os.Setenv("API_HOST", apiHost)
				os.Setenv("API_PORT", apiPort)
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "{\"Name\":\"Testserver\",\"Time\":\"07.04.2021 13:01:01.001\"}")
			},
			want:       "Host Testserver sent time: 07.04.2021 13:01:01.001",
			statusCode: http.StatusOK,
		},
		{
			name: "expect error environment not set",
			setenv: func(apiHost string, apiPort string) {
				os.Setenv("API_HOST", "")
				os.Setenv("API_PORT", "")
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "")
			},
			want:       "Cannot connect to API",
			statusCode: http.StatusInternalServerError,
		},

		{
			name: "handle internal api error",
			setenv: func(apiHost string, apiPort string) {
				os.Setenv("API_HOST", apiHost)
				os.Setenv("API_PORT", apiPort)
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "internal error", http.StatusInternalServerError)
			},
			want:       "API responded with status 500",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, singletest := range tests {
		t.Run(singletest.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(singletest.handler))
			defer ts.Close()
			u, _ := url.Parse(ts.URL)
			singletest.setenv(u.Hostname(), u.Port())

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			responseRecorder := httptest.NewRecorder()
			noLog.showTime(responseRecorder, request)
			if responseRecorder.Code != singletest.statusCode {
				t.Errorf("Want status %d, got %d", singletest.statusCode, responseRecorder.Code)
			}

			if !strings.Contains(responseRecorder.Body.String(), singletest.want) {
				t.Errorf("expected in body:\n%v", singletest.want)
				t.Errorf("got in body:\n%v", responseRecorder.Body.String())
			}

		})
	}
}
