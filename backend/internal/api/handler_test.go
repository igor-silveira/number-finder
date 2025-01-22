package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"number-finder-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type MockFinder struct {
	FindFunc func(target int, thresholdPercentage float64) (*service.Result, error)
}

func (m *MockFinder) Find(target int, thresholdPercentage float64) (*service.Result, error) {
	return m.FindFunc(target, thresholdPercentage)
}

func setupTest(mockFinder *MockFinder) (*fiber.App, *Handler) {
	app := fiber.New()
	handler := NewHandler(mockFinder, slog.Default())
	handler.RegisterRoutes(app)
	return app, handler
}

func performRequest(t *testing.T, app *fiber.App, method, url string) *http.Response {
	req := httptest.NewRequest(method, url, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	return resp
}

func assertResponse(t *testing.T, resp *http.Response, expectedStatusCode int, expectedResponse interface{}) {
	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	var response interface{}
	if expectedStatusCode == fiber.StatusOK {
		var res Response
		err := json.NewDecoder(resp.Body).Decode(&res)
		assert.NoError(t, err)
		response = res
	} else {
		var res map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&res)
		assert.NoError(t, err)
		response = res
	}

	assert.Equal(t, expectedResponse, response)
}

func TestHandler_handleFind(t *testing.T) {
	tests := []struct {
		name               string
		value              string
		threshold          string
		mockFindFunc       func(target int, thresholdPercentage float64) (*service.Result, error)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:      "Valid request with exact match",
			value:     "42",
			threshold: "0",
			mockFindFunc: func(target int, thresholdPercentage float64) (*service.Result, error) {
				return &service.Result{Index: 1, Value: 42, IsApproximate: false}, nil
			},
			expectedStatusCode: fiber.StatusOK,
			expectedResponse: Response{
				Index:         1,
				Value:         42,
				IsApproximate: false,
			},
		},
		{
			name:      "Valid request with approximate match",
			value:     "42",
			threshold: "10",
			mockFindFunc: func(target int, thresholdPercentage float64) (*service.Result, error) {
				return &service.Result{Index: 2, Value: 45, IsApproximate: true}, nil
			},
			expectedStatusCode: fiber.StatusOK,
			expectedResponse: Response{
				Index:         2,
				Value:         45,
				IsApproximate: true,
			},
		},
		{
			name:      "Invalid value parameter",
			value:     "not-a-number",
			threshold: "0",
			mockFindFunc: func(target int, thresholdPercentage float64) (*service.Result, error) {
				return nil, errors.New("invalid value parameter")
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   map[string]interface{}{"error": "Invalid value parameter"},
		},
		{
			name:      "Invalid thresholdPercentage parameter",
			value:     "42",
			threshold: "not-a-number",
			mockFindFunc: func(target int, thresholdPercentage float64) (*service.Result, error) {
				return nil, errors.New("invalid thresholdPercentage parameter")
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   map[string]interface{}{"error": "Invalid thresholdPercentage parameter"},
		},
		{
			name:      "Find operation error",
			value:     "42",
			threshold: "0",
			mockFindFunc: func(target int, thresholdPercentage float64) (*service.Result, error) {
				return nil, errors.New("value not found")
			},
			expectedStatusCode: fiber.StatusNotFound,
			expectedResponse:   map[string]interface{}{"message": "value not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFinder := &MockFinder{FindFunc: tt.mockFindFunc}
			app, _ := setupTest(mockFinder)

			url := "/api/find/" + tt.value + "?thresholdPercentage=" + tt.threshold
			resp := performRequest(t, app, "GET", url)

			assertResponse(t, resp, tt.expectedStatusCode, tt.expectedResponse)
		})
	}
}

func TestHandler_handleHealthCheck(t *testing.T) {
	app := fiber.New()
	handler := NewHandler(nil, slog.Default())

	handler.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"status": "ok"}, response)
}
