package httputil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestStruct struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

func TestBindJSON(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		target        interface{}
		expectedError bool
	}{
		{
			name:          "Valid JSON",
			payload:       `{"name": "John", "age": 30, "country": "USA"}`,
			target:        &TestStruct{},
			expectedError: false,
		},
		{
			name:          "Invalid JSON Format",
			payload:       `{"name": "John", "age": 30, "country": "USA"`,
			target:        &TestStruct{},
			expectedError: true,
		},
		{
			name:          "Unknown Field",
			payload:       `{"name": "John", "age": 30, "country": "USA", "extra": "field"}`,
			target:        &TestStruct{},
			expectedError: true,
		},
		{
			name:          "Empty Body",
			payload:       "",
			target:        &TestStruct{},
			expectedError: true,
		},
		{
			name:          "Null JSON",
			payload:       "null",
			target:        &TestStruct{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.payload))
			w := httptest.NewRecorder()

			err := BindJSON(w, req, tt.target)

			if (err != nil) != tt.expectedError {
				t.Errorf("BindJSON() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name           string
		httpCode       int
		response       interface{}
		expectedError  bool
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "Valid Response",
			httpCode:       http.StatusOK,
			response:       TestStruct{Name: "John", Age: 30, Country: "USA"},
			expectedBody:   `{"name":"John","age":30,"country":"USA"}`,
			expectedHeader: "application/json",
		},
		{
			name:           "Nil Response",
			httpCode:       http.StatusNoContent,
			response:       nil,
			expectedBody:   "",
			expectedHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := JSON(w, tt.httpCode, tt.response)

			if (err != nil) != tt.expectedError {
				t.Errorf("JSON() error = %v, expectedError %v", err, tt.expectedError)
			}

			if w.Code != tt.httpCode {
				t.Errorf("JSON() status code = %v, want %v", w.Code, tt.httpCode)
			}

			if contentType := w.Header().Get("Content-Type"); contentType != tt.expectedHeader {
				t.Errorf("JSON() Content-Type = %v, want %v", contentType, tt.expectedHeader)
			}

			if !tt.expectedError && tt.response != nil {
				got := strings.TrimSpace(w.Body.String())
				want := strings.TrimSpace(tt.expectedBody)
				if got != want {
					t.Errorf("JSON() body = %v, want %v", got, want)
				}
			}
		})
	}
}
