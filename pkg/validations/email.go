package validations

import "fmt"

func ValidateEmail(email string) error {
	if len(email) == 0 {
		return fmt.Errorf("empty email")
	}
	return nil
}

func ValidatePhone(phone string) error {
	return nil
}
