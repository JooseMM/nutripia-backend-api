curl -v -X POST http://localhost:3000/nutritionist/client/$1 \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "juanete",
    "lastname": "moreno",
    "emailAddress": "john.doe@example.com",
    "birthDate": "1995-11-10T00:00:00Z"
  }' | jq
