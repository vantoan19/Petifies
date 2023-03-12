package mapper

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user/entities"
	sqlc "github.com/vantoan19/Petifies/server/services/user-service/internal/infra/db/sqlc"
)

func DbSessionToEntity(s *sqlc.Session) entities.Session {
	return entities.Session{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpriresAt,
		ClientIP:     s.ClientIp,
		IsDisabled:   s.IsDisabled,
		CreatedAt:    s.CreatedAt,
	}
}

func EntitySessionToDb(s entities.Session) *sqlc.Session {
	return &sqlc.Session{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		ExpriresAt:   s.ExpiresAt,
		ClientIp:     s.ClientIP,
		IsDisabled:   s.IsDisabled,
		CreatedAt:    s.CreatedAt,
	}
}

func EntitySessionsToUpsertParams(sessions []entities.Session) *sqlc.BulkUpsertSessionsParams {
	return &sqlc.BulkUpsertSessionsParams{
		ID:           common.Map2(sessions, func(s entities.Session) uuid.UUID { return s.ID }),
		UserID:       common.Map2(sessions, func(s entities.Session) uuid.UUID { return s.UserID }),
		RefreshToken: common.Map2(sessions, func(s entities.Session) string { return s.RefreshToken }),
		ExpriresAt:   common.Map2(sessions, func(s entities.Session) time.Time { return s.ExpiresAt }),
		ClientIp:     common.Map2(sessions, func(s entities.Session) string { return s.ClientIP }),
		IsDisabled:   common.Map2(sessions, func(s entities.Session) bool { return s.IsDisabled }),
	}
}

func EntitySessionsToCreateParams(sessions []entities.Session) *sqlc.BulkCreateSessionParams {
	return &sqlc.BulkCreateSessionParams{
		ID:           common.Map2(sessions, func(s entities.Session) uuid.UUID { return s.ID }),
		UserID:       common.Map2(sessions, func(s entities.Session) uuid.UUID { return s.UserID }),
		RefreshToken: common.Map2(sessions, func(s entities.Session) string { return s.RefreshToken }),
		ExpriresAt:   common.Map2(sessions, func(s entities.Session) time.Time { return s.ExpiresAt }),
		ClientIp:     common.Map2(sessions, func(s entities.Session) string { return s.ClientIP }),
		IsDisabled:   common.Map2(sessions, func(s entities.Session) bool { return s.IsDisabled }),
	}
}
