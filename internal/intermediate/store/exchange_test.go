package store

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/stretchr/testify/assert"
)

func TestStore_validateAuthRequirementsSatisfiedInner(t *testing.T) {
	testCases := []struct {
		name                 string
		qIntermediateSession queries.IntermediateSession
		emailVerified        bool
		qOrg                 queries.Organization
		wantErr              bool
	}{
		{
			name: "google happy path",
			qIntermediateSession: queries.IntermediateSession{
				GoogleUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithGoogle: true,
			},
			wantErr: false,
		},
		{
			name: "google email not verified",
			qIntermediateSession: queries.IntermediateSession{
				GoogleUserID: aws.String("foo"),
			},
			emailVerified: false,
			qOrg: queries.Organization{
				LogInWithGoogle: true,
			},
			wantErr: true,
		},
		{
			name: "google not enabled",
			qIntermediateSession: queries.IntermediateSession{
				GoogleUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithGoogle: false,
			},
			wantErr: true,
		},

		{
			name: "microsoft happy path",
			qIntermediateSession: queries.IntermediateSession{
				MicrosoftUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithMicrosoft: true,
			},
			wantErr: false,
		},
		{
			name: "microsoft email not verified",
			qIntermediateSession: queries.IntermediateSession{
				MicrosoftUserID: aws.String("foo"),
			},
			emailVerified: false,
			qOrg: queries.Organization{
				LogInWithMicrosoft: true,
			},
			wantErr: true,
		},
		{
			name: "microsoft not enabled",
			qIntermediateSession: queries.IntermediateSession{
				MicrosoftUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithMicrosoft: false,
			},
			wantErr: true,
		},

		{
			name: "password happy path",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithPassword: true,
			},
			wantErr: false,
		},
		{
			name: "password email not verified",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
			},
			emailVerified: false,
			qOrg: queries.Organization{
				LogInWithPassword: true,
			},
			wantErr: true,
		},
		{
			name: "password not verified",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: false,
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithPassword: true,
			},
			wantErr: true,
		},
		{
			name: "password org not enabled",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithPassword: false,
			},
			wantErr: true,
		},

		{
			name: "require mfa happy path passkey",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
				PasskeyVerified:  true,
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithPassword: true,
				LogInWithPasskey:  true,
				RequireMfa:        true,
			},
			wantErr: false,
		},
		{
			name: "require mfa happy path authenticator app",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified:         true,
				AuthenticatorAppVerified: true,
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithPassword:         true,
				LogInWithAuthenticatorApp: true,
				RequireMfa:                true,
			},
			wantErr: false,
		},
		{
			name: "require mfa no mfa",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithPassword: true,
				RequireMfa:        true,
			},
			wantErr: true,
		},
		{
			name: "require mfa not allowed mfa method",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
				PasskeyVerified:  true,
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithPassword:         true,
				LogInWithAuthenticatorApp: true,
				RequireMfa:                true,
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAuthRequirementsSatisfiedInner(tt.qIntermediateSession, tt.emailVerified, tt.qOrg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
