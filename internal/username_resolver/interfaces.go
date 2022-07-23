package username_resolver

type UsernameResolver interface {
	Resolve(string) (string, error)
}
