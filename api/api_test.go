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
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			name:               "OK",
			handler:            strHandler{str: "testvalue"},
			payload:            `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE"}`,
			expectedResponse:   `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE", "hash": "testvalue"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid payload",
			handler:            nil,
			payload:            `not valid json`,
			expectedResponse:   `{"msg": "bad request"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Error in handler",
			handler:            errHandler{err: errors.New("handler error")},
			payload:            `{"url": "https://github.com/gourses/miniurl/blob/main/LICENSE"}`,
			expectedResponse:   `{"msg": "internal server error"}`,
			expectedStatusCode: http.StatusInternalServerError,
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
