curl -v -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "juanete",
    "lastname": "moreno",
    "emailAddress": "john.doe@example.com",
    "password": "SecurePassword123!",
    "birthDate": "1995-11-10T00:00:00Z"
  }' | jq
