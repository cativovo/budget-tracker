package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/require"
)

func testServerUrl() string {
	return fmt.Sprintf("http://127.0.0.1%s", testServer.Addr())
}

// fetch sends an HTTP request to the test server and returns the decoded success body, response, and decoded error body if the status code is >= 400.
func fetchFromServer[T any](t *testing.T, method string, path string, header http.Header, payload any) (T, *http.Response, *huma.ErrorModel) {
	return fetch[T](t, method, testServerUrl()+path, header, payload)
}

// fetch sends an HTTP request and returns the decoded success body, response, and decoded error body if the status code is >= 400.
func fetch[T any](t *testing.T, method string, url string, header http.Header, payload any) (T, *http.Response, *huma.ErrorModel) {
	t.Helper()

	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		require.NoError(t, err)
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)

	if header != nil {
		req.Header = header.Clone()
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	var data T

	if resp.StatusCode >= http.StatusBadRequest {
		em := new(huma.ErrorModel)
		err := json.NewDecoder(resp.Body).Decode(em)
		require.NoError(t, err)
		return data, resp, em
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	require.NoError(t, err)
	return data, resp, nil
}

func mustCreateCategory(t *testing.T, header http.Header) category.Category {
	t.Helper()

	c, resp, _ := fetchFromServer[category.Category](t, http.MethodPost, "/api/category", header, category.CreateCategoryInput{
		Name:  gofakeit.Noun(),
		Color: gofakeit.HexColor(),
		Icon:  gofakeit.Emoji(),
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Failed to create category")
	return c
}
