package signingkeys

type SigningKeys struct {
	IntermediateSessions *intermediateSessionSigningKeys
	Sessions *sessionSigningKeys
}

type NewSigningKeysParams struct {
}

func NewSigningKeys(params *NewSigningKeysParams) *SigningKeys {
	return &SigningKeys{
		IntermediateSessions: NewIntermediateSessionSigningKeys(),
		Sessions: NewSessionSigningKeys(),
	}
}