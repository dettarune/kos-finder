package model

type ValidationError struct {
	Status bool
	Message string
	StatusCode int
	Errors interface{}
}

func (v *ValidationError) Error() string {
	return v.Message
}