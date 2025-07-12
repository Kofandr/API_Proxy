package pathbuilder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildTargetURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		expected string
		errMsg   string
	}{
		{
			name:     "valid path without ID",
			baseURL:  "https://api.example.com",
			path:     "/api/posts",
			expected: "https://api.example.com/posts",
		},
		{
			name:     "valid path with numeric ID",
			baseURL:  "https://api.example.com",
			path:     "/api/posts/123",
			expected: "https://api.example.com/posts/123",
		},
		{
			name:     "path without /api/ prefix",
			baseURL:  "https://api.example.com",
			path:     "posts/456",
			expected: "https://api.example.com/posts/456",
		},
		{
			name:     "empty path after trim",
			baseURL:  "https://api.example.com",
			path:     "/api/",
			expected: "",
			errMsg:   "empty path after /api/",
		},
		{
			name:     "invalid endpoint",
			baseURL:  "https://api.example.com",
			path:     "/api/users",
			expected: "",
			errMsg:   "invalid endpoint",
		},
		{
			name:     "non-numeric ID",
			baseURL:  "https://api.example.com",
			path:     "/api/posts/abc",
			expected: "",
			errMsg:   "invalid post ID",
		},
		{
			name:     "too many path parts",
			baseURL:  "https://api.example.com",
			path:     "/api/posts/123/comments",
			expected: "",
			errMsg:   "invalid path",
		},
		{
			name:     "zero ID",
			baseURL:  "https://api.example.com",
			path:     "/api/posts/0",
			expected: "https://api.example.com/posts/0",
		},
		{
			name:     "large numeric ID",
			baseURL:  "https://api.example.com",
			path:     "/api/posts/18446744073709551615",
			expected: "https://api.example.com/posts/18446744073709551615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildTargetURL(tt.baseURL, tt.path)

			if tt.errMsg != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.errMsg)
				require.Empty(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
