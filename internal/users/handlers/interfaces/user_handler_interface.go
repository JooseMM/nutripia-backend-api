package interfaces

import "net/http"

type UserRepository interface {
	Create(w http.ResponseWriter, r *http.Request)
}
