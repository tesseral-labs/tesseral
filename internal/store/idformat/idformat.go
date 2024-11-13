package idformat

import "github.com/ssoready/prettyuuid"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

var (
	APIKey                = prettyuuid.MustNewFormat("api_key_", alphabet)
	APISecretKey          = prettyuuid.MustNewFormat("openauth_sk_", alphabet)
	Organization          = prettyuuid.MustNewFormat("org_", alphabet)
	Project							  = prettyuuid.MustNewFormat("project_", alphabet)
	ProjectOrganization		= prettyuuid.MustNewFormat("project_org_", alphabet)
	ProjectUser						= prettyuuid.MustNewFormat("project_user_", alphabet)
	User                  = prettyuuid.MustNewFormat("user_", alphabet)
)
