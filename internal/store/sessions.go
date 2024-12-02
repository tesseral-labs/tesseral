package store

import (
	openauthv1 "github.com/openauth/openauth/internal/gen/openauth/v1"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseSession(qSession *queries.Session) *openauthv1.Session {
	return &openauthv1.Session{
		Id:         idformat.Session.Format(qSession.ID),
		UserId:     idformat.User.Format(qSession.UserID),
		CreateTime: timestamppb.New(*qSession.CreateTime),
		ExpireTime: timestamppb.New(derefOrEmpty(qSession.ExpireTime)),
		Revoked:    qSession.Revoked,
	}
}
