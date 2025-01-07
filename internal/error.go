package internal

type Causes struct {
	Field   string
	Message string
}

type DomainError struct {
	Message string
	Causes  []Causes
}

func (d DomainError) Error() string {
	return d.Message
}
