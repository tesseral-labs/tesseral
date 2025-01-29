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
		qProject             queries.Project
		qOrg                 queries.Organization
		wantErr              bool
	}{
		{
			name: "google happy path",
			qIntermediateSession: queries.IntermediateSession{
				GoogleUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg:          queries.Organization{},
			qProject: queries.Project{
				LogInWithGoogleEnabled: true,
			},
			wantErr: false,
		},
		{
			name: "google email not verified",
			qIntermediateSession: queries.IntermediateSession{
				GoogleUserID: aws.String("foo"),
			},
			emailVerified: false,
			qOrg:          queries.Organization{},
			qProject: queries.Project{
				LogInWithGoogleEnabled: true,
			},
			wantErr: true,
		},
		{
			name: "google project not enabled",
			qIntermediateSession: queries.IntermediateSession{
				GoogleUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				OverrideLogInMethods:   true,
				DisableLogInWithGoogle: aws.Bool(false), // this shouldn't matter
			},
			qProject: queries.Project{
				LogInWithGoogleEnabled: false,
			},
			wantErr: true,
		},
		{
			name: "google org not enabled",
			qIntermediateSession: queries.IntermediateSession{
				GoogleUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				OverrideLogInMethods:   true,
				DisableLogInWithGoogle: aws.Bool(true),
			},
			qProject: queries.Project{
				LogInWithGoogleEnabled: true,
			},
			wantErr: true,
		},

		{
			name: "microsoft happy path",
			qIntermediateSession: queries.IntermediateSession{
				MicrosoftUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg:          queries.Organization{},
			qProject: queries.Project{
				LogInWithMicrosoftEnabled: true,
			},
			wantErr: false,
		},
		{
			name: "microsoft email not verified",
			qIntermediateSession: queries.IntermediateSession{
				MicrosoftUserID: aws.String("foo"),
			},
			emailVerified: false,
			qOrg:          queries.Organization{},
			qProject: queries.Project{
				LogInWithMicrosoftEnabled: true,
			},
			wantErr: true,
		},
		{
			name: "microsoft project not enabled",
			qIntermediateSession: queries.IntermediateSession{
				MicrosoftUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				OverrideLogInMethods:      true,
				DisableLogInWithMicrosoft: aws.Bool(false), // this shouldn't matter
			},
			qProject: queries.Project{
				LogInWithMicrosoftEnabled: false,
			},
			wantErr: true,
		},
		{
			name: "microsoft org not enabled",
			qIntermediateSession: queries.IntermediateSession{
				MicrosoftUserID: aws.String("foo"),
			},
			emailVerified: true,
			qOrg: queries.Organization{
				OverrideLogInMethods:      true,
				DisableLogInWithMicrosoft: aws.Bool(true),
			},
			qProject: queries.Project{
				LogInWithMicrosoftEnabled: true,
			},
			wantErr: true,
		},

		{
			name: "password happy path",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
			},
			emailVerified: true,
			qOrg:          queries.Organization{},
			qProject: queries.Project{
				LogInWithPasswordEnabled: true,
			},
			wantErr: false,
		},
		{
			name: "password email not verified",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
			},
			emailVerified: false,
			qOrg:          queries.Organization{},
			qProject: queries.Project{
				LogInWithPasswordEnabled: true,
			},
			wantErr: true,
		},
		{
			name: "password not verified",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: false,
			},
			emailVerified: true,
			qOrg:          queries.Organization{},
			qProject: queries.Project{
				LogInWithPasswordEnabled: true,
			},
			wantErr: true,
		},
		{
			name: "password project not enabled",
			qIntermediateSession: queries.IntermediateSession{
				PasswordVerified: true,
			},
			emailVerified: true,
			qOrg:          queries.Organization{},
			qProject: queries.Project{
				LogInWithPasswordEnabled: false,
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
				OverrideLogInMethods:     true,
				DisableLogInWithPassword: aws.Bool(true),
			},
			qProject: queries.Project{
				LogInWithPasswordEnabled: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAuthRequirementsSatisfiedInner(tt.qIntermediateSession, tt.emailVerified, tt.qProject, tt.qOrg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
