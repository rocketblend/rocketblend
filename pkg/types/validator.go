package types

type Validator interface {
	Validate(interface{}) error
}
