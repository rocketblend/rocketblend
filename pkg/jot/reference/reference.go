package reference

type Reference string

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
