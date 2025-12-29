package sessionRepository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5"
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
			session.Uid, session.UserUid,
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

func (r *SessionRepository) GetSessionByUserUid(ctx context.Context, userUid uuid.UUID) (entity.Session, bool, error) {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{UserUid: userUid})

	if err := r.GetOne(ctx, sessionRow, sessionRow.ConditionUserUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Session{}, false, nil
		}
		log.Printf("failed to get session by user uid")
		return entity.Session{}, false, errors.Wrapf(err, "failed to get sessions by user uid %s", userUid)
	}

	log.Printf("successfuly get session by user uid")
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

func (r *SessionRepository) DeleteSessionByUid(ctx context.Context, uid uuid.UUID) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{Uid: uid})

	if err := r.Delete(ctx, sessionRow, sessionRow.ConditionUidEqual()); err != nil {
		log.Printf("failed to delete session by uid")
		return errors.Wrapf(err, "failed to delete session with uid %s", uid)
	}

	log.Printf("successfuly delete session by uid")
	return nil
}

func (r *SessionRepository) DeleteSessionByUserUid(ctx context.Context, userUid uuid.UUID) error {
	sessionRow := pgEntity.NewSessionRow().FromEntity(entity.Session{UserUid: userUid})

	if err := r.Delete(ctx, sessionRow, sessionRow.ConditionUserUidEqual()); err != nil {
		log.Printf("failed to delete session by user uid")
		return errors.Wrapf(err, "failed to delete sessions with user uid %s", userUid)
	}

	log.Printf("successfuly delete session by user uid")
	return nil
}

func (r *SessionRepository) DeleteSessionsByUsersUids(ctx context.Context, usersUids []uuid.UUID) error {
	if len(usersUids) == 0 {
		log.Printf("sessionsRepository.DeleteSessionsByUsersUids: empty usersUids")
		return errors.New("empty users uids")
	}
	sessionRow := pgEntity.NewSessionRow()
	args := make([]interface{}, len(usersUids))

	stmt := strings.Builder{}
	stmt.WriteString(fmt.Sprintf("delete from %s where user_uid in ($1", sessionRow.Table()))

	args[0] = pgtype.UUID{
		Bytes:  usersUids[0],
		Valid: true,
	}

	for i := 1; i < len(usersUids); i++ {
		stmt.WriteString(fmt.Sprintf(",$%d", i+1))
		args[i] = pgtype.UUID{
			Bytes:  usersUids[i],
			Valid: true,
		}
	}

	stmt.WriteByte(')')

	_, err := r.Exec(ctx, stmt.String(), args...)
	if err != nil {
		log.Printf("sessionsRepository.DeleteSessionsByUsersUids: empty usersUids")
		return nil
	}
	return nil
}
