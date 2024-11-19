package idformat

import "github.com/ssoready/prettyuuid"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

var (
	APIKey                			= prettyuuid.MustNewFormat("openauth_api_key_", alphabet)
	APISecretKey          			= prettyuuid.MustNewFormat("openauth_secret_", alphabet)
	IntermediateSession   			= prettyuuid.MustNewFormat("openauth_intermediate_session_", alphabet)
	MethodVerificationChallenge = prettyuuid.MustNewFormat("openauth_method_verification_challenge_", alphabet)
	Organization          			= prettyuuid.MustNewFormat("openauth_org_", alphabet)
	Project							  			= prettyuuid.MustNewFormat("openauth_project_", alphabet)
	Session               			= prettyuuid.MustNewFormat("openauth_session_", alphabet)
	User                  			= prettyuuid.MustNewFormat("openauth_user_", alphabet)
)
