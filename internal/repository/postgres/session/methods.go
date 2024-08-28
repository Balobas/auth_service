package sessionRepository

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *SessionRepository) CreateSession(ctx context.Context, session entity.Session) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(session)

	if err := r.Create(ctx, sessionRow); err != nil {
		return errors.Wrapf(
			err,
			"failed to create session with uid %s, user uid %s, device uid %s",
			session.Uid, session.UserUid, session.DeviceUid,
		)
	}
	return nil
}

func (r *SessionRepository) GetSessionByUid(ctx context.Context, uid uuid.UUID) (entity.Session, error) {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{Uid: uid})

	if err := r.GetOne(ctx, sessionRow, sessionRow.ConditionUidEqual()); err != nil {
		return entity.Session{}, errors.Wrapf(err, "failed to get session by uid %s", uid)
	}

	return sessionRow.ToEntity(), nil
}

func (r *SessionRepository) GetSessionsByUserUid(ctx context.Context, userUid uuid.UUID) ([]entity.Session, error) {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{UserUid: userUid})
	rows := pgEntity.NewSessionRows()

	if err := r.GetSome(ctx, sessionRow, rows, sessionRow.ConditionUserUidEqual()); err != nil {
		return nil, errors.Wrapf(err, "failed to get sessions by user uid %s", userUid)
	}

	return rows.ToEntities(), nil
}

func (r *SessionRepository) GetSessionByDeviceUid(ctx context.Context, deviceUid uuid.UUID) (entity.Session, error) {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{DeviceUid: deviceUid})

	if err := r.GetOne(ctx, sessionRow, sessionRow.ConditionDeviceUidEqual()); err != nil {
		return entity.Session{}, errors.Wrapf(err, "failed to get session by device uid %s", deviceUid)
	}

	return sessionRow.ToEntity(), nil
}

func (r *SessionRepository) DeleteSessionByUid(ctx context.Context, uid uuid.UUID) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{Uid: uid})

	if err := r.Delete(ctx, sessionRow, sessionRow.ConditionUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete session with uid %s", uid)
	}
	return nil
}

func (r *SessionRepository) DeleteSessionsByUserUid(ctx context.Context, userUid uuid.UUID) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{UserUid: userUid})

	if err := r.Delete(ctx, sessionRow, sessionRow.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete sessions with user uid %s", userUid)
	}
	return nil
}

func (r *SessionRepository) DeleteSessionByDeviceUid(ctx context.Context, deviceUid uuid.UUID) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{DeviceUid: deviceUid})

	if err := r.Delete(ctx, sessionRow, sessionRow.ConditionDeviceUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete session with device uid %s", deviceUid)
	}
	return nil
}
