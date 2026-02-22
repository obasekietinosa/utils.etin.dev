package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"utils.etin.dev/internal/email"
	"utils.etin.dev/internal/handlers"
)

func TestContactHandler_ServeHTTP(t *testing.T) {
	mockEmailConfig := &email.Config{
		From: "test@example.com",
		To:   "admin@example.com",
	}

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSendEmail  func(cfg *email.Config, subject, body string) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Request",
			requestBody: map[string]interface{}{
				"fullName":            "John Doe",
				"location":            "Abuja",
				"address":             "123 Main St",
				"preferredDate":       "2023-10-27",
				"numberOfRooms":       "3",
				"estimatedSquareFeet": "2500",
				"agreement":           true,
			},
			mockSendEmail: func(cfg *email.Config, subject, body string) error {
				if subject != "New Quote Request from John Doe" {
					return errors.New("unexpected subject: " + subject)
				}
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Request received"}`,
		},
		{
			name: "Missing Agreement",
			requestBody: map[string]interface{}{
				"fullName": "John Doe",
				"location": "Abuja",
				"address":  "123 Main St",
			},
			mockSendEmail: func(cfg *email.Config, subject, body string) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "You must agree to the policy",
		},
		{
			name: "Missing Required Fields",
			requestBody: map[string]interface{}{
				"agreement": true,
			},
			mockSendEmail: func(cfg *email.Config, subject, body string) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing required fields",
		},
		{
			name: "Email Send Failure",
			requestBody: map[string]interface{}{
				"fullName":            "John Doe",
				"location":            "Abuja",
				"address":             "123 Main St",
				"preferredDate":       "2023-10-27",
				"numberOfRooms":       "3",
				"estimatedSquareFeet": "2500",
				"agreement":           true,
			},
			mockSendEmail: func(cfg *email.Config, subject, body string) error {
				return errors.New("smtp error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to send email",
		},
		{
			name: "Header Injection Attempt",
			requestBody: map[string]interface{}{
				"fullName":            "John Doe\nCc: hacker@example.com",
				"location":            "Abuja",
				"address":             "123 Main St",
				"preferredDate":       "2023-10-27",
				"numberOfRooms":       "3",
				"estimatedSquareFeet": "2500",
				"agreement":           true,
			},
			mockSendEmail: func(cfg *email.Config, subject, body string) error {
				if subject != "New Quote Request from John DoeCc: hacker@example.com" {
					return errors.New("subject was not sanitized: " + subject)
				}
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Request received"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/contact", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler := handlers.NewContactHandler(mockEmailConfig, tt.mockSendEmail)
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				// Trim newline usually added by http.Error
				gotBody := w.Body.String()
				// http.Error adds a newline
				if tt.expectedStatus != http.StatusOK {
					gotBody = gotBody[:len(gotBody)-1]
				}
				if gotBody != tt.expectedBody {
					t.Errorf("expected body %q, got %q", tt.expectedBody, gotBody)
				}
			}
		})
	}
}
