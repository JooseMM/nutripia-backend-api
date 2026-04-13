package core

import "time"

func ValidateEmailAddress(email string) []string {
	var errList []string
	if email == "" {
		errList = append(errList, "email: is required")
	}
	return errList
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

func ValidatePassword(password string) []string {
	var errList []string

	if password == "" {
		errList = append(errList, "password: is required")
	}

	return errList
}

func ValidateBirthDate(birthDate time.Time) []string {
	var errList []string
	now := time.Now().UTC()

	if now.Before(birthDate) {
		errList = append(errList, "birthDate: cannot be in the future.")
	}

	if len(errList) == 0 {
		return nil
	}

	return errList
}

func ValidateRUT(rut *string) []string {
	var errList []string

	if *rut == "" {
		errList = append(errList, "rut: RUT is required")
	}

	if len(errList) == 0 {
		return nil
	}

	return errList
}
