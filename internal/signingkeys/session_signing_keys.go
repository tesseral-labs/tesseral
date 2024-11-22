package signingkeys

type sessionSigningKeys struct {
}

func NewSessionSigningKeys() *sessionSigningKeys {
	return &sessionSigningKeys{}
}

func (s *sessionSigningKeys) Create() error {
	return nil
}

func (s *sessionSigningKeys) Get() error {
	return nil
}
