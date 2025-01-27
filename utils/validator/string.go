package validator

import "regexp"

func String(str string, minV int, maxV int) bool {
	return len(str) >= minV && len(str) <= maxV && !BlankString(str)
}

func EmptyString(str string) bool {
	return len(str) <= 0
}

func BlankString(str string) bool {
	return str == " " || len(str) <= 0
}

func IsCep(cep string) bool {
	cepRegex := regexp.MustCompile(`^\d{5}-\d{3}$`)
	return cepRegex.MatchString(cep)
}

func IsEmail(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9]+([._-][0-9a-zA-Z]+)*@[a-zA-Z0-9]+([.-][0-9a-zA-Z]+)*\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(e)
}

func IsTelephone(t string) bool {
	telephoneRegex := regexp.MustCompile(`^(\(?\d{2}\)?\s)?(\d{4,5}-\d{4})$`)
	return telephoneRegex.MatchString(t)
}
