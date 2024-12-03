package store

import (
	openauthv1 "github.com/openauth/openauth/internal/gen/openauth/v1"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
)

func parseSession(qSession *queries.Session) *openauthv1.Session {
	return &openauthv1.Session{
		Id:         idformat.Session.Format(qSession.ID),
		UserId:     idformat.User.Format(qSession.UserID),
		CreateTime: derefTimeOrNil(qSession.CreateTime),
		ExpireTime: derefTimeOrNil(qSession.ExpireTime),
		Revoked:    qSession.Revoked,
	}
}
