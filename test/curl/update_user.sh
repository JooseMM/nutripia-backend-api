curl -v -X PUT http://localhost:3000/nutritionist/$1 \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "update",
    "lastname": "moreno",
    "emailAddress": "john.doe.updated@example.com",
    "password": "SecurePassword123!",
    "birthDate": "1998-11-10T00:00:00Z"
  }' | jq
