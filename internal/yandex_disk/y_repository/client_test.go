package yrepository

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Test case for successful client creation
	t.Run("Success", func(t *testing.T) {
		// Подготовка данных для теста
		token := &Token{AccessToken: "abc123"}
		baseURL := "http://example.com"
		httpClient := &http.Client{}

		// Вызов тестируемого метода
		client, err := NewClient(token, baseURL, httpClient)

		// Проверки
		assert.NoError(t, err, "Expected no error")
		assert.Equal(t, client.httpClient, httpClient, "Expected httpClient to be set correctly")
		assert.Equal(t, client.token, token, "Expected token to be set correctly")
		assert.Equal(t, client.baseUrl.String(), baseURL, "Expected baseUrl to be set correctly")
	})

	// Test case for nil httpClient
	t.Run("NilHTTPClient", func(t *testing.T) {
		token := &Token{AccessToken: "abc123"}
		baseURL := "http://example.com"
		client, err := NewClient(token, baseURL, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if client.httpClient != http.DefaultClient {
			t.Errorf("httpClient not set to DefaultClient")
		}
	})

	// Test case for empty token
	t.Run("EmptyToken", func(t *testing.T) {
		baseURL := "http://example.com"
		_, err := NewClient(&Token{}, baseURL, http.DefaultClient)
		if err == nil {
			t.Error("expected an error, got nil")
		}
	})

	// Test case for empty baseURL
	t.Run("EmptyBaseURL", func(t *testing.T) {
		token := &Token{AccessToken: "abc123"}
		_, err := NewClient(token, "", http.DefaultClient)
		if err == nil {
			t.Error("expected an error, got nil")
		}
	})

	// Test case for invalid baseURL
	t.Run("InvalidBaseURL", func(t *testing.T) {
		token := &Token{AccessToken: "abc123"}
		baseURL := ""
		_, err := NewClient(token, baseURL, http.DefaultClient)
		if err == nil {
			t.Error("expected an error, got nil")
		}
	})
}

func TestSetRequestHeaders(t *testing.T) {

	req1, _ := http.NewRequest("GET", "http://example.com", nil)
	token := &Token{AccessToken: "abc123"}
	baseURL := "http://example.com"
	httpClient := &http.Client{}

	c, _ := NewClient(token, baseURL, httpClient)

	c.SetRequestHeaders(req1)

	assert.NotEmpty(t, req1.Header.Get("Access"), "Access header should not be empty")
	assert.NotEmpty(t, req1.Header.Get("Content-Type"), "Content-Type header should not be empty")
	assert.NotEmpty(t, req1.Header.Get("Authorization"), "Authorization header should not be empty")

	assert.Equal(t, "application/json", req1.Header.Get("Access"))
	assert.Equal(t, "application/json", req1.Header.Get("Content-Type"))
	assert.Equal(t, "OAuth abc123", req1.Header.Get("Authorization"))
}

func TestMakeRequest(t *testing.T) {
	// Test case for successful request creation
	t.Run("Success", func(t *testing.T) {
		// Подготовка данных для теста
		token := &Token{AccessToken: "abc123"}
		baseURL := "http://example.com"
		httpClient := &http.Client{}

		client, _ := NewClient(token, baseURL, httpClient)

		method := "GET"
		pathUrl := "/users"

		req, err := client.MakeRequest(method, pathUrl, nil)

		assert.NoError(t, err, "Expected no error")
		assert.Equal(t, method, req.Method, "Expected method to be set correctly")
		assert.Equal(t, "http://example.com/users", req.URL.String(), "Expected URL to be set correctly")
	})

	// // Test case for base URL containing trailing slash
	t.Run("BaseURLContainingTrailingSlash", func(t *testing.T) {
		// Подготовка данных для теста
		token := &Token{AccessToken: "abc123"}
		baseURL := "http://example.com/"
		httpClient := &http.Client{}

		client, _ := NewClient(token, baseURL, httpClient)

		method := "GET"
		pathUrl := "/users"

		req, err := client.MakeRequest(method, pathUrl, nil)

		assert.NoError(t, err, "Expected an error")
		assert.Equal(t, "http://example.com/users", req.URL.String(), "Expected URL to be set not correctly")
	})

	// Test case for invalid base URL
	t.Run("InvalidMethod", func(t *testing.T) {
		token := &Token{AccessToken: "abc123"}
		baseURL := "http://example.com/"
		httpClient := &http.Client{}

		client, _ := NewClient(token, baseURL, httpClient)

		method := ""
		pathUrl := "/users"

		req, err := client.MakeRequest(method, pathUrl, nil)

		assert.NoError(t, err, "Expected an error")
		assert.NotEqual(t, req.Method, "POST", "The expected method was not set correctly")
		assert.Equal(t, req.Method, "GET", "The expected method was set correctly")
	})
}

func TestBodyClose(t *testing.T) {
	tests := []struct {
		name     string
		closer   io.Closer
		expected error
	}{
		{
			name:     "Successful close",
			closer:   &mockCloser{},
			expected: nil,
		},
		{
			name:     "Error on close",
			closer:   &mockCloser{err: errors.New("close error")},
			expected: errors.New("close error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bodyClose(tt.closer)
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("Expected error %v, got %v", tt.expected, err)
			}
		})
	}
}

type mockCloser struct {
	err error
}

func (m *mockCloser) Close() error {
	return m.err
}
