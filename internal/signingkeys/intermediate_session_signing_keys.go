package signingkeys

type intermediateSessionSigningKeys struct{}

func NewIntermediateSessionSigningKeys() *intermediateSessionSigningKeys {
	return &intermediateSessionSigningKeys{}
}

func (i *intermediateSessionSigningKeys) Create() error {
	return nil
}

func (i *intermediateSessionSigningKeys) Get() error {
	return nil
}
