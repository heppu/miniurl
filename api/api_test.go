package api_test

import (
	"errors"
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
	tests := []struct {
		name               string
		payload            string
		handler            api.Handler
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "OK",
			payload:            `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE"}`,
			handler:            &strHandler{str: "testvalue"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE", "hash":"testvalue"}`,
		},
		{
			name:               "invalid payload",
			payload:            `invalid json`,
			handler:            nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"msg": "bad request"}`,
		},
		{
			name:               "handler error",
			payload:            `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE"}`,
			handler:            &errHandler{err: errors.New("handler error")},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"msg": "internal server error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/url", strings.NewReader(tc.payload))
			rr := httptest.NewRecorder()

			r := httprouter.New()
			api.Bind(r, tc.handler)
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatusCode, rr.Result().StatusCode)
			body, err := io.ReadAll(rr.Result().Body)
			require.NoError(t, err)
			assert.JSONEq(t, tc.expectedBody, string(body))
		})
	}
}

type strHandler struct {
	str string
}

func (h *strHandler) AddUrl(url string) (hash string, err error) {
	return h.str, nil
}

type errHandler struct {
	err error
}

func (h *errHandler) AddUrl(url string) (hash string, err error) {
	return "", h.err
}
