package main

import (
	"log/slog"
	"net/http"
	"os"

	"utils.etin.dev/internal/email"
	"utils.etin.dev/internal/handlers"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	emailConfig, err := email.LoadConfig()
	if err != nil {
		slog.Error("Failed to load email configuration", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// Create the handler with the loaded configuration and the email sending function
	contactHandler := handlers.NewContactHandler(emailConfig, email.SendEmail)

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"project": "utils.etin.dev", "description": "Bespoke utilities exposed as endpoints for my use"}`))
	})

	mux.Handle("POST /asocleans/contact", contactHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
