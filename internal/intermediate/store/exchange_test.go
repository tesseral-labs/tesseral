package store

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
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
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorGoogle),
				GoogleUserID:      aws.String("foo"),
				Email:             aws.String("foo@bar.com"),
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
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorGoogle),
				GoogleUserID:      aws.String("foo"),
				Email:             aws.String("foo@bar.com"),
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
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorGoogle),
				GoogleUserID:      aws.String("foo"),
				Email:             aws.String("foo@bar.com"),
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
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorMicrosoft),
				MicrosoftUserID:   aws.String("foo"),
				Email:             aws.String("foo@bar.com"),
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
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorMicrosoft),
				MicrosoftUserID:   aws.String("foo"),
				Email:             aws.String("foo@bar.com"),
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
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorMicrosoft),
				MicrosoftUserID:   aws.String("foo"),
				Email:             aws.String("foo@bar.com"),
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
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorEmail),
				PasswordVerified:  true,
				Email:             aws.String("foo@bar.com"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithEmail:    true,
				LogInWithPassword: true,
			},
			wantErr: false,
		},
		{
			name: "password email not verified",
			qIntermediateSession: queries.IntermediateSession{
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorEmail),
				PasswordVerified:  true,
				Email:             aws.String("foo@bar.com"),
			},
			emailVerified: false,
			qOrg: queries.Organization{
				LogInWithEmail:    true,
				LogInWithPassword: true,
			},
			wantErr: true,
		},
		{
			name: "password not verified",
			qIntermediateSession: queries.IntermediateSession{
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorEmail),
				PasswordVerified:  false,
				Email:             aws.String("foo@bar.com"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithEmail:    true,
				LogInWithPassword: true,
			},
			wantErr: true,
		},

		{
			name: "require mfa happy path passkey",
			qIntermediateSession: queries.IntermediateSession{
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorEmail),
				PasswordVerified:  true,
				PasskeyVerified:   true,
				Email:             aws.String("foo@bar.com"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithEmail:    true,
				LogInWithPassword: true,
				LogInWithPasskey:  true,
				RequireMfa:        true,
			},
			wantErr: false,
		},
		{
			name: "require mfa happy path authenticator app",
			qIntermediateSession: queries.IntermediateSession{
				PrimaryAuthFactor:        primaryAuthFactor(queries.PrimaryAuthFactorEmail),
				PasswordVerified:         true,
				AuthenticatorAppVerified: true,
				Email:                    aws.String("foo@bar.com"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithEmail:            true,
				LogInWithPassword:         true,
				LogInWithAuthenticatorApp: true,
				RequireMfa:                true,
			},
			wantErr: false,
		},
		{
			name: "require mfa no mfa",
			qIntermediateSession: queries.IntermediateSession{
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorEmail),
				PasswordVerified:  true,
				Email:             aws.String("foo@bar.com"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithEmail:    true,
				LogInWithPassword: true,
				RequireMfa:        true,
			},
			wantErr: true,
		},
		{
			name: "require mfa not allowed mfa method",
			qIntermediateSession: queries.IntermediateSession{
				PrimaryAuthFactor: primaryAuthFactor(queries.PrimaryAuthFactorEmail),
				PasswordVerified:  true,
				PasskeyVerified:   true,
				Email:             aws.String("foo@bar.com"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				LogInWithEmail:            true,
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

func primaryAuthFactor(v queries.PrimaryAuthFactor) *queries.PrimaryAuthFactor {
	return &v
}
