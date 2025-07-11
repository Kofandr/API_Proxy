package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockClient struct {
	resp *http.Response
	err  error
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func setupHandler() *Handler {
	logger := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return New(logger)
}

func TestHandler_Proxy(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		mockStatus     int
		mockStatusText string
		mockBody       string
		mockHeaders    map[string]string
		mockErr        error
		expectedStatus int
		expectedBody   string
		expectedHeader map[string]string
	}{
		{
			name:           "Empty path",
			method:         http.MethodGet,
			path:           "/api/",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Empty path after /api/\n",
		},
		{
			name:           "Invalid endpoint",
			method:         http.MethodGet,
			path:           "/api/users",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Invalid endpoint\n",
		},
		{
			name:           "Invalid post ID",
			method:         http.MethodGet,
			path:           "/api/posts/abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid post ID\n",
		},
		{
			name:           "Invalid path too long",
			method:         http.MethodGet,
			path:           "/api/posts/1/extra",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Invalid path\n",
		},
		{
			name:           "GET /api/posts success",
			method:         http.MethodGet,
			path:           "/api/posts",
			mockStatus:     http.StatusOK,
			mockStatusText: "200 OK",
			mockBody:       `[{"id":1,"title":"test"}]`,
			mockHeaders:    map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"title":"test"}]`,
			expectedHeader: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:           "POST /api/posts success",
			method:         http.MethodPost,
			path:           "/api/posts",
			body:           `{"title":"foo","body":"bar","userId":1}`,
			mockStatus:     http.StatusCreated,
			mockStatusText: "201 Created",
			mockBody:       `{"id":101,"title":"foo"}`,
			mockHeaders:    map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":101,"title":"foo"}`,
			expectedHeader: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:           "PUT /api/posts not allowed",
			method:         http.MethodPut,
			path:           "/api/posts",
			body:           `{"title":"foo"}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
			expectedHeader: map[string]string{"Allow": "GET, POST"},
		},
		{
			name:           "PUT /api/posts/1 success",
			method:         http.MethodPut,
			path:           "/api/posts/1",
			body:           `{"id":1,"title":"foo"}`,
			mockStatus:     http.StatusOK,
			mockStatusText: "200 OK",
			mockBody:       `{"id":1,"title":"foo"}`,
			mockHeaders:    map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"foo"}`,
			expectedHeader: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:           "PATCH /api/posts/1 success",
			method:         http.MethodPatch,
			path:           "/api/posts/1",
			body:           `{"title":"updated"}`,
			mockStatus:     http.StatusOK,
			mockStatusText: "200 OK",
			mockBody:       `{"id":1,"title":"updated"}`,
			mockHeaders:    map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"updated"}`,
			expectedHeader: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:           "DELETE /api/posts/1 success",
			method:         http.MethodDelete,
			path:           "/api/posts/1",
			mockStatus:     http.StatusNoContent,
			mockStatusText: "204 No Content",
			mockBody:       "",
			mockHeaders:    map[string]string{},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "Upstream error 404",
			method:         http.MethodGet,
			path:           "/api/posts/999",
			mockStatus:     http.StatusNotFound,
			mockStatusText: "404 Not Found",
			mockBody:       `{"error":"not found"}`,
			mockHeaders:    map[string]string{},
			expectedStatus: http.StatusBadGateway,
			expectedBody:   "Upstream error: 404 Not Found\n",
		},
		{
			name:           "Upstream client error",
			method:         http.MethodGet,
			path:           "/api/posts",
			mockErr:        http.ErrHandlerTimeout,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Server Error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setupHandler()

			// Настраиваем мок клиента
			mockResp := &http.Response{
				StatusCode: tt.mockStatus,
				Status:     tt.mockStatusText,
				Body:       io.NopCloser(strings.NewReader(tt.mockBody)),
				Header:     make(http.Header),
			}
			for k, v := range tt.mockHeaders {
				mockResp.Header.Set(k, v)
			}
			handler.client = &mockClient{resp: mockResp, err: tt.mockErr}

			body := strings.NewReader(tt.body)
			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			rr := httptest.NewRecorder()

			handler.Proxy(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch")

			require.Equal(t, tt.expectedBody, rr.Body.String(), "Response body mismatch")

			for k, v := range tt.expectedHeader {
				require.Equal(t, v, rr.Header().Get(k), "Header %q mismatch", k)
			}
		})
	}
}
