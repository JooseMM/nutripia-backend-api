curl -v -X POST http://localhost:3000/nutritionist/verify-reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "'"$1"'"
  }' | jq
