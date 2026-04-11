curl -v -X POST http://localhost:3000/nutritionist/login \
  -H "Content-Type: application/json" \
  -d '{
    "emailAddress": "john.doe@example.com",
    "password": "SecurePassword123!"
  }' | jq
