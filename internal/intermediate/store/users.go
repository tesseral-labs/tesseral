package store

import (
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseUser(qUser queries.User) *intermediatev1.User {
	return &intermediatev1.User{
		Id:                  idformat.User.Format(qUser.ID),
		CreateTime:          timestamppb.New(*qUser.CreateTime),
		UpdateTime:          timestamppb.New(*qUser.UpdateTime),
		Email:               qUser.Email,
		Owner:               &qUser.IsOwner,
		GoogleUserId:        derefOrEmpty(qUser.GoogleUserID),
		MicrosoftUserId:     derefOrEmpty(qUser.MicrosoftUserID),
		GithubUserId:        derefOrEmpty(qUser.GithubUserID),
		HasAuthenticatorApp: qUser.AuthenticatorAppSecretCiphertext != nil,
		DisplayName:         qUser.DisplayName,
		ProfilePictureUrl:   qUser.ProfilePictureUrl,
	}
}
