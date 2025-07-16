package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestHandler_Get_Simple(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
		serverResponse func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:           "successful GET request",
			method:         http.MethodGet,
			path:           "/posts/1",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello from server",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello from server"))
			},
		},
		{
			name:           "invalid HTTP method",
			method:         http.MethodPost,
			path:           "/posts/1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid method\n",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello from server"))
			},
		},
		{
			name:           "non-existent ID returns 404",
			method:         http.MethodGet,
			path:           "/posts/999",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Not found",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/999" {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Not found"))
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Hello from server"))
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			testServer := httptest.NewServer(http.HandlerFunc(test.serverResponse))
			defer testServer.Close()

			h := &Handler{
				baseURL:     testServer.URL + "/",
				restyClient: resty.New(),
			}

			req := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()

			h.Get(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.expectedStatus, res.StatusCode)

			if test.expectedBody != "" {
				body := w.Body.String()
				assert.Equal(t, test.expectedBody, body)
			}
		})
	}
}
