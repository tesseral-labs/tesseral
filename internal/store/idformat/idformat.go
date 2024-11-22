package idformat

import "github.com/ssoready/prettyuuid"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

var (
	APIKey                        = prettyuuid.MustNewFormat("api_key_", alphabet)
	APISecretKey                  = prettyuuid.MustNewFormat("openauth_secret_key_", alphabet)
	IntermediateSession           = prettyuuid.MustNewFormat("intermediate_session_", alphabet)
	IntermediateSessionSigningKey = prettyuuid.MustNewFormat("intermediate_session_signing_key_", alphabet)
	MethodVerificationChallenge   = prettyuuid.MustNewFormat("method_verification_challenge_", alphabet)
	Organization                  = prettyuuid.MustNewFormat("org_", alphabet)
	Project                       = prettyuuid.MustNewFormat("project_", alphabet)
	Session                       = prettyuuid.MustNewFormat("session_", alphabet)
	SessionSigningKey             = prettyuuid.MustNewFormat("session_signing_key_", alphabet)
	User                          = prettyuuid.MustNewFormat("user_", alphabet)
)
