package authenticationDtos

type SendResetPasswordTokenRequest struct {
	EmailAddress string `json:"emailAddress"`
}
