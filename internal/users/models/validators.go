package userModels

func ValidateEmailAddress(email string) *string {
	return nil
}

func ValidateName(firstname string, lastname string) []string {
	var errList []string

	if firstname == "" {
		errList = append(errList, "firstname: is required")
	}

	if lastname == "" {
		errList = append(errList, "lastname: is required")
	}

	return errList
}

func ValidatePassword(password string) *string {
	if password == "" {
		var description = "lastname: is required"
		return &description
	}

	return nil
}
