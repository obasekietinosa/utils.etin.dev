package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"utils.etin.dev/internal/email"
)

type ContactFormRequest struct {
	FullName            string `json:"fullName"`
	Location            string `json:"location"`
	Address             string `json:"address"`
	PreferredDate       string `json:"preferredDate"`
	NumberOfRooms       string `json:"numberOfRooms"`
	EstimatedSquareFeet string `json:"estimatedSquareFeet"`
	Agreement           bool   `json:"agreement"`
}

type ContactHandler struct {
	EmailConfig *email.Config
	SendEmail   func(cfg *email.Config, subject, body string) error
}

func NewContactHandler(cfg *email.Config, sendEmail func(cfg *email.Config, subject, body string) error) *ContactHandler {
	return &ContactHandler{
		EmailConfig: cfg,
		SendEmail:   sendEmail,
	}
}

func (h *ContactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ContactFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !req.Agreement {
		http.Error(w, "You must agree to the policy", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.FullName == "" || req.Location == "" || req.Address == "" || req.PreferredDate == "" || req.NumberOfRooms == "" || req.EstimatedSquareFeet == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Sanitize FullName to prevent header injection
	sanitizedFullName := strings.ReplaceAll(req.FullName, "\r", "")
	sanitizedFullName = strings.ReplaceAll(sanitizedFullName, "\n", "")

	subject := fmt.Sprintf("New Quote Request from %s", sanitizedFullName)
	body := fmt.Sprintf(`
New Quote Request Received:

Full Name: %s
Location: %s
Address: %s
Preferred Date: %s
Number of Rooms: %s
Estimated Square Feet: %s

Agreed to Policy: %v
`, req.FullName, req.Location, req.Address, req.PreferredDate, req.NumberOfRooms, req.EstimatedSquareFeet, req.Agreement)

	if err := h.SendEmail(h.EmailConfig, subject, body); err != nil {
		slog.Error("Failed to send email", "error", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Request received"}`))
}
