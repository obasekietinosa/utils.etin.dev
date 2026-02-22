#!/bin/bash

# Start the server in the background
# export SMTP_HOST=smtp.example.com
# export SMTP_PORT=587
# export SMTP_USERNAME=user
# export SMTP_PASSWORD=pass
# export SMTP_FROM_EMAIL=from@example.com
# export SMTP_TO_EMAIL=to@example.com
# export PORT=8080

# go run ./cmd/server &
# SERVER_PID=$!
# sleep 2

echo "Sending test request..."
curl -X POST http://localhost:8080/contact \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Test User",
    "location": "Lagos",
    "address": "123 Test Lane",
    "preferredDate": "2023-11-01",
    "numberOfRooms": "2",
    "estimatedSquareFeet": "1500",
    "agreement": true
  }'

# kill $SERVER_PID
