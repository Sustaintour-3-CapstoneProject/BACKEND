package controllers_test

import (
	"backend/config"
	"backend/controllers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	config.TestInitDB()

	e := echo.New()

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
		expectedKey  string
	}{
		{
			name: "Successful Registration",
			input: map[string]interface{}{
				"username":   "testuser10",
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "testuser10@example.com",
				"city":       "New York",
				"password":   "password123",
				"role":       "",
			},
			expectedCode: http.StatusOK,
			expectedKey:  "token",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			body, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = controllers.RegisterHandler(c)
			if err != nil {
				t.Fatalf("Handler error: %v", err)
			}

			assert.Equal(t, tc.expectedCode, rec.Code)

			var response map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			assert.Contains(t, response, "meta")
			meta := response["meta"].(map[string]interface{})
			assert.Equal(t, tc.expectedCode, int(meta["code"].(float64)))

			if tc.expectedKey != "" {
				assert.Contains(t, response["data"], tc.expectedKey)
			}
		})
	}
}
