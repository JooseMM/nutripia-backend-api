package userModels

type Role int

const (
	CLIENT Role = iota
	NUTRITIONIST
	ADMIN
)
