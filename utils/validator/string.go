package validator

import "regexp"

func String(str string, min int, max int) bool {
	return len(str) >= min && len(str) <= max && !BlankString(str)
}

func EmptyString(str string) bool {
	return len(str) <= 0
}

func BlankString(str string) bool {
	return str == " " || len(str) <= 0
}

func IsCep(cep string) bool {
	cepRegex := regexp.MustCompile(`^\d{5}-\d{3}`)
	return cepRegex.MatchString(cep)
}

func IsEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
