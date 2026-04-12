curl -v -X POST http://localhost:3000/nutritionist/send-reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "emailAddress": "josexmoreno1998@gmail.com"
  }' | jq
