package response

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Data    *T      `json:"data"`
}
