#!/bin/bash

echo "Testing Gin Basic API..."

echo "1. Testing health check endpoint:"
curl -X GET http://localhost:8080/health

echo -e "\n\n2. Testing user registration:"
curl -X POST http://localhost:8080/v1/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

echo -e "\n\n3. Testing user login:"
curl -X POST http://localhost:8080/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

echo -e "\n\n4. Testing user info retrieval:"
curl -X GET "http://localhost:8080/v1/user/info?email=test@example.com"

echo -e "\n\nTest completed."