curl -v http://localhost:3000/clients/$1 \
  -H "Content-Type: application/json" | jq
