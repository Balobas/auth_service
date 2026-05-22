package sessionRepository

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *SessionRepository) CreateSession(ctx context.Context, session entity.Session) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(session)

	if err := r.Create(ctx, sessionRow); err != nil {
		log.Printf("failed to create session")
		return errors.Wrapf(
			err,
			"failed to create session with uid %s, user uid %s",
			session.Uid, session.MaintainerUid,
		)
	}
	log.Printf("successfuly create session")
	return nil
}

func (r *SessionRepository) GetSessionByUid(ctx context.Context, uid uuid.UUID) (entity.Session, bool, error) {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{Uid: uid})

	if err := r.GetOne(ctx, sessionRow, sessionRow.ConditionUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Session{}, false, nil
		}
		log.Printf("failed to get session by uid")
		return entity.Session{}, false, errors.Wrapf(err, "failed to get session by uid %s", uid)
	}

	log.Printf("successfuly get session by uid")
	return sessionRow.ToEntity(), true, nil
}

func (r *SessionRepository) GetSessionByMaintainerAndDevice(ctx context.Context, maintainerUid uuid.UUID, deviceUid uuid.UUID) (entity.Session, bool, error) {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{MaintainerUid: maintainerUid, DeviceUid: deviceUid})

	if err := r.GetOne(ctx, sessionRow, sessionRow.ConditionMaintainerUidAndDeviceUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Session{}, false, nil
		}
		log.Printf("failed to get session by user uid and device uid")
		return entity.Session{}, false, errors.Wrapf(err, "failed to get sessions by maintainer uid %s and device uid %s", maintainerUid, deviceUid)
	}

	log.Printf("successfuly get session by user uid and device uid")
	return sessionRow.ToEntity(), true, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session entity.Session) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(session)

	if err := r.Update(ctx, sessionRow, sessionRow.ConditionUidEqual()); err != nil {
		log.Printf("failed to update session")
		return errors.Wrapf(err, "failed to update session with uid %s", session.Uid)
	}

	log.Printf("successfuly update session")
	return nil
}

func (r *SessionRepository) DeleteUserSessions(ctx context.Context, userUid uuid.UUID, excludedUids ...uuid.UUID) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{MaintainerUid: userUid})

	condition := sq.And{
		sessionRow.ConditionMaintainerUidEqual(),
		sessionRow.ConditionTypeEqual(entity.SessionTypeUser),
	}

	if len(excludedUids) != 0 {
		pgUids := make([]pgtype.UUID, len(excludedUids))
		for i := 0; i < len(excludedUids); i++ {
			pgUids[i] = pgtype.UUID{
				Bytes: excludedUids[i],
				Valid: true,
			}
		}

		condition = append(condition, sq.NotEq{
			"uid": pgUids,
		})
	}

	if err := r.Delete(ctx, sessionRow, condition); err != nil {
		log.Printf("failed to delete sessions by user uid")
		return errors.Wrapf(err, "failed to delete sessions with user uid %s", userUid)
	}

	log.Printf("successfuly delete sessions by user uid")
	return nil
}

func (r *SessionRepository) DeleteSessionByMaintainerAndDevice(ctx context.Context, maintainerUid uuid.UUID, deviceUid uuid.UUID) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{MaintainerUid: maintainerUid, DeviceUid: deviceUid})

	if err := r.Delete(ctx, sessionRow, sessionRow.ConditionMaintainerUidAndDeviceUidEqual()); err != nil {
		log.Printf("failed to delete session by user uid and device uid")
		return errors.Wrapf(err, "failed to delete sessions with maintainer uid %s and device uid %s", maintainerUid, deviceUid)
	}

	log.Printf("successfuly delete session by user uid and device uid")
	return nil
}

// func (r *SessionRepository) DeleteSessionsByUsersUids(ctx context.Context, usersUids []uuid.UUID) error {
// 	if len(usersUids) == 0 {
// 		log.Printf("sessionsRepository.DeleteSessionsByUsersUids: empty usersUids")
// 		return errors.New("empty users uids")
// 	}
// 	sessionRow := pgEntity.NewSessionRow()
// 	args := make([]interface{}, len(usersUids))

// 	stmt := strings.Builder{}
// 	stmt.WriteString(fmt.Sprintf("delete from %s where user_uid in ($1", sessionRow.Table()))

// 	args[0] = pgtype.UUID{
// 		Bytes: usersUids[0],
// 		Valid: true,
// 	}

// 	for i := 1; i < len(usersUids); i++ {
// 		stmt.WriteString(fmt.Sprintf(",$%d", i+1))
// 		args[i] = pgtype.UUID{
// 			Bytes: usersUids[i],
// 			Valid: true,
// 		}
// 	}

// 	stmt.WriteByte(')')

// 	_, err := r.Exec(ctx, stmt.String(), args...)
// 	if err != nil {
// 		log.Printf("sessionsRepository.DeleteSessionsByUsersUids: empty usersUids")
// 		return nil
// 	}
// 	return nil
// }
