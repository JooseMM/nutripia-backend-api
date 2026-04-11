curl -v http://localhost:3000/clients/ \
  -H "Content-Type: application/json" \
  -H "x-session-token: $1" | jq
