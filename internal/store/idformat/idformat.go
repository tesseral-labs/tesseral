package idformat

import "github.com/ssoready/prettyuuid"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

var (
	EmailVerificationChallenge    = prettyuuid.MustNewFormat("email_verification_challenge_", alphabet)
	IntermediateSession           = prettyuuid.MustNewFormat("intermediate_session_", alphabet)
	IntermediateSessionSigningKey = prettyuuid.MustNewFormat("intermediate_session_signing_key_", alphabet)
	Organization                  = prettyuuid.MustNewFormat("org_", alphabet)
	Project                       = prettyuuid.MustNewFormat("project_", alphabet)
	Session                       = prettyuuid.MustNewFormat("session_", alphabet)
	SessionRefreshToken           = prettyuuid.MustNewFormat("tesseral_secret_session_refresh_token_", alphabet)
	RelayedSessionToken           = prettyuuid.MustNewFormat("tesseral_secret_relayed_session_token_", alphabet)
	RelayedSessionRefreshToken    = prettyuuid.MustNewFormat("tesseral_secret_relayed_session_refresh_token_", alphabet)
	SessionSigningKey             = prettyuuid.MustNewFormat("session_signing_key_", alphabet)
	User                          = prettyuuid.MustNewFormat("user_", alphabet)
	VerifiedEmail                 = prettyuuid.MustNewFormat("verified_email_", alphabet)
	SAMLConnection                = prettyuuid.MustNewFormat("saml_connection_", alphabet)
	Passkey                       = prettyuuid.MustNewFormat("passkey_", alphabet)
	UserInvite                    = prettyuuid.MustNewFormat("user_invite_", alphabet)
	AuthenticatorAppRecoveryCode  = prettyuuid.MustNewFormat("authenticator_app_recovery_code_", alphabet)
	PasswordResetCode             = prettyuuid.MustNewFormat("password_reset_code_", alphabet)

	IntermediateSessionSecretToken = prettyuuid.MustNewFormat("tesseral_secret_intermediate_session_token_", alphabet)

	EmailVerificationChallengeCode = prettyuuid.MustNewFormat("email_verification_challenge_code_", alphabet)

	BackendAPIKey            = prettyuuid.MustNewFormat("backend_api_key_", alphabet)
	BackendAPIKeySecretToken = prettyuuid.MustNewFormat("tesseral_secret_key_", alphabet)

	PublishableKey = prettyuuid.MustNewFormat("publishable_key_", alphabet)

	SCIMAPIKey            = prettyuuid.MustNewFormat("scim_api_key_", alphabet)
	SCIMAPIKeySecretToken = prettyuuid.MustNewFormat("tesseral_secret_scim_api_key_", alphabet)

	UserImpersonationToken       = prettyuuid.MustNewFormat("user_impersonation_token_", alphabet)
	UserImpersonationSecretToken = prettyuuid.MustNewFormat("tesseral_secret_user_impersonation_token_", alphabet)

	ProjectRedirectURI = prettyuuid.MustNewFormat("project_redirect_uri_", alphabet)
	ProjectUISettings  = prettyuuid.MustNewFormat("project_ui_settings_", alphabet)

	ProjectWebhookSettings = prettyuuid.MustNewFormat("project_webhook_settings_", alphabet)
)
