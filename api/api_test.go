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
		handler            api.Handler
		payload            string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "OK",
			handler:            strHandler{str: "testvalue"},
			payload:            `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE", "hash": "testvalue"}`,
		},
		{
			name:               "Invalid payload",
			handler:            nil,
			payload:            `not valid json`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"msg": "bad request"}`,
		},
		{
			name:               "Error adding URL",
			handler:            errHandler{err: errors.New("handler error")},
			payload:            `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg": "internal server error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/url", strings.NewReader(tc.payload))
			rr := httptest.NewRecorder()

			router := httprouter.New()
			api.Bind(router, tc.handler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatusCode, rr.Result().StatusCode)
			body, err := io.ReadAll(rr.Result().Body)
			require.NoError(t, err)
			assert.JSONEq(t, tc.expectedResponse, string(body))
		})
	}
}

type strHandler struct {
	str string
}

func (h strHandler) AddUrl(string) (string, error) { return h.str, nil }

type errHandler struct {
	err error
}

func (h errHandler) AddUrl(string) (string, error) { return "", h.err }
