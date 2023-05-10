package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/heppu/miniurl/api"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_AddUrl(t *testing.T) {
	const (
		payload            = `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE"}`
		expectedBody       = `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE", "hash": "testvalue"}`
		expectedStatusCode = http.StatusOK
	)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/url", strings.NewReader(payload))
	rr := httptest.NewRecorder()

	r := httprouter.New()
	api.Bind(r, nil)
	r.ServeHTTP(rr, req)

	assert.Equal(t, expectedStatusCode, rr.Result().StatusCode)
	body, err := io.ReadAll(rr.Result().Body)
	require.NoError(t, err)
	assert.Equal(t, expectedBody, string(body))
}
