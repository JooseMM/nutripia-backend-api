package authenticationTypes

type UserRoles int

const (
	CLIENT UserRoles = iota
	NUTRITIONIST
)

func (r *UserRoles) IsNutritionist() bool {
	return *r == NUTRITIONIST
}

func (r *UserRoles) IsClient() bool {
	return *r == CLIENT
}
