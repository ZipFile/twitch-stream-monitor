package testing

type UsernameResolver struct {
	Usernames map[string]string
	Error     error
}

func (ur *UsernameResolver) Resolve(username string) (string, error) {
	broadcasterID, _ := ur.Usernames[username]

	return broadcasterID, ur.Error
}
