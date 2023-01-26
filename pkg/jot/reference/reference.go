package reference

import "fmt"

type Reference string

func Parse(referenceStr string) (Reference, error) {
	ref := Reference(referenceStr)
	if !ref.IsValid() {
		return "", fmt.Errorf("%q is not a valid reference string", referenceStr)
	}

	return ref, nil
}

func (r *Reference) Url() string {
	ret, _ := convertToUrl(string(*r))
	return ret.String()
}

func (r *Reference) String() string {
	return string(*r)
}

func (r *Reference) IsValid() bool {
	// TODO: validate this properly
	_, err := convertToUrl(string(*r))
	return err == nil
}
