package reference

type Reference string

func (r *Reference) Url() string {
	ret, _ := convertToUrl(string(*r))
	return ret.String()
}

func (r *Reference) String() string {
	return string(*r)
}
