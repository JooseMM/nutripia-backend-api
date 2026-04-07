curl -v -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "juanete",
    "lastname": "moreno",
    "email_address": "john.doe@example.com",
    "password": "SecurePassword123!",
    "birth_date": "1995-11-10T00:00:00Z"
  }' | jq
