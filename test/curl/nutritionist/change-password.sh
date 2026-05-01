curl -v -X POST http://localhost:3000/nutritionist/complete-reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "'"$1"'",
    "password": "ResetPassword123!"
  }' | jq
