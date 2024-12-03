package store

import (
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func parseSession(qSession *queries.Session) *frontendv1.Session {
	return &frontendv1.Session{
		Id:         idformat.Session.Format(qSession.ID),
		UserId:     idformat.User.Format(qSession.UserID),
		CreateTime: derefTimeOrNil(qSession.CreateTime),
		ExpireTime: derefTimeOrNil(qSession.ExpireTime),
		Revoked:    qSession.Revoked,
	}
}
