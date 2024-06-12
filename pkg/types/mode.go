package types

const (
	DefaultInjectionMode = RelaxedInjectionMode
)

type (
	InjectionMode string
)

const (
	EmptyInjectionMode   InjectionMode = ""
	StrictInjectionMode  InjectionMode = "strict"
	RelaxedInjectionMode InjectionMode = "relaxed"
	IgnoreInjectionMode  InjectionMode = "ignore"
)
